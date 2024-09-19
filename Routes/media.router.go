package Routes

import (
	"gin-gonic-gom/Controllers"
	"github.com/gin-gonic/gin"
)

func MediaRoutes(media *gin.RouterGroup) {
	mediaRoute := media.Group("/upload")
	{
		//mediaRoute.Use(Middlewares.AuthValidationBearerMiddleware)
		{
			mediaRoute.POST("/image", Controllers.UploadImage)
		}
	}
	adminMediaRoute := media.Group("/admin/excel/upload")
	{
		//adminMediaRoute.Use(Middlewares.AuthValidationBearerMiddleware)
		//adminMediaRoute.Use(Middlewares.RoleMiddleware("admin"))
		{
			adminMediaRoute.POST("/user", Controllers.UploadUserByExcel)
		}
	}
}
