package Routes

import (
	"gin-gonic-gom/Controllers"
	"gin-gonic-gom/Middlewares"
	"github.com/gin-gonic/gin"
)

func SubjectRoutes(subject *gin.RouterGroup) {
	subjectAdminRoute := subject.Group("/admin/subject")
	{
		subjectAdminRoute.Use(Middlewares.AuthValidationBearerMiddleware)
		subjectAdminRoute.Use(Middlewares.RoleMiddleware("admin"))
		{
			subjectAdminRoute.GET("/list", Controllers.GetAllSubjects)
			subjectAdminRoute.GET("/details/:id", Controllers.GetSubjectDetail)
			subjectAdminRoute.POST("/add", Controllers.CreateSubject)
			subjectAdminRoute.GET("/search", Controllers.SearchSubject)
			subjectAdminRoute.DELETE("/:id", Controllers.DeleteSubject)
			subjectAdminRoute.PATCH("/:id", Controllers.UpdateSubject)
		}
	}
}
