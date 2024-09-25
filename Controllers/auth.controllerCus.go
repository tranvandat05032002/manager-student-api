package Controllers

import (
	"fmt"
	"gin-gonic-gom/Collections"
	"gin-gonic-gom/Common"
	"gin-gonic-gom/Models"
	"gin-gonic-gom/config"
	"gin-gonic-gom/constant"
	"gin-gonic-gom/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
	"time"
)

type LoginForm struct {
	Email    string `json:"email" bson:"email" binding:"required,email"`
	Password string `json:"password" bson:"password" binding:"required,min=6,max=30"`
}
type FindEmailForm struct {
	Email string `json:"email" bson:"email" binding:"required,email"`
}
type OTPForm struct {
	Email   string `json:"email" bson:"-"`
	OTPCode string `json:"otp_code" bson:"-"`
}
type ForgotPasswordForm struct {
	Email           string `json:"email" bson:"-" binding:"required,email"`
	OtpToken        string `json:"otp_token" bson:"-" binding:"required"`
	NewPassword     string `json:"new_password,omitempty" bson:"-" binding:"required,min=6,max=30"`
	ConfirmPassword string `json:"confirm_password,omitempty" bson:"-" binding:"required,min=6,max=30"`
}
type ChangePasswordForm struct {
	OlderPassword      string `json:"older_password" binding:"required,min=6,max=30"`
	NewPassword        string `json:"new_password,omitempty" binding:"required,min=6,max=30"`
	ConfirmNewPassword string `json:"confirm_new_password,omitempty" binding:"required,min=6,max=30"`
}
type UserUpdate struct {
	MajorId        primitive.ObjectID `json:"major_id"`
	Email          string             `json:"email" bson:"email"`
	Role           string             `json:"role_type" bson:"role_type"`
	Phone          string             `json:"phone" bson:"phone"`
	Name           string             `json:"name" bson:"name"`
	Avatar         string             `json:"avatar" bson:"avatar"`
	Gender         int                `json:"gender" bson:"gender"`
	Department     string             `json:"department" bson:"department"`
	DateOfBirth    time.Time          `json:"date_of_birth" bson:"dateOfBirth"`
	EnrollmentDate time.Time          `json:"enrollment_date" bson:"enrollmentDate"`
	HireDate       time.Time          `json:"hire_date" bson:"hireDate"`
	Address        string             `json:"address" bson:"address"`
}

