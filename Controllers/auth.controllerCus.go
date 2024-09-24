package Controllers

import (
	"gin-gonic-gom/Collections"
	"gin-gonic-gom/Common"
	"gin-gonic-gom/Models"
	"gin-gonic-gom/config"
	"gin-gonic-gom/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
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
