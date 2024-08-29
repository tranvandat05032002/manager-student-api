package Controllers

import (
	"gin-gonic-gom/Middlewares"
	"gin-gonic-gom/Models"
	statiscal "gin-gonic-gom/Services/statistical"
	"gin-gonic-gom/common"
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
	ctx.JSON(http.StatusOK, common.NewSuccessResponse(http.StatusOK, "Lấy thống kê theo học kỳ thành công!", res, 0, page, limit))
}
func (statisticalController *StatisticalController) ExportFileToExcel(ctx *gin.Context) {
	var statisticalExportInput []Models.StatisticalExportInput
	if err := ctx.ShouldBindJSON(&statisticalExportInput); err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, err.Error())
		return
	}
	err := statisticalController.StatisticalService.ExportStatisticalOfTerm(statisticalExportInput)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorExportFile, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "Lấy thống kê theo học kỳ thành công!", nil))
}
func (statisticalController *StatisticalController) RegisterStatisticalRoutes(rg *gin.RouterGroup) {
	statisticalroute := rg.Group("/statistical")
	{
		statisticalroute.Use(Middlewares.AuthValidationBearerMiddleware)
		{
			statisticalroute.POST("/term/export/excel", statisticalController.ExportFileToExcel)
		}
	}
	adminstatisticalroute := rg.Group("/admin/statistical")
	{
		adminstatisticalroute.Use(Middlewares.AuthValidationBearerMiddleware)
		adminstatisticalroute.Use(Middlewares.RoleMiddleware("admin"))
		{
			adminstatisticalroute.GET("/term", statisticalController.StatisticalOfTerm)
		}
	}
}
