package user

import (
	"context"
	"errors"
	"fmt"
	"gin-gonic-gom/Models"
	"gin-gonic-gom/common"
	"gin-gonic-gom/constant"
	"gin-gonic-gom/helper"
	"gin-gonic-gom/utils"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserImplementService struct {
	usercollection  *mongo.Collection
	tokencollection *mongo.Collection
	otpcollection   *mongo.Collection
	majorcollection *mongo.Collection
	ctx             context.Context
}

func NewUserService(usercollection *mongo.Collection, majorcollection *mongo.Collection, tokencollection *mongo.Collection, otpcollection *mongo.Collection, ctx context.Context) UserService {
	return &UserImplementService{
		usercollection:  usercollection,
		tokencollection: tokencollection,
		otpcollection:   otpcollection,
		majorcollection: majorcollection,
		ctx:             ctx,
	}
}
func (a *UserImplementService) FindOneData(collection *mongo.Collection, filter interface{}, result interface{}, opts ...*options.FindOneOptions) error {
	err := collection.FindOne(a.ctx, filter, opts...).Decode(result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("no document found")
		}
		return err
	}
	return nil
}
func (a *UserImplementService) CreateUser(input *Models.CreateUserInput) error {
	//user.Id = primitive.NewObjectID()
	var user Models.UserModel
	var err error
	timeLocalHoChiMinh, _ := utils.GetCurrentTimeInLocal("Asia/Ho_Chi_Minh")
	user = Models.UserModel{
		Id:             primitive.NewObjectID(),
		MajorId:        nil,
		Role:           input.Role,
		Name:           input.Name,
		Email:          input.Email,
		Password:       input.Password,
		Avatar:         input.Avatar,
		Phone:          input.Phone,
		Address:        "",
		Department:     "",
		Gender:         nil,
		DateOfBirth:    timeLocalHoChiMinh,
		EnrollmentDate: timeLocalHoChiMinh,
		HireDate:       timeLocalHoChiMinh,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	_, err = a.usercollection.InsertOne(a.ctx, user)
	if err != nil {
		return fmt.Errorf("Error inserting user: %v", err)
	}
	return err
}
func (a *UserImplementService) CheckExistEmail(email string) (bool, error) {
	var user *Models.UserModel
	query := bson.M{"email": email}
	err := a.usercollection.FindOne(a.ctx, query).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func (a *UserImplementService) DeleteTokenExp() {
	expiryThreshold := time.Now().UTC()
	_, err := a.tokencollection.DeleteMany(a.ctx, bson.M{"exp": bson.M{"$lte": expiryThreshold}})
	if err != nil {
		log.Printf("Error deleting expired tokens: %v", err)
	}
	fmt.Println("Cleaned expired tokens ----------> at", expiryThreshold)
}
func (a *UserImplementService) DeleteOTPExp() {
	expiryThreshold := time.Now().UTC()
	_, err := a.otpcollection.DeleteMany(a.ctx, bson.M{"expires_at": bson.M{"$lte": expiryThreshold}})
	if err != nil {
		log.Printf("Error deleting expired OTP: %v", err)
	}
}
func (a *UserImplementService) LoginUser(authInput *Models.AuthInput, c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var foundUser Models.UserModel
	deviced := c.Request.Header.Get("User-Agent")
	expAccessToken := time.Minute * 30
	expRefreshToken := time.Hour * 24 * 3
	err := a.usercollection.FindOne(ctx, bson.M{"email": authInput.Email}).Decode(&foundUser)
	defer cancel()
	if err != nil {
		common.NewErrorResponse(c, http.StatusUnauthorized, common.ErrorEmailOrPassword, err.Error())
		return
	}
	passwordIsValid := utils.CheckPasswordHash(authInput.Password, foundUser.Password)
	if passwordIsValid != true {
		common.NewErrorResponse(c, http.StatusUnauthorized, common.ErrorPassword, "Mật khẩu không hợp lệ")
		return
	}
	defer cancel()
	if err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorInternetServer, err.Error())
		return
	}
	// Generator token and Deviced
	access_token, err_access_token := utils.GenerateAccessToken(expAccessToken, &foundUser, os.Getenv("ACCESS_TOKEN_SECRET"))
	helper.ErrorPanic(err_access_token)
	refresh_token, err_refresh_token := utils.GenerateRefreshToken(expRefreshToken, &foundUser, os.Getenv("REFRESH_TOKEN_SECRET"))
	helper.ErrorPanic(err_refresh_token)
	// set cookie
	c.SetCookie("access_token", access_token, 3600, "/", "localhost", false, true)
	// cap nhat hoac them token neu chua ton tai
	filterDeviced := bson.D{{"deviced", deviced}}
	updateToken := bson.D{{"$set", bson.D{
		{"user_id", foundUser.Id},
		{"refresh_token", refresh_token},
		{"exp", utils.ConvertDurationToTimeUTC(expRefreshToken)},
		{"deviced", deviced},
		{"created_at", time.Now().UTC()},
		{"updated_at", time.Now().UTC()},
	}}}
	_, _ = a.tokencollection.UpdateOne(ctx, filterDeviced, updateToken, options.Update().SetUpsert(true))
	// dem token
	filter := bson.M{"user_id": foundUser.Id}
	opts := options.Find().SetSort(bson.D{{"created_at", 1}})

	cursor, err := a.tokencollection.Find(ctx, filter, opts)
	if err != nil {
		common.NewErrorResponse(c, http.StatusNotFound, common.ErrorFindUser, err.Error())
		return
	}
	defer cursor.Close(ctx)

	var tokens []Models.TokenModel
	if err := cursor.All(ctx, &tokens); err != nil {
		common.NewErrorResponse(c, http.StatusInternalServerError, common.ErrorInternetServer, err.Error())
		return
	}
	numDevices := len(tokens)
	// kiem tra token >= 2 --> xoa cac token dau tien va chi giu lai 2 token cuoi
	if numDevices >= 2 {
		devicesToDelete := tokens[:numDevices-2]
		if len(devicesToDelete) > 0 {
			var idsToDelete []primitive.ObjectID
			for _, token := range devicesToDelete {
				idsToDelete = append(idsToDelete, token.Id)
			}
			_, err := a.tokencollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": idsToDelete}})
			if err != nil {
				common.NewErrorResponse(c, http.StatusBadRequest, common.ErrorDeleteTokenData, err.Error())
				return
			}
		}
	}
	token := Models.Token{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}
	c.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, common.SuccessLogin, token))
}
func (a *UserImplementService) GetMe(userId, role string) (*Models.UserModel, error) {
	var account *Models.UserModel
	query := bson.D{bson.E{Key: "_id", Value: utils.ConvertStringToObjectId(userId)}}
	var err error
	switch role {
	case "student":
		studentProjection := bson.D{{"password", 0}, {"hire_date", 0}, {"department", 0}}
		err = a.usercollection.FindOne(a.ctx, query, options.FindOne().SetProjection(studentProjection)).Decode(&account)
	case "teacher":
		teacherProjection := bson.D{{"password", 0}, {"address", 0}, {"major_id", 0}, {"enrollment_date", 0}}
		err = a.usercollection.FindOne(a.ctx, query, options.FindOne().SetProjection(teacherProjection)).Decode(&account)
	case "admin":
		adminProjection := bson.D{{"password", 0}, {"address", 0}, {"major_id", 0}, {"enrollment_date", 0}, {"hire_date", 0}, {"department", 0}}
		err = a.usercollection.FindOne(a.ctx, query, options.FindOne().SetProjection(adminProjection)).Decode(&account)
	default:
		return nil, errors.New("Invalid role")
	}
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Không tìm thấy tài liệu
			log.Printf("Không tìm thấy user!")
		} else {
			log.Fatal(err)
		}
	}
	return account, err
}
func (a *UserImplementService) GetAccount(userId string) (*Models.UserModel, error) {
	var account *Models.UserModel
	query := bson.D{bson.E{Key: "_id", Value: utils.ConvertStringToObjectId(userId)}}
	err := a.usercollection.FindOne(a.ctx, query, options.FindOne().SetProjection(bson.D{{"password", 0}})).Decode(&account)
	return account, err
}
func (a *UserImplementService) GetAll(page, limit int) ([]*Models.UserModel, int, error) {
	skip := limit * (page - 1)
	opts := options.Find().SetProjection(bson.D{{"password", 0}}).SetSkip(int64(skip)).SetLimit(int64(limit))
	total, err := a.usercollection.CountDocuments(a.ctx, bson.M{})
	cur, err := a.usercollection.Find(a.ctx, bson.M{}, opts)
	defer cur.Close(a.ctx)
	var users []*Models.UserModel
	for cur.Next(a.ctx) {
		var user *Models.UserModel
		err := cur.Decode(&user)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return nil, 0, err
	}
	return users, int(total), err
}

