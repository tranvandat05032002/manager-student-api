package Routes

import (
	"gin-gonic-gom/Middlewares"
	"github.com/gin-gonic/gin"
)

func Router(rg *gin.RouterGroup) {
	// middleware global server
	// prefix router
	rg.Use(Middlewares.CorsMiddleware())
	AuthRoutes(rg)
	MajorRoutes(rg)
	SubjectRoutes(rg)
	MediaRoutes(rg)
	TermRoutes(rg)
	ScheduleRoutes(rg)
}
