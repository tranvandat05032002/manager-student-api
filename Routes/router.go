package Routes

import (
	"github.com/gin-gonic/gin"
)

func Router(rg *gin.RouterGroup) {
	// middleware global server
	// prefix router
	rg.Group("/v1/api")
	MajorRouter(rg)
}
