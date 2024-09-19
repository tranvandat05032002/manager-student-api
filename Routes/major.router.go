package Routes

import (
	"gin-gonic-gom/Controllers"
	"github.com/gin-gonic/gin"
)

func MajorRouter(major *gin.RouterGroup) {
	majorAdminRoute := major.Group("/admin/major")
	{
		//majorAdminRoute.Use(Middlewares.AuthValidationBearerMiddleware)
		//majorAdminRoute.Use(Middlewares.RoleMiddleware("admin"))
		{
			majorAdminRoute.GET("/details/:id", Controllers.GetDetailMajor)
			majorAdminRoute.GET("/all", Controllers.GetAllMajors)
			majorAdminRoute.POST("/add", Controllers.PostMajor)
			majorAdminRoute.GET("/search", Controllers.SearchMajor)
			majorAdminRoute.DELETE("/:id", Controllers.DeleteMajor)
			majorAdminRoute.PATCH("/:id", Controllers.UpdateMajor)
		}
	}
}
