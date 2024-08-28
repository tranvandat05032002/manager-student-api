package Controllers

import (
	"gin-gonic-gom/Middlewares"
	"gin-gonic-gom/Models"
	"gin-gonic-gom/Services/term"
	"gin-gonic-gom/common"
	"gin-gonic-gom/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	err := termController.TermService.CreateTerm(termInput)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorAddDataMessage, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, common.SuccessAddDataMessage, nil))
}

func (termController *TermController) UpdateTermController(ctx *gin.Context) {
	var termInput Models.TermInput
	id := ctx.Param("id")
	if err := ctx.ShouldBindJSON(&termInput); err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, err.Error())
		return
	}
	res, err := termController.TermService.UpdateTerm(utils.ConvertStringToObjectId(id), termInput)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Cập nhật học kỳ trong năm thất bại!", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "Cập nhật học kỳ thành công!", res))
}

func (termController *TermController) GetAllTermController(ctx *gin.Context) {
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
	terms, total, err := termController.TermService.GetAllTerm(page, limit)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Không thể lấy thông tin!", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.NewSuccessResponse(http.StatusOK, "Lấy danh sách học kỳ thành công", terms, total, page, limit))
}
func (termController *TermController) DeleteTermController(ctx *gin.Context) {
	termId := ctx.Param("id")
	res, err := termController.TermService.DeleteTerm(utils.ConvertStringToObjectId(termId))
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Xóa học kỳ không thành công!", err.Error())
		return
	}
	if res < 1 {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Học kỳ không tồn tại!", "")
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "Xóa học kỳ thành công", nil))
}
func (termController *TermController) GetTermDetailsController(ctx *gin.Context) {
	termId := ctx.Param("id")
	term, err := termController.TermService.GetTermDetails(utils.ConvertStringToObjectId(termId))
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Học kỳ không tồn tại!", err.Error())
		return
	}
	ctx.JSON(
		http.StatusOK,
		common.SimpleSuccessResponse(http.StatusOK, "Lấy thông tin môn học thành công!", term),
	)
}
func (termController *TermController) RegisterTermRoutes(rg *gin.RouterGroup) {
	termroute := rg.Group("/term")
	{
		termroute.Use(Middlewares.AuthValidationBearerMiddleware)
		{
			termroute.GET("/all", termController.GetAllTermController)
			termroute.GET("/details/:id", termController.GetTermDetailsController)
		}
	}
	termadminroute := rg.Group("/admin/term")
	{
		termadminroute.Use(Middlewares.AuthValidationBearerMiddleware)
		termadminroute.Use(Middlewares.RoleMiddleware("admin"))
		{
			termadminroute.POST("/add", termController.CreateTermController)
			termadminroute.DELETE("/:id", termController.DeleteTermController)
			termadminroute.PATCH("/:id", termController.UpdateTermController)
		}
	}
}
