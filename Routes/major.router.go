package Routes

import (
	"gin-gonic-gom/Controllers"
	"gin-gonic-gom/Middlewares"
	"github.com/gin-gonic/gin"
)

func MajorRoutes(major *gin.RouterGroup) {
	majorAdminRoute := major.Group("/admin/major")
	{
		majorAdminRoute.Use(Middlewares.AuthValidationBearerMiddleware)
		majorAdminRoute.Use(Middlewares.RoleMiddleware("admin"))
		{
			majorAdminRoute.GET("/details/:id", Controllers.GetDetailMajor)
			majorAdminRoute.GET("/list", Controllers.GetAllMajors)
			majorAdminRoute.POST("/add", Controllers.CreateMajor)
			majorAdminRoute.GET("/search", Controllers.SearchMajor)
			majorAdminRoute.DELETE("/:id", Controllers.DeleteMajor)
			majorAdminRoute.PATCH("/:id", Controllers.UpdateMajor)
		}
	}
}