func (a *UserImplementService) Logout(deviced string, userId primitive.ObjectID) error {
	filter := bson.M{"deviced": deviced, "user_id": userId}
	_, err := a.tokencollection.DeleteOne(a.ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (a *UserImplementService) UpdateMe(userId string, userData *Models.UserUpdate) (*Models.UserModel, error) {
	filter := bson.M{"_id": utils.ConvertStringToObjectId(userId)}
	updateFields := bson.M{}

	if userData.Name != "" {
		updateFields["name"] = userData.Name
	}
	if userData.Gender != constant.SEX_OTHER || userData.Gender != constant.SEX_MALE || userData.Gender != constant.SEX_FEMALE {
		updateFields["gender"] = userData.Gender
	}
	if userData.Avatar != "" {
		updateFields["avatar"] = userData.Avatar
	}
	if userData.Address != "" {
		updateFields["address"] = userData.Address
	}
	if userData.Phone != "" {
		updateFields["phone"] = userData.Phone
	}
	if userData.Department != "" {
		updateFields["department"] = userData.Department
	}
	if !userData.DateOfBirth.IsZero() {
		updateFields["date_of_birth"] = userData.DateOfBirth
	}
	if !userData.HireDate.IsZero() {
		updateFields["hire_date"] = userData.HireDate
	}
	if !userData.EnrollmentDate.IsZero() {
		updateFields["enrollment_date"] = userData.EnrollmentDate
	}
	updateFields["updated_at"] = time.Now().UTC()
	if len(updateFields) > 0 {
		userDataUpdate := bson.M{"$set": updateFields}

		var userRes *Models.UserModel
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetProjection(bson.M{"password": 0, "role": 0})
		err := a.usercollection.FindOneAndUpdate(a.ctx, filter, userDataUpdate, opts).Decode(&userRes)
		if err != nil {
			return nil, err
		}
		return userRes, nil
	}

	return nil, fmt.Errorf("Hủy update!")
}
func (a *UserImplementService) UpdateUser(account *Models.AccountUpdate, id primitive.ObjectID) (*Models.UserModel, error) {
	var err error
	filter := bson.M{"_id": id}
	passwordHash, _ := utils.HashPassword(account.Password)
	updateFields := bson.M{}
	var userFind *Models.UserModel
	err = a.FindOneData(a.usercollection, filter, &userFind)
	if userFind.Role == "student" {
		updateFields["major_id"] = account.MajorId
	}
	if account.Name != "" {
		updateFields["name"] = account.Name
	}
	if account.Gender != constant.SEX_OTHER || account.Gender != constant.SEX_MALE || account.Gender != constant.SEX_FEMALE {
		updateFields["gender"] = account.Gender
	}
	if account.Email != "" {
		updateFields["email"] = account.Email
	}
	if account.Password != "" {
		updateFields["password"] = passwordHash
	}
	if account.Role != "" {
		updateFields["role_type"] = account.Role
	}
	if account.Avatar != "" {
		updateFields["avatar"] = account.Avatar
	}
	if account.Address != "" {
		updateFields["address"] = account.Address
	}
	if account.Phone != "" {
		updateFields["phone"] = account.Phone
	}
	if account.Department != "" {
		updateFields["department"] = account.Department
	}
	if !account.DateOfBirth.IsZero() {
		updateFields["date_of_birth"] = account.DateOfBirth
	}
	if !account.HireDate.IsZero() {
		updateFields["hire_date"] = account.HireDate
	}
	if !account.EnrollmentDate.IsZero() {
		updateFields["enrollment_date"] = account.EnrollmentDate
	}

	updateFields["updated_at"] = time.Now().UTC()
	if len(updateFields) > 0 {
		userDataUpdate := bson.M{"$set": updateFields}
		var userRes *Models.UserModel
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		err = a.usercollection.FindOneAndUpdate(a.ctx, filter, userDataUpdate, opts).Decode(&userRes)
		if err != nil {
			return nil, err
		}
		return userRes, nil
	}
	return nil, fmt.Errorf("Hủy update!")
}
func (a *UserImplementService) DeleteUser(userId string) (int, error) {
	filter := bson.D{{"_id", utils.ConvertStringToObjectId(userId)}}
	res, err := a.usercollection.DeleteOne(a.ctx, filter)
	return int(res.DeletedCount), err
}
func (a *UserImplementService) ChangePassword(userId string, passwordInput *Models.ChangePasswordInput) error {
	// Kiểm tra tra mật khẩu mới và confirmPassword
	if passwordInput.NewPassword != passwordInput.ConfirmNewPassword {
		return errors.New("Mật khẩu mới không khớp")
	}
	var user *Models.UserModel
	filter := bson.M{"_id": utils.ConvertStringToObjectId(userId)}
	err := a.usercollection.FindOne(a.ctx, filter).Decode(&user)
	if err != nil {
		return errors.New("Không tìm thấy user!")
	}
	// kiểm tra mật khâ cũ
	if !utils.CheckPasswordHash(passwordInput.OlderPassword, user.Password) {
		return errors.New("Mật khẩu cũ không đúng!")
	}
	// Hash mật khẩu mới
	hashNewPassword, err := utils.HashPassword(passwordInput.NewPassword)
	if err != nil {
		return errors.New("Lỗi khi hash mật khẩu mới!")
	}
	// Cập nhật mật khẩu
	update := bson.M{"$set": bson.M{"password": hashNewPassword}}
	_, err = a.usercollection.UpdateOne(a.ctx, filter, update)
	if err != nil {
		return errors.New("Thay đổi mật khẩu thất bại!")
	}
	return nil
}
func (a *UserImplementService) SaveOTPForUser(userId primitive.ObjectID, otpHash string) error {
	updateData := Models.OTPModel{
		Id:        primitive.NewObjectID(),
		UserId:    userId,
		OTPCode:   otpHash,
		ExpiresAt: time.Now().UTC().Add(5 * time.Minute),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_, err := a.otpcollection.InsertOne(a.ctx, updateData)
	if err != nil {
		return err
	}
	return nil
}
func (a *UserImplementService) FindEmail(email string) (*Models.UserModel, error) {
	var user *Models.UserModel
	query := bson.D{bson.E{Key: "email", Value: email}}
	err := a.usercollection.FindOne(a.ctx, query).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (a *UserImplementService) VerifyOTP(email, otp string) (bool, error) {
	var user *Models.UserModel
	var otpRes *Models.OTPRes
	user, err := a.FindEmail(email)
	if err != nil {
		return false, err
	}
	filter := bson.M{"user_id": user.Id}
	err = a.otpcollection.FindOne(a.ctx, filter).Decode(&otpRes)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, errors.New("Đã xảy ra lỗi hệ thống!")
		}
		return false, errors.New("Email không tồn tại!")
	}
	timeNowUTC := time.Now().UTC()
	if timeNowUTC.After(otpRes.ExpiresAt) {
		return false, errors.New("OTP đã hết hạn!")
	}
	if !utils.CheckOTPHash(otp, otpRes.OTPCode) {
		return false, errors.New("OTP không đúng!")
	}
	return true, nil
}
func (a *UserImplementService) ResendOTP(userId primitive.ObjectID, otpHash string) (bool, error) {
	var userOTP *Models.OTPModel
	//var otpRes *Models.OTPRes
	filter := bson.M{"user_id": userId}
	err := a.otpcollection.FindOne(a.ctx, filter).Decode(&userOTP)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, errors.New("Chưa thực hiện gửi OTP!")
		}
		return false, errors.New("Chưa thực hiện gửi OTP!")
	}
	//newOTPHash, err := utils.HashPassword(otp)
	if err != nil {
		return false, errors.New("Đã xảy ra lỗi hệ thống khi mã hóa OTP!")
	}

	// Bước 4: Cập nhật hoặc tạo mới OTP trong cơ sở dữ liệu
	updateOTPData := bson.D{{
		"$set", bson.D{
			{"user_id", userId},
			{"otp_code", otpHash},
			{"expires_at", time.Now().UTC().Add(5 * time.Minute)},
			{"created_at", time.Now().UTC()},
			{"updated_at", time.Now().UTC()},
		},
	}}
	_, err = a.otpcollection.UpdateOne(a.ctx, filter, updateOTPData, options.Update().SetUpsert(true))
	if err != nil {
		return false, errors.New("Không thể cập nhật OTP trong cơ sở dữ liệu!")
	}
	return true, nil
}
func (a *UserImplementService) ForgotPasswordByOTP(forgotPasswordInput *Models.ForgotPasswordInput) (bool, error) {
	var user *Models.UserModel
	var otpRes *Models.OTPRes
	user, err := a.FindEmail(forgotPasswordInput.Email)
	if err != nil {
		return false, errors.New("Không tìm thấy user!")
	}
	filter := bson.D{
		{"user_id", user.Id},
		{"otp_code", forgotPasswordInput.OtpToken},
	}
	err = a.otpcollection.FindOne(a.ctx, filter).Decode(&otpRes)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, errors.New("Đã xảy ra lỗi hệ thống!")
		}
		return false, errors.New("Email không tồn tại!")
	}
	// kiểm tra thời gian hết hạn của token
	timeNowUTC := time.Now().UTC()
	if timeNowUTC.After(otpRes.ExpiresAt) {
		return false, errors.New("OTP đã hết hạn!")
	}
	if forgotPasswordInput.OtpToken != otpRes.OTPCode {
		return false, errors.New("OTP không đúng!")
	}
	// hash new password
	hashPassword, _ := utils.HashPassword(forgotPasswordInput.NewPassword)
	// update password
	filterUpdate := bson.M{"_id": user.Id}
	update := bson.D{
		{"$set", bson.D{{"password", hashPassword}}},
	}
	_, err = a.usercollection.UpdateOne(a.ctx, filterUpdate, update)
	if err != nil {
		return false, errors.New("Lấy lại mật khẩu thất bại!")
	}
	_, errDelete := a.otpcollection.DeleteOne(a.ctx, filter)
	if errDelete != nil {
		return false, errors.New("Lỗi khi xóa OTP!")
	}
	return true, nil
}

