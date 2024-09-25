package Routes

import (
	"gin-gonic-gom/Controllers"
	"gin-gonic-gom/Middlewares"
	"github.com/gin-gonic/gin"
)

func StatisticalRoutes(statistical *gin.RouterGroup) {
	statisticalRoute := statistical.Group("/statistical")
	{
		statisticalRoute.Use(Middlewares.AuthValidationBearerMiddleware)
		{
			statisticalRoute.POST("/term/export/excel", Controllers.ExportFileToExcel)
		}
	}
	adminStatisticalRoute := statistical.Group("/admin/statistical")
	{
		adminStatisticalRoute.Use(Middlewares.AuthValidationBearerMiddleware)
		adminStatisticalRoute.Use(Middlewares.RoleMiddleware("admin"))
		{
			adminStatisticalRoute.GET("/term", Controllers.StatisticalOfTerm)
		}
	}
}
