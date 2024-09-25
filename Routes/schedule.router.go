package Routes

import (
	"gin-gonic-gom/Controllers"
	"gin-gonic-gom/Middlewares"
	"github.com/gin-gonic/gin"
)

func ScheduleRoutes(rg *gin.RouterGroup) {
	scheduleRoute := rg.Group("/schedule")
	{
		scheduleRoute.Use(Middlewares.AuthValidationBearerMiddleware)
		{
			//scheduleroute.POST("/term/export/excel", statisticalController.ExportFileToExcel)
		}
	}
	adminsScheduleRoute := rg.Group("/admin/schedule")
	{
		adminsScheduleRoute.Use(Middlewares.AuthValidationBearerMiddleware)
		adminsScheduleRoute.Use(Middlewares.RoleMiddleware("admin"))
		{
			adminsScheduleRoute.GET("/list", Controllers.GetAllSchedules)
			adminsScheduleRoute.GET("/details/:id", Controllers.GetDetailSchedule)
			adminsScheduleRoute.POST("/add", Controllers.CreateSchedule)
			adminsScheduleRoute.PATCH("/:id", Controllers.UpdateSchedule)
			adminsScheduleRoute.DELETE("/:id", Controllers.DeleteSchedule)
		}
	}
}
