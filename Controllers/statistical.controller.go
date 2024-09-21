package Controllers

import (
	"gin-gonic-gom/Common"
	"gin-gonic-gom/Middlewares"
	"gin-gonic-gom/Models"
	statiscal "gin-gonic-gom/Services/statistical"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type StatisticalController struct {
	StatisticalService statiscal.StatisticalService
}

func NewStatistical(statisticalService statiscal.StatisticalService) StatisticalController {
	return StatisticalController{
		StatisticalService: statisticalService,
	}
}
func (statisticalController *StatisticalController) StatisticalOfTerm(ctx *gin.Context) {
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
	res, err := statisticalController.StatisticalService.StatisticalOfTerm(page, limit)
	ctx.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Lấy thống kê theo học kỳ thành công!", res, 0, page, limit))
}
func (statisticalController *StatisticalController) ExportFileToExcel(ctx *gin.Context) {
	var statisticalExportInput []Models.StatisticalExportInput
	if err := ctx.ShouldBindJSON(&statisticalExportInput); err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, err.Error())
		return
	}
	err := statisticalController.StatisticalService.ExportStatisticalOfTerm(statisticalExportInput)
	if err != nil {
		Common.NewErrorResponse(ctx, http.StatusBadRequest, Common.ErrorExportFile, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Lấy thống kê theo học kỳ thành công!", nil))
}
func (statisticalController *StatisticalController) RegisterStatisticalRoutes(rg *gin.RouterGroup) {
	statisticalRoute := rg.Group("/statistical")
	{
		statisticalRoute.Use(Middlewares.AuthValidationBearerMiddleware)
		{
			statisticalRoute.POST("/term/export/excel", statisticalController.ExportFileToExcel)
		}
	}
	adminStatisticalRoute := rg.Group("/admin/statistical")
	{
		adminStatisticalRoute.Use(Middlewares.AuthValidationBearerMiddleware)
		adminStatisticalRoute.Use(Middlewares.RoleMiddleware("admin"))
		{
			adminStatisticalRoute.GET("/term", statisticalController.StatisticalOfTerm)
		}
	}
}
