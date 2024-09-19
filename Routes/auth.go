package Routes

import (
	"gin-gonic-gom/Controllers"
)

var (
	authController Controllers.UserController
)

//func RegisterAuthRoutes(rg *gin.RouterGroup) {
//	userRoute := rg.Group("/user") // Client
//	{
//		userRoute.POST("/create", authController.CreateUser)
//		userRoute.POST("/login", authController.LoginUser)
//		userRoute.POST("/find-email", authController.FindEmail)
//		userRoute.POST("/otp/verify-otp", authController.VerifyOTP)
//		userRoute.POST("/otp/resend-otp", authController.ResendOTP)
//		userRoute.POST("/forgot-password", authController.ForgotPasswordByOTP)
//		userRoute.Use(Middlewares.AuthValidationBearerMiddleware)
//		{
//			userRoute.GET("/me", authController.GetMe)
//			userRoute.POST("/logout", authController.Logout)
//			userRoute.GET("/:user_id", authController.GetAccount)
//			userRoute.PATCH("/me/update", authController.UpdateMe)
//			userRoute.PUT("/change-password", authController.ChangePassword)
//		}
//	}
//	// Admin
//	adminRoute := rg.Group("/admin")
//	adminRoute.Use(Middlewares.AuthValidationBearerMiddleware)
//	adminRoute.Use(Middlewares.RoleMiddleware("admin"))
//	{
//		adminRoute.GET("/all", authController.GetAll)
//		adminRoute.GET("/user/search", authController.SearchUser)
//		adminRoute.GET("/user/student/all", authController.GetAllUserRoleIsStudent)
//		adminRoute.GET("/user/teacher/all", authController.GetAllUserRoleIsTeacher)
//		adminRoute.PATCH("/update/:id", authController.UpdateUser)
//		adminRoute.DELETE("/delete/:user_id", authController.DeleteUser)
//		adminRoute.GET("/user/pending-deletion", authController.GetListUserDependingDeletion)
//		adminRoute.PATCH("/user/restore/:id", authController.RestoreUser)
//	}
//}
