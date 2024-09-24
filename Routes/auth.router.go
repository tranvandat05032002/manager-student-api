package Routes

import (
	"gin-gonic-gom/Controllers"
	"gin-gonic-gom/Middlewares"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(auth *gin.RouterGroup) {
	authRoute := auth.Group("/user") // Client
	{
		authRoute.POST("/create", Controllers.PostUser)
		authRoute.POST("/login", Controllers.SignIn)
		authRoute.POST("/find-email", Controllers.FindEmail)
		authRoute.POST("/otp/verify-otp", Controllers.VerifyOTP)
		authRoute.POST("/otp/resend-otp", Controllers.ResendOTP)
		//userRoute.POST("/forgot-password", authController.ForgotPasswordByOTP)
		authRoute.Use(Middlewares.AuthValidationBearerMiddleware)
		{
			authRoute.GET("/me", Controllers.GetMe)
			authRoute.POST("/logout", Controllers.SignOut)
			//userRoute.GET("/:user_id", authController.GetAccount)
			//userRoute.PATCH("/me/update", authController.UpdateMe)
			//userRoute.PUT("/change-password", authController.ChangePassword)
		}
	}
	// Admin
	//adminRoute := rg.Group("/admin")
	//adminRoute.Use(Middlewares.AuthValidationBearerMiddleware)
	//adminRoute.Use(Middlewares.RoleMiddleware("admin"))
	//{
	//	adminRoute.GET("/all", authController.GetAll)
	//	adminRoute.GET("/user/search", authController.SearchUser)
	//	adminRoute.GET("/user/student/all", authController.GetAllUserRoleIsStudent)
	//	adminRoute.GET("/user/teacher/all", authController.GetAllUserRoleIsTeacher)
	//	adminRoute.PATCH("/update/:id", authController.UpdateUser)
	//	adminRoute.DELETE("/delete/:user_id", authController.DeleteUser)
	//	adminRoute.GET("/user/pending-deletion", authController.GetListUserDependingDeletion)
	//	adminRoute.PATCH("/user/restore/:id", authController.RestoreUser)
	//}
}
