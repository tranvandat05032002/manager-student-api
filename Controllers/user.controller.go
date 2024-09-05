package Controllers

import (
	"fmt"
	"gin-gonic-gom/Middlewares"
	"gin-gonic-gom/Models"
	"gin-gonic-gom/Services/user"
	"gin-gonic-gom/common"
	"gin-gonic-gom/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService user.UserService
}

func New(userService user.UserService) UserController {
	return UserController{
		UserService: userService,
	}
}

// CreateUser godoc
// @Description Tạo mới một người dùng
// @Tags User
// @Accept  json
// @Produce  json
// @Param createUserInput body Models.CreateUserInput true "Create User Input"
// @Success 200 {object} common.Response{status=int,message=string,data=interface{}} "Success response"
// @Failure 400 {object} common.ErrorResponse "Bad Request response"
// @Router /user/create [post]
func (userController *UserController) CreateUser(ctx *gin.Context) {
	var user Models.CreateUserInput
	if err := ctx.ShouldBindJSON(&user); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	switch user.Role {
	case "student":
		user.Role = "student"
	case "teacher":
		user.Role = "teacher"
	case "admin":
		user.Role = "admin"
	default:
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorRoleMessage, "")
		return
	}
	//Check exist email in DB
	result, errCheckMail := userController.UserService.CheckExistEmail(user.Email)
	if result == true {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorEmailExistMessage, "Email already exists")
		return
	}

	if errCheckMail != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorEmailExistMessage, errCheckMail.Error())
		return
	}
	var password string = user.Password
	passwordHash, _ := utils.HashPassword(password)
	user.Password = passwordHash
	// add user to DB
	err := userController.UserService.CreateUser(&user)

	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorAddDataMessage, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, common.SuccessAddDataMessage, nil))
}
func (useController *UserController) LoginUser(ctx *gin.Context) {
	var user Models.AuthInput
	if err := ctx.ShouldBindJSON(&user); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	useController.UserService.LoginUser(&user, ctx)
}
func (userController *UserController) Logout(ctx *gin.Context) {
	userId, _ := ctx.MustGet("userId").(string)
	deviced, _ := ctx.MustGet("deviced").(string)
	userController.UserService.Logout(deviced, utils.ConvertStringToObjectId(userId))
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, common.SuccessLogout, nil))
}
func (userController *UserController) GetMe(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	role, exists := ctx.Get("role")
	if !exists {
		common.NewErrorResponse(ctx, http.StatusUnauthorized, common.ErrorRoleOrUserIDMessage, "")
		return
	}
	user, _ := userController.UserService.GetMe(userId.(string), role.(string))
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, common.SuccessGetMe, user))
}
func (userController *UserController) GetAccount(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	account, err := userController.UserService.GetAccount(userId)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusNotFound, common.ErrorFindUser, err.Error())
		return
	}
	ctx.JSON(
		http.StatusOK,
		common.SimpleSuccessResponse(http.StatusOK, "Get user thành công!", account),
	)
}
func (userController *UserController) GetAll(ctx *gin.Context) {
	limitStr := ctx.Query("limit")
	pageStr := ctx.Query("page")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	users, total, err := userController.UserService.GetAll(page, limit)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.NewSuccessResponse(http.StatusOK, "Lấy danh sách users thành công", users, total, page, limit))
}
func (userController *UserController) UpdateMe(ctx *gin.Context) {
	var UserUpdate Models.UserUpdate
	if err := ctx.ShouldBindJSON(&UserUpdate); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	userId, exists := ctx.Get("userId")
	if !exists {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorFindUser, "")

		return
	}
	user, err := userController.UserService.UpdateMe(userId.(string), &UserUpdate)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorUpdateUser, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, common.SuccessUpdateData, user))
}
func (userController *UserController) UpdateUser(ctx *gin.Context) {
	var AccountUpdate Models.AccountUpdate
	id := ctx.Param("id")
	if err := ctx.ShouldBindJSON(&AccountUpdate); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	user, err := userController.UserService.UpdateUser(&AccountUpdate, utils.ConvertStringToObjectId(id))
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorUpdateUser, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, common.SuccessUpdateData, user))
}
func (userController *UserController) DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	res, err := userController.UserService.DeleteUser(userId)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorDeleteUser, err.Error())
		return
	}
	if res < 1 {
		common.NewErrorResponse(ctx, http.StatusNotFound, common.ErrorFindUser, err.Error())
		return
	}
	message := fmt.Sprintf("Xóa thành công %d account", res)
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, message, nil))
}
func (userController *UserController) ChangePassword(ctx *gin.Context) {
	var PasswordInput Models.ChangePasswordInput
	userId, exists := ctx.Get("userId")
	if !exists {
		common.NewErrorResponse(ctx, http.StatusNotFound, common.ErrorFindUser, "")
		return
	}
	if err := ctx.ShouldBindJSON(&PasswordInput); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	err := userController.UserService.ChangePassword(userId.(string), &PasswordInput)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorChangePassword, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "Thay đổi mật khẩu thành công", nil))
}
func (userController *UserController) FindEmail(ctx *gin.Context) {
	var EmailInput Models.FindEmailInput
	if err := ctx.ShouldBindJSON(&EmailInput); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	userByEmail, err := userController.UserService.FindEmail(EmailInput.Email)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorEmailNotFound, err.Error())
		return
	}
	// Send OTP to email
	otp, err := utils.GeneratorOTP(6)
	otpHash, _ := utils.HashPassword(otp)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorCreateOTP, err.Error())
		return
	}
	errOTP := userController.UserService.SaveOTPForUser(userByEmail.Id, otpHash)
	if errOTP != nil {
		common.NewErrorResponse(ctx, http.StatusInternalServerError, common.ErrorInternetServer, errOTP.Error())
		return
	}
	_ = utils.SendSecretCodeToEmail(userByEmail.Email, otp, otpHash)
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "tìm thấy user theo email success!!!", otp))
}
func (usercontroller *UserController) VerifyOTP(ctx *gin.Context) {
	var OTPReq Models.OTPReq
	if err := ctx.ShouldBindJSON(&OTPReq); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	_, err := usercontroller.UserService.VerifyOTP(OTPReq.Email, OTPReq.OTPCode)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorOTP, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "Xác thực OTP thành công!", nil))
}
func (usercontroller *UserController) ResendOTP(ctx *gin.Context) {
	var EmailInput Models.FindEmailInput
	if err := ctx.ShouldBindJSON(&EmailInput); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	userByEmail, err := usercontroller.UserService.FindEmail(EmailInput.Email)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorEmailNotFound, err.Error())
		return
	}
	// Send OTP to email
	otp, err := utils.GeneratorOTP(6)
	otpHash, _ := utils.HashOTP(otp)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorCreateOTP, err.Error())
		return
	}
	_, errOTP := usercontroller.UserService.ResendOTP(userByEmail.Id, otpHash)
	if errOTP != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorResendOTP, errOTP.Error())
		return
	}
	_ = utils.SendSecretCodeToEmail(userByEmail.Email, otp, otpHash)
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "Đã gửi lại mã OTP thành công!!!", otp))
}
func (userController *UserController) ForgotPasswordByOTP(ctx *gin.Context) {
	var ForgotPasswordInput Models.ForgotPasswordInput
	if err := ctx.ShouldBindJSON(&ForgotPasswordInput); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	if ForgotPasswordInput.ConfirmPassword != ForgotPasswordInput.NewPassword {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Mật khẩu không khớp!", "")
		return
	}
	_, err := userController.UserService.ForgotPasswordByOTP(&ForgotPasswordInput)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorUpdatePassword, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "Reset Password Success!!", ForgotPasswordInput))
}
func (userController *UserController) SearchUser(ctx *gin.Context) {
	nameQuery := ctx.Query("name")
	limitStr := ctx.Query("limit")
	pageStr := ctx.Query("page")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	name := string(nameQuery)
	users, total, _ := userController.UserService.SearchUser(name, page, limit)
	ctx.JSON(http.StatusOK, common.NewSuccessResponse(http.StatusOK, "Tìm kiếm user thành công!!", users, total, page, limit))
}
func (userController *UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
	userroute := rg.Group("/user") // Client
	{
		userroute.POST("/create", userController.CreateUser)
		userroute.POST("/login", userController.LoginUser)
		userroute.POST("/find-email", userController.FindEmail)
		userroute.POST("/otp/verify-otp", userController.VerifyOTP)
		userroute.POST("/otp/resend-otp", userController.ResendOTP)
		userroute.POST("/forgot-password", userController.ForgotPasswordByOTP)
		userroute.Use(Middlewares.AuthValidationBearerMiddleware)
		{
			userroute.GET("/me", userController.GetMe)
			userroute.POST("/logout", userController.Logout)
			userroute.GET("/:user_id", userController.GetAccount)
			userroute.PATCH("/me/update", userController.UpdateMe)
			userroute.PUT("/change-password", userController.ChangePassword)
		}
	}
	// Admin
	adminroute := rg.Group("/admin")
	adminroute.Use(Middlewares.AuthValidationBearerMiddleware)
	adminroute.Use(Middlewares.RoleMiddleware("admin"))
	{
		adminroute.GET("/all", userController.GetAll)
		adminroute.GET("/user/search", userController.SearchUser)
		adminroute.PATCH("/update/:id", userController.UpdateUser)
		adminroute.DELETE("/delete/:user_id", userController.DeleteUser)
	}
}
