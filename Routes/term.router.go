package Routes

import (
	"gin-gonic-gom/Controllers"
	"github.com/gin-gonic/gin"
)

func TermRoutes(term *gin.RouterGroup) {
	termAdminRoute := term.Group("/admin/term")
	{
		//termAdminRoute.Use(Middlewares.AuthValidationBearerMiddleware)
		//termAdminRoute.Use(Middlewares.RoleMiddleware("admin"))
		{
			termAdminRoute.GET("/all", Controllers.GetAllTerms)
			termAdminRoute.GET("/details/:id", Controllers.GetTermDetail)
			termAdminRoute.POST("/add", Controllers.CreateTerm)
			termAdminRoute.DELETE("/:id", Controllers.DeleteTerm)
			termAdminRoute.PATCH("/:id", Controllers.UpdateTerm)
		}
	}
}
