package Controllers

import (
	"fmt"
	"gin-gonic-gom/Middlewares"
	"gin-gonic-gom/Models"
	"gin-gonic-gom/Services/term"
	"gin-gonic-gom/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type TermController struct {
	TermService term.TermService
}

func NewTerm(termService term.TermService) TermController {
	return TermController{
		TermService: termService,
	}
}
func (termController *TermController) CreateTermController(ctx *gin.Context) {
	var termInput Models.TermInput
	if err := ctx.ShouldBindJSON(&termInput); err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, err.Error())
		return
	}
	fmt.Println("Type ---> ", reflect.TypeOf(termInput.StartDate))
	fmt.Println("Data term ---> ", termInput)
	err := termController.TermService.CreateTerm(termInput)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorAddDataMessage, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, common.SuccessAddDataMessage, nil))
}

//	func (subjectController *SubjectController) UpdateSubjectController(ctx *gin.Context) {
//		var subjectInput bson.M
//		id := ctx.Param("id")
//		if err := ctx.ShouldBindJSON(&subjectInput); err != nil {
//			common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, err.Error())
//			return
//		}
//		fmt.Println("DataInput ----> ", subjectInput)
//		res, err := subjectController.SubjectService.UpdateSubject(utils.ConvertStringToObjectId(id), subjectInput)
//		if err != nil {
//			common.NewErrorResponse(ctx, http.StatusBadRequest, "Cập nhật môn học thất bại!", err.Error())
//			return
//		}
//		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "Cập nhật môn học thành công!", res))
//	}
//
//	func (subjectController *SubjectController) GetAllSubjectController(ctx *gin.Context) {
//		limitStr := ctx.Query("limit")
//		pageStr := ctx.Query("page")
//		limit, err := strconv.Atoi(limitStr)
//		if err != nil || limit <= 0 {
//			limit = 5
//		}
//		page, err := strconv.Atoi(pageStr)
//		if err != nil || page <= 0 {
//			page = 1
//		}
//		fmt.Println("Limit --> ", limit, "page ---> ", page)
//		subjects, total, err := subjectController.SubjectService.GetAllSubject(page, limit)
//		if err != nil {
//			common.NewErrorResponse(ctx, http.StatusBadRequest, "Không thể lấy thông tin!", err.Error())
//			return
//		}
//		ctx.JSON(http.StatusOK, common.NewSuccessResponse(http.StatusOK, "Lấy danh sách môn học thành công", subjects, total, page, limit))
//	}
//
//	func (subjectController *SubjectController) DeleteSubjectController(ctx *gin.Context) {
//		majorId := ctx.Param("id")
//		res, err := subjectController.SubjectService.DeleteSubject(utils.ConvertStringToObjectId(majorId))
//		if err != nil {
//			common.NewErrorResponse(ctx, http.StatusBadRequest, "Xóa môn học không thành công!", err.Error())
//			return
//		}
//		if res < 1 {
//			common.NewErrorResponse(ctx, http.StatusBadRequest, "Môn học không tồn tại!", "")
//			return
//		}
//		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "Xóa môn học thành công", nil))
//	}
//
//	func (subjectController *SubjectController) GetSubjectDetailsController(ctx *gin.Context) {
//		subjectId := ctx.Param("id")
//		subject, err := subjectController.SubjectService.GetSubjectDetails(utils.ConvertStringToObjectId(subjectId))
//		if err != nil {
//			common.NewErrorResponse(ctx, http.StatusBadRequest, "Môn học không tồn tại!", err.Error())
//			return
//		}
//		ctx.JSON(
//			http.StatusOK,
//			common.SimpleSuccessResponse(http.StatusOK, "Lấy thông tin môn học thành công!", subject),
//		)
//	}
func (termController *TermController) RegisterTermRoutes(rg *gin.RouterGroup) {
	termroute := rg.Group("/term")
	{
		termroute.Use(Middlewares.AuthValidationBearerMiddleware)
		{
			//termroute.GET("/all", termController.GetAllTermController)
			//termroute.GET("/details/:id", termController.GetTermDetailsController)
		}
	}
	termadminroute := rg.Group("/admin/term")
	{
		termadminroute.Use(Middlewares.AuthValidationBearerMiddleware)
		termadminroute.Use(Middlewares.RoleMiddleware("admin"))
		{
			termadminroute.POST("/add", termController.CreateTermController)
			//termadminroute.DELETE("/:id", termController.DeleteTermController)
			//termadminroute.PATCH("/:id", termController.UpdateTermController)
		}
	}
}
