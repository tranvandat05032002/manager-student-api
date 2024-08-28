package Controllers

import (
	"gin-gonic-gom/Middlewares"
	"gin-gonic-gom/Models"
	"gin-gonic-gom/Services/subject"
	"gin-gonic-gom/common"
	"gin-gonic-gom/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type SubjectController struct {
	SubjectService subject.SubjectService
}

func NewSubject(subjectService subject.SubjectService) SubjectController {
	return SubjectController{
		SubjectService: subjectService,
	}
}
func (subjectController *SubjectController) CreateSubjectController(ctx *gin.Context) {
	var subjectInput Models.SubjectInput
	if err := ctx.ShouldBindJSON(&subjectInput); err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, err.Error())
		return
	}
	err := subjectController.SubjectService.CreateSubject(subjectInput)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorAddDataMessage, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, common.SuccessAddDataMessage, nil))
}

func (subjectController *SubjectController) UpdateSubjectController(ctx *gin.Context) {
	var subjectInput Models.SubjectInput
	id := ctx.Param("id")
	if err := ctx.ShouldBindJSON(&subjectInput); err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, err.Error())
		return
	}
	res, err := subjectController.SubjectService.UpdateSubject(utils.ConvertStringToObjectId(id), subjectInput)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Cập nhật môn học thất bại!", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "Cập nhật môn học thành công!", res))
}

func (subjectController *SubjectController) GetAllSubjectController(ctx *gin.Context) {
	limitStr := ctx.Query("limit")
	pageStr := ctx.Query("page")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	subjects, total, err := subjectController.SubjectService.GetAllSubject(page, limit)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Không thể lấy thông tin!", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.NewSuccessResponse(http.StatusOK, "Lấy danh sách môn học thành công", subjects, total, page, limit))
}
func (subjectController *SubjectController) DeleteSubjectController(ctx *gin.Context) {
	majorId := ctx.Param("id")
	res, err := subjectController.SubjectService.DeleteSubject(utils.ConvertStringToObjectId(majorId))
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Xóa môn học không thành công!", err.Error())
		return
	}
	if res < 1 {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Môn học không tồn tại!", "")
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "Xóa môn học thành công", nil))
}
func (subjectController *SubjectController) GetSubjectDetailsController(ctx *gin.Context) {
	subjectId := ctx.Param("id")
	subject, err := subjectController.SubjectService.GetSubjectDetails(utils.ConvertStringToObjectId(subjectId))
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Môn học không tồn tại!", err.Error())
		return
	}
	ctx.JSON(
		http.StatusOK,
		common.SimpleSuccessResponse(http.StatusOK, "Lấy thông tin môn học thành công!", subject),
	)
}
func (subjectController *SubjectController) RegisterSubjectRoutes(rg *gin.RouterGroup) {
	subjectroute := rg.Group("/subject")
	{
		subjectroute.Use(Middlewares.AuthValidationBearerMiddleware)
		{
			subjectroute.GET("/all", subjectController.GetAllSubjectController)
			subjectroute.GET("/details/:id", subjectController.GetSubjectDetailsController)
		}
	}
	subjectadminroute := rg.Group("/admin/subject")
	{
		subjectadminroute.Use(Middlewares.AuthValidationBearerMiddleware)
		subjectadminroute.Use(Middlewares.RoleMiddleware("admin"))
		{
			subjectadminroute.POST("/add", subjectController.CreateSubjectController)
			subjectadminroute.DELETE("/:id", subjectController.DeleteSubjectController)
			subjectadminroute.PATCH("/:id", subjectController.UpdateSubjectController)
		}
	}
}