func PostUser(c *gin.Context) {
	entry := Collections.UserModel{}
	var (
		DB  = config.GetMongoDB()
		err error
		//Other config
		//.......
	)
	if err = c.ShouldBindBodyWith(&entry, binding.JSON); err != nil {
		//Logger
		//Response
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	switch entry.Role {
	case "student":
		entry.Role = "student"
	case "teacher":
		entry.Role = "teacher"
	case "admin":
		entry.Role = "admin"
	default:
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorRoleMessage, "")
		return
	}
	res, errCheckTerm := entry.CheckExit(DB, entry.Email, entry.Phone)
	if errCheckTerm != nil {
		Common.NewErrorResponse(c, http.StatusInternalServerError, "Lỗi hệ thống! ", nil)
		return
	}
	if res {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Email hoặc số điện thoại đã tồn tại!", nil)
		return
	}
	password := entry.Password
	passwordHash, _ := utils.HashPassword(password)
	entry.Password = passwordHash
	err = entry.Create(DB)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Lỗi hệ thống!", nil)
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Thêm dữ liệu thành công!", nil))
}
func SignIn(c *gin.Context) {
	var request LoginForm
	device := c.Request.Header.Get("User-Agent")
	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		//Logger
		//Response
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}

	var (
		err        error
		entry      = Collections.UserModel{}
		tokenEntry = Collections.TokenModel{}
		DB         = config.GetMongoDB()
	)
	res, err := entry.CheckIsValid(DB, request.Email, request.Password)
	if err != nil {
		Common.NewInternal()
		return
	}
	if !res {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Email hoặc mật khẩu không hợp lệ", nil)
		return
	}
	// Kiem tra tai khoan co dang bi khoa hoac bi xoa khong
	isDelete := entry.CheckIsLocked(DB, request.Email)
	if isDelete {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Tài khoản đang bị khóa hoặc không tồn tại trong hệ thống!", nil)
		return
	}
	user, _ := entry.FindByEmail(DB, request.Email)
	//Tao token
	accessToken, _ := utils.GenerateAccessToken(user.Id, user.Role, os.Getenv("ACCESS_TOKEN_SECRET"))
	refreshToken, _ := utils.GenerateRefreshToken(user.Id, user.Role, os.Getenv("ACCESS_TOKEN_SECRET"))
	c.SetCookie("refresh_token", accessToken, 3600, "/", c.Request.Host, true, true)
	c.SetCookie("access_token", refreshToken, 259200, "/", c.Request.Host, true, true)
	token := Models.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	// Cap nhat thiet bi
	filterDevice := bson.D{{"device", device}}
	expToken := utils.Now().Add(72 * time.Hour)
	updateToken := bson.D{{"$set", bson.D{
		{"user_id", user.Id},
		{"refresh_token", refreshToken},
		{"access_token", accessToken},
		{"exp", utils.ConvertToVietnamTime(expToken)},
		{"device", device},
		{"created_at", time.Now().UTC()},
		{"updated_at", time.Now().UTC()},
	}}}
	err = tokenEntry.FindAndUpdate(DB, filterDevice, updateToken)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusNotFound, Common.ErrorFindUser, nil)
		return
	}
	err = tokenEntry.CheckAndDeleteDevice(DB, user.Id)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusInternalServerError, Common.ErrorInternetServer, nil)
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, Common.SuccessLogin, token))
}
func GetMe(c *gin.Context) {
	var (
		DB    = config.GetMongoDB()
		entry Collections.UserModel
	)
	userId, exists := c.MustGet("userId").(string)
	role, exists := c.Get("role")
	if !exists {
		Common.NewErrorResponse(c, http.StatusUnauthorized, Common.ErrorRoleOrUserIDMessage, "")
		return
	}
	filter := bson.M{
		"_id": utils.ConvertStringToObjectId(userId),
	}
	switch role {
	case "student":
		opts := options.FindOne().SetProjection(bson.D{{"password", 0}, {"hire_date", 0}, {"department", 0}})
		_, _ = entry.FindOne(DB, filter, opts)
	case "teacher":
		teacherProjection := bson.D{{"password", 0}, {"address", 0}, {"major_id", 0}, {"enrollment_date", 0}}
		opts := options.FindOne().SetProjection(teacherProjection)
		_, _ = entry.FindOne(DB, filter, opts)
	case "admin":
		adminProjection := bson.D{{"password", 0}, {"address", 0}, {"major_id", 0}, {"enrollment_date", 0}, {"hire_date", 0}, {"department", 0}}
		opts := options.FindOne().SetProjection(adminProjection)
		_, _ = entry.FindOne(DB, filter, opts)
	default:
		Common.NewErrorResponse(c, http.StatusUnauthorized, "không thể lấy dữ liệu", "")
	}
	//res, _ := userController.UserService.GetMe(userId.(string), role.(string))
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, Common.SuccessGetMe, entry))
}
func SignOut(c *gin.Context) {
	var (
		DB         = config.GetMongoDB()
		tokenEntry Collections.TokenModel
	)
	var token = c.Request.Header.Get("Authorization")
	device := c.Request.Header.Get("User-Agent")
	if token == "" {
		cookie, _ := c.Request.Cookie("access_token")
		token = cookie.String()
		if len(token) > 6 {
			token = token[6:]
		}
	}
	filter := bson.M{"device": device}

	err := tokenEntry.DeleteOne(DB, filter)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Đăng xuất thất bại", "")
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, Common.SuccessLogout, nil))
}