func (a *UserImplementService) SearchUser(name string, page, limit int) ([]Models.UserModel, int, error) {
	skip := (page - 1) * limit
	totalCount, err := a.usercollection.CountDocuments(a.ctx, bson.D{{"$text", bson.D{{"$search", name}}}})
	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"$text", bson.D{{"$search", name}}}}}},
		bson.D{{"$skip", skip}},
		bson.D{{"$limit", limit}},
	}
	var userRes []Models.UserModel
	cursor, err := a.usercollection.Aggregate(a.ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(a.ctx)
	if err = cursor.All(a.ctx, &userRes); err != nil {
		return nil, 0, err
	}
	return userRes, int(totalCount), nil
}
func (a *UserImplementService) GetAllUserRoleIsStudent(page, limit int) ([]Models.UserModel, int, error) {
	skip := limit * (page - 1)
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cur, err := a.usercollection.Find(a.ctx, bson.M{"role_type": constant.STUDENT}, opts)
	total, err := a.usercollection.CountDocuments(a.ctx, bson.M{"role_type": constant.STUDENT})
	defer cur.Close(a.ctx)
	var students []Models.UserModel
	for cur.Next(a.ctx) {
		var student Models.UserModel
		err := cur.Decode(&student)
		if err != nil {
			return nil, 0, err
		}
		students = append(students, student)
	}
	if err := cur.Err(); err != nil {
		return nil, 0, err
	}
	return students, int(total), err
}
func (a *UserImplementService) GetAllUserRoleIsTeacher(page, limit int) ([]Models.UserModel, int, error) {
	skip := limit * (page - 1)
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cur, err := a.usercollection.Find(a.ctx, bson.M{"role_type": constant.TEACHER}, opts)
	total, err := a.usercollection.CountDocuments(a.ctx, bson.M{"role_type": constant.TEACHER})
	defer cur.Close(a.ctx)
	var teachers []Models.UserModel
	for cur.Next(a.ctx) {
		var teacher Models.UserModel
		err := cur.Decode(&teacher)
		if err != nil {
			return nil, 0, err
		}
		teachers = append(teachers, teacher)
	}
	if err := cur.Err(); err != nil {
		return nil, 0, err
	}
	return teachers, int(total), err
}
