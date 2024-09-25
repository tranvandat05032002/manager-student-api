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
		authRoute.POST("/forgot-password", Controllers.ResetPassword)
		authRoute.Use(Middlewares.AuthValidationBearerMiddleware)
		{
			authRoute.GET("/me", Controllers.GetMe)
			authRoute.POST("/logout", Controllers.SignOut)
			authRoute.GET("/:user_id", Controllers.GetDetailUser)
			authRoute.PATCH("/me/update", Controllers.UpdateMe)
			authRoute.PUT("/change-password", Controllers.ChangePassword)
		}
	}
	// Admin
	adminRoute := auth.Group("/admin")
	adminRoute.Use(Middlewares.AuthValidationBearerMiddleware)
	adminRoute.Use(Middlewares.RoleMiddleware("admin"))
	{
		adminRoute.GET("/list", Controllers.GetAll)
		adminRoute.GET("/user/search", Controllers.SearchUser)
		adminRoute.GET("/user/students", Controllers.GetStudent)
		adminRoute.GET("/user/teachers", Controllers.GetTeacher)
		adminRoute.PATCH("/update/:id", Controllers.UpdateUser)
		adminRoute.DELETE("/delete/:user_id", Controllers.DeleteUser)
		adminRoute.GET("/user/pending-deletion", Controllers.GetUserDepending)
		adminRoute.PATCH("/user/restore/:id", Controllers.RestoreUser)
	}
}