func FindEmail(c *gin.Context) {
	var request FindEmailForm
	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		//Logger
		//Response
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	var (
		DB       = config.GetMongoDB()
		err      error
		entry    Collections.UserModel
		otpEntry Collections.OTPModel
	)
	_, err = entry.FindByEmail(DB, request.Email)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusNotFound, Common.ErrorEmailNotFound, nil)
	}
	otp, err := utils.GeneratorOTP(6)
	otpHash, _ := utils.HashPassword(otp)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorCreateOTP, err.Error())
		return
	}
	errOTP := otpEntry.SaveOTP(DB, entry.Id, otpHash)
	if errOTP != nil {
		Common.NewErrorResponse(c, http.StatusInternalServerError, Common.ErrorInternetServer, errOTP.Error())
		return
	}
	_ = utils.SendSecretCodeToEmail(entry.Email, otp, otpHash)
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "tìm thấy user theo email success!!!", otp))
}
func VerifyOTP(c *gin.Context) {
	var request OTPForm
	if err := c.ShouldBindJSON(&request); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	var (
		DB       = config.GetMongoDB()
		err      error
		otpEntry Collections.OTPModel
	)
	_, err = otpEntry.VerifyOTP(DB, request.Email, request.OTPCode)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorOTP, err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Xác thực OTP thành công!", nil))
}
func ResendOTP(c *gin.Context) {
	var request FindEmailForm
	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		//Logger
		//Response
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	var (
		DB        = config.GetMongoDB()
		err       error
		userEntry Collections.UserModel
		otpEntry  Collections.OTPModel
	)
	_, err = userEntry.FindByEmail(DB, request.Email)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorEmailNotFound, err.Error())
		return
	}
	// Send OTP to email
	otp, err := utils.GeneratorOTP(6)
	otpHash, _ := utils.HashOTP(otp)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorCreateOTP, err.Error())
		return
	}
	_, errOTP := otpEntry.ResendOTP(DB, userEntry.Id, otpHash)
	if errOTP != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorResendOTP, nil)
		return
	}
	_ = utils.SendSecretCodeToEmail(userEntry.Email, otp, otpHash)
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Đã gửi lại mã OTP thành công!!!", otp))
}
func ResetPassword(c *gin.Context) {
	var request ForgotPasswordForm
	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		//Logger
		//Response
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	var (
		DB        = config.GetMongoDB()
		err       error
		userEntry Collections.UserModel
	)
	if request.NewPassword != request.ConfirmPassword {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Mật khẩu không khớp!", "")
		return
	}
	err = userEntry.ResetPasswordByOTP(DB, request.Email, request.OtpToken, request.NewPassword)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorUpdatePassword, nil)
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Reset Password Success!!", nil))
}

func ChangePassword(c *gin.Context) {
	var request ChangePasswordForm
	userId, exists := c.Get("userId")
	if !exists {
		Common.NewErrorResponse(c, http.StatusNotFound, Common.ErrorFindUser, "")
		return
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	if request.NewPassword != request.ConfirmNewPassword {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Mật khẩu không khớp!", nil)
		return
	}
	var (
		DB        = config.GetMongoDB()
		userEntry Collections.UserModel
	)
	filter := bson.M{
		"_id": utils.ConvertStringToObjectId(userId.(string)),
	}

	_, err := userEntry.FindOne(DB, filter, nil)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Không tìm thấy người dùng!", nil)
		return
	}
	if !utils.CheckPasswordHash(request.OlderPassword, userEntry.Password) {
		Common.NewErrorResponse(c, http.StatusBadRequest, " Mật khẩu cũ không đúng", nil)
		return
	}
	// Hash lai mat khau moi
	hashNewPassword, err := utils.HashPassword(request.NewPassword)
	// Cập nhật mật khẩu
	update := bson.M{"$set": bson.M{"password": hashNewPassword}}
	err = userEntry.UpdateOne(DB, filter, update)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Thay đổi mật khẩu thất bại!", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Thay đổi mật khẩu thành công", nil))
}
func GetDetailUser(c *gin.Context) {
	userId := c.Param("user_id")
	var (
		DB        = config.GetMongoDB()
		userEntry Collections.UserModel
	)
	filter := bson.M{
		"_id": utils.ConvertStringToObjectId(userId),
	}
	res, err := userEntry.FindOne(DB, filter, options.FindOne())
	if err != nil {
		Common.NewErrorResponse(c, http.StatusNotFound, Common.ErrorFindUser, err.Error())
		return
	}
	c.JSON(
		http.StatusOK,
		Common.SimpleSuccessResponse(http.StatusOK, "Get user thành công!", res),
	)
}
func UpdateMe(c *gin.Context) {
	request := UserUpdate{}
	userId, exists := c.Get("userId")
	if !exists {
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorFindUser, "")
		return
	}
	var (
		userEntry Collections.UserModel
		DB        = config.GetMongoDB()
		err       error
		//Other config
		//.......
	)
	if err = c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		//Logger
		//Response
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	filter := bson.M{
		"_id": utils.ConvertStringToObjectId(userId.(string)),
	}
	err = userEntry.UpdateMe(DB, filter, request)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Cập nhật thất bại!", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Cập nhật thành công!", nil))
}
func GetAll(c *gin.Context) {
	var (
		userEntry Collections.UserModel
		DB        = config.GetMongoDB()
		err       error
		//Other config
		//.......
	)
	page, limit, skip := utils.Pagination(c)
	filter := bson.D{
		{"role_type", bson.D{{"$ne", "admin"}}},
	}
	total, _ := userEntry.Count(DB, filter)
	res, err := userEntry.Find(DB, limit, skip)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Không thể lấy thông tin!", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Lấy danh sách users thành công", res, int(total), page, limit))
}
func SearchUser(c *gin.Context) {
	var (
		userEntry Collections.UserModel
		DB        = config.GetMongoDB()
		err       error
		//Other config
		//.......
	)

	query := c.Query("name")
	page, limit, skip := utils.Pagination(c)
	res, total, err := userEntry.Search(DB, query, skip, limit)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Tìm kiếm xảy ra lỗi!", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Tìm kiếm môn học thành công!!", res, total, page, limit))
}

