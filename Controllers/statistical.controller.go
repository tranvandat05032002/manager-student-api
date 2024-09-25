package Controllers

import (
	"gin-gonic-gom/Collections"
	"gin-gonic-gom/Common"
	"gin-gonic-gom/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func StatisticalOfTerm(c *gin.Context) {
	var (
		DB = config.GetMongoDB()
	)
	limitStr := c.Query("limit")
	pageStr := c.Query("page")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}
	res, err := Collections.StatisticalOfTerm(DB, page, limit)
	c.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Lấy thống kê theo học kỳ thành công!", res, 0, page, limit))
}

func ExportFileToExcel(c *gin.Context) {
	var request []Collections.StatisticalExportInput
	if err := c.ShouldBindJSON(&request); err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, err.Error())
		return
	}
	err := Collections.ExportStatisticalOfTerm(request)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorExportFile, err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Lấy thống kê theo học kỳ thành công!", nil))
}
