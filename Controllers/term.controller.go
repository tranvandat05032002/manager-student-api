package Controllers

import (
	"gin-gonic-gom/Common"
	"gin-gonic-gom/Middlewares"
	"gin-gonic-gom/Models"
	"gin-gonic-gom/Services/term"
	"gin-gonic-gom/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	var termEnity Models.TermInput
	if err := ctx.ShouldBindJSON(&termEnity); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	err := termController.TermService.CreateTerm(termEnity)
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorAddDataMessage, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, Common.SuccessAddDataMessage, nil))
}

func (termController *TermController) UpdateTermController(ctx *gin.Context) {
	var termEnity Models.TermInput
	id := ctx.Param("id")
	if err := ctx.ShouldBindJSON(&termEnity); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	res, err := termController.TermService.UpdateTerm(utils.ConvertStringToObjectId(id), termEnity)
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, "Cập nhật học kỳ trong năm thất bại!", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Cập nhật học kỳ thành công!", res))
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
	res, total, err := termController.TermService.GetAllTerm(page, limit)
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, "Không thể lấy thông tin!", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Lấy danh sách học kỳ thành công", res, total, page, limit))
}
func (termController *TermController) DeleteTermController(ctx *gin.Context) {
	termId := ctx.Param("id")
	res, err := termController.TermService.DeleteTerm(utils.ConvertStringToObjectId(termId))
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, "Xóa học kỳ không thành công!", err.Error())
		return
	}
	if res < 1 {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, "Học kỳ không tồn tại!", "")
		return
	}
	ctx.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Xóa học kỳ thành công", nil))
}
func (termController *TermController) GetTermDetailsController(ctx *gin.Context) {
	termId := ctx.Param("id")
	res, err := termController.TermService.GetTermDetails(utils.ConvertStringToObjectId(termId))
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, "Học kỳ không tồn tại!", err.Error())
		return
	}
	ctx.JSON(
		http.StatusOK,
		Common.SimpleSuccessResponse(http.StatusOK, "Lấy thông tin học kỳ thành công!", res),
	)
}
func (termController *TermController) RegisterTermRoutes(rg *gin.RouterGroup) {
	termAdminRoute := rg.Group("/admin/term")
	{
		termAdminRoute.Use(Middlewares.AuthValidationBearerMiddleware)
		termAdminRoute.Use(Middlewares.RoleMiddleware("admin"))
		{
			termAdminRoute.GET("/all", termController.GetAllTermController)
			termAdminRoute.GET("/details/:id", termController.GetTermDetailsController)
			termAdminRoute.POST("/add", termController.CreateTermController)
			termAdminRoute.DELETE("/:id", termController.DeleteTermController)
			termAdminRoute.PATCH("/:id", termController.UpdateTermController)
		}
	}
}