func GetStudent(c *gin.Context) {
	var (
		userEntry Collections.UserModel
		DB        = config.GetMongoDB()
		err       error
		//Other config
		//.......
	)

	page, limit, skip := utils.Pagination(c)
	res, total, err := userEntry.GetByRole(DB, constant.STUDENT, skip, limit)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Không thể lấy thông tin!", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Lấy danh sách sinh viên thành công", res, total, page, limit))
}
func GetTeacher(c *gin.Context) {
	var (
		userEntry Collections.UserModel
		DB        = config.GetMongoDB()
		err       error
		//Other config
		//.......
	)

	page, limit, skip := utils.Pagination(c)
	res, total, err := userEntry.GetByRole(DB, constant.TEACHER, skip, limit)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Không thể lấy thông tin!", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Lấy danh sách giáo viên thành công", res, total, page, limit))
}
func UpdateUser(c *gin.Context) {
	request := UserUpdate{}
	id := c.Param("id")
	var (
		userEntry Collections.UserModel
		DB        = config.GetMongoDB()
		err       error
		//Other config
		//.......
	)
	if err = c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		//Logger
		//Response
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	filter := bson.M{
		"_id": utils.ConvertStringToObjectId(id),
	}
	err = userEntry.Update(DB, filter, request, request.MajorId)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Cập nhật thất bại!", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Cập nhật thành công!", nil))
}
func DeleteUser(c *gin.Context) {
	var (
		DB        = config.GetMongoDB()
		err       error
		userEntry Collections.UserModel
	)
	userId := c.Param("user_id")
	filter := bson.M{
		"_id": utils.ConvertStringToObjectId(userId),
	}
	err = userEntry.DeleteBackUp(DB, filter)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorDeleteUser, err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Xóa thành công account", nil))
}
func GetUserDepending(c *gin.Context) {
	var (
		DB        = config.GetMongoDB()
		err       error
		userEntry Collections.UserModel
	)

	page, limit, skip := utils.Pagination(c)
	filter := bson.M{"depending_delete": constant.ISDEPENDING}
	res, total, err := userEntry.FindByStatus(DB, filter, skip, limit)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Không thể lấy thông tin!", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Lấy danh sách sinh viên đang chờ xóa thành công", res, total, page, limit))
}

func RestoreUser(c *gin.Context) {
	userId := c.Param("id")
	var (
		DB        = config.GetMongoDB()
		err       error
		userEntry Collections.UserModel
	)
	filter := bson.M{
		"_id": utils.ConvertStringToObjectId(userId),
	}
	err = utils.DelCache(utils.ConvertStringToObjectId(userId).Hex())
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Xảy ra lỗi khi xóa cache", err.Error())
		return
	}
	err = userEntry.RestoreUser(DB, filter)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorDeleteUser, err.Error())
		return
	}
	message := fmt.Sprintf("khôi phục tài khoản %s", userId)
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, message, nil))
}
