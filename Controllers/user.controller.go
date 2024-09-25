package Controllers

import (
	"fmt"
	"gin-gonic-gom/Common"
	"gin-gonic-gom/Middlewares"
	"gin-gonic-gom/Models"
	"gin-gonic-gom/Services/user"
	"gin-gonic-gom/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
// @Success 200 {object} Common.Response{status=int,message=string,data=interface{}} "Success response"
// @Failure 400 {object} Common.ErrorResponse "Bad Request response"
// @Router /user/create [post]
func (userController *UserController) GetAccount(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	res, err := userController.UserService.GetAccount(userId)
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusNotFound, Common.ErrorFindUser, err.Error())
		return
	}
	ctx.JSON(
		http.StatusOK,
		Common.SimpleSuccessResponse(http.StatusOK, "Get user thành công!", res),
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
	res, total, err := userController.UserService.GetAll(page, limit)
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Lấy danh sách users thành công", res, total, page, limit))
}
func (userController *UserController) UpdateMe(ctx *gin.Context) {
	var userEntity Models.UserUpdate
	if err := ctx.ShouldBindJSON(&userEntity); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	userId, exists := ctx.Get("userId")
	if !exists {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorFindUser, "")

		return
	}
	res, err := userController.UserService.UpdateMe(userId.(string), &userEntity)
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorUpdateUser, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, Common.SuccessUpdateData, res))
}
func (userController *UserController) UpdateUser(ctx *gin.Context) {
	var userEntity Models.AccountUpdate
	id := ctx.Param("id")
	if err := ctx.ShouldBindJSON(&userEntity); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	res, err := userController.UserService.UpdateUser(&userEntity, utils.ConvertStringToObjectId(id))
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorUpdateUser, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, Common.SuccessUpdateData, res))
}
func (userController *UserController) DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("user_id")
	err := userController.UserService.DeleteUser(userId)
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorDeleteUser, err.Error())
		return
	}
	message := fmt.Sprintf("Xóa thành công account")
	ctx.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, message, nil))
}
func (userController *UserController) ChangePassword(ctx *gin.Context) {
	var passwordInput Models.ChangePasswordInput
	userId, exists := ctx.Get("userId")
	if !exists {
		Common.NewErrorResponse(ctx, http.StatusNotFound, Common.ErrorFindUser, "")
		return
	}
	if err := ctx.ShouldBindJSON(&passwordInput); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	err := userController.UserService.ChangePassword(userId.(string), &passwordInput)
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorChangePassword, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Thay đổi mật khẩu thành công", nil))
}
func (userController *UserController) FindEmail(ctx *gin.Context) {
	var emailInput Models.FindEmailInput
	if err := ctx.ShouldBindJSON(&emailInput); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	userByEmail, err := userController.UserService.FindEmail(emailInput.Email)
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorEmailNotFound, err.Error())
		return
	}
	// Send OTP to email
	otp, err := utils.GeneratorOTP(6)
	otpHash, _ := utils.HashPassword(otp)
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorCreateOTP, err.Error())
		return
	}
	errOTP := userController.UserService.SaveOTPForUser(userByEmail.Id, otpHash)
	if errOTP != nil {
		Common.NewErrorResponse(ctx, http.StatusInternalServerError, Common.ErrorInternetServer, errOTP.Error())
		return
	}
	_ = utils.SendSecretCodeToEmail(userByEmail.Email, otp, otpHash)
	ctx.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "tìm thấy user theo email success!!!", otp))
}
func (userController *UserController) VerifyOTP(ctx *gin.Context) {
	var OTPReq Models.OTPReq
	if err := ctx.ShouldBindJSON(&OTPReq); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	_, err := userController.UserService.VerifyOTP(OTPReq.Email, OTPReq.OTPCode)
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorOTP, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Xác thực OTP thành công!", nil))
}
func (userController *UserController) ResendOTP(ctx *gin.Context) {
	var emailInput Models.FindEmailInput
	if err := ctx.ShouldBindJSON(&emailInput); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	userByEmail, err := userController.UserService.FindEmail(emailInput.Email)
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorEmailNotFound, err.Error())
		return
	}
	// Send OTP to email
	otp, err := utils.GeneratorOTP(6)
	otpHash, _ := utils.HashOTP(otp)
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorCreateOTP, err.Error())
		return
	}
	_, errOTP := userController.UserService.ResendOTP(userByEmail.Id, otpHash)
	if errOTP != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorResendOTP, errOTP.Error())
		return
	}
	_ = utils.SendSecretCodeToEmail(userByEmail.Email, otp, otpHash)
	ctx.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Đã gửi lại mã OTP thành công!!!", otp))
}
func (userController *UserController) ForgotPasswordByOTP(ctx *gin.Context) {
	var forgotPasswordInput Models.ForgotPasswordInput
	if err := ctx.ShouldBindJSON(&forgotPasswordInput); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	if forgotPasswordInput.ConfirmPassword != forgotPasswordInput.NewPassword {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, "Mật khẩu không khớp!", "")
		return
	}
	_, err := userController.UserService.ForgotPasswordByOTP(&forgotPasswordInput)
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorUpdatePassword, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Reset Password Success!!", forgotPasswordInput))
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
	res, total, _ := userController.UserService.SearchUser(nameQuery, page, limit)
	ctx.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Tìm kiếm user thành công!!", res, total, page, limit))
}
func (userController *UserController) GetAllUserRoleIsStudent(ctx *gin.Context) {
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
	res, total, err := userController.UserService.GetAllUserRoleIsStudent(page, limit)
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, "Không thể lấy thông tin!", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Lấy danh sách sinh viên thành công", res, total, page, limit))
}
func (userController *UserController) GetAllUserRoleIsTeacher(ctx *gin.Context) {
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
	res, total, err := userController.UserService.GetAllUserRoleIsTeacher(page, limit)
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, "Không thể lấy thông tin!", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Lấy danh sách giáo viên thành công", res, total, page, limit))
}
func (userController *UserController) GetListUserDependingDeletion(ctx *gin.Context) {
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
	res, total, err := userController.UserService.GetListUserDependingDeletion(page, limit)
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, "Không thể lấy danh sách!", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Lấy danh sách sinh viên đang chờ xóa thành công", res, total, page, limit))
}
func (userController *UserController) RestoreUser(ctx *gin.Context) {
	userId := ctx.Param("id")
	err := userController.UserService.RestoreUser(utils.ConvertStringToObjectId(userId))
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorDeleteUser, err.Error())
		return
	}
	message := fmt.Sprintf("khôi phục tài khoản %s", userId)
	ctx.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, message, nil))
}
func (authController *UserController) RegisterAuthRoutes(rg *gin.RouterGroup) {
	userRoute := rg.Group("/user") // Client
	{
		userRoute.POST("/forgot-password", authController.ForgotPasswordByOTP)
		userRoute.Use(Middlewares.AuthValidationBearerMiddleware)
		{
			userRoute.GET("/:user_id", authController.GetAccount)
			userRoute.PATCH("/me/update", authController.UpdateMe)
			userRoute.PUT("/change-password", authController.ChangePassword)
		}
	}
	// Admin
	adminRoute := rg.Group("/admin")
	adminRoute.Use(Middlewares.AuthValidationBearerMiddleware)
	adminRoute.Use(Middlewares.RoleMiddleware("admin"))
	{
		adminRoute.GET("/all", authController.GetAll)
		adminRoute.GET("/user/search", authController.SearchUser)
		adminRoute.GET("/user/student/all", authController.GetAllUserRoleIsStudent)
		adminRoute.GET("/user/teacher/all", authController.GetAllUserRoleIsTeacher)
		adminRoute.PATCH("/update/:id", authController.UpdateUser)
		adminRoute.DELETE("/delete/:user_id", authController.DeleteUser)
		adminRoute.GET("/user/pending-deletion", authController.GetListUserDependingDeletion)
		adminRoute.PATCH("/user/restore/:id", authController.RestoreUser)
	}
}
