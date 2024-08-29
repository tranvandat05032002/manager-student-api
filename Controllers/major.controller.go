package Controllers

import (
	"fmt"
	"gin-gonic-gom/Middlewares"
	"gin-gonic-gom/Models"
	"gin-gonic-gom/Services/major"
	"gin-gonic-gom/common"
	"gin-gonic-gom/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MajorController struct {
	MajorService major.MajorService
}

func NewMajor(majorService major.MajorService) MajorController {
	return MajorController{
		MajorService: majorService,
	}
}
func (majorController *MajorController) CreateMajorController(ctx *gin.Context) {
	var major Models.MajorModel
	if err := ctx.ShouldBindJSON(&major); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	err := majorController.MajorService.CreateMajor(&major)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorAddDataMessage, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "Tạo ngành thành công!!!", nil))
}
func (majorController *MajorController) UpdateMajorController(ctx *gin.Context) {
	var majorUpdateReq Models.MajorUpdateReq
	id := ctx.Param("id")
	if err := ctx.ShouldBindJSON(&majorUpdateReq); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	res, err := majorController.MajorService.UpdateMajor(utils.ConvertStringToObjectId(id), &majorUpdateReq)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Cập nhật ngành thất bại!", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "Cập nhật ngành thành công!", res))
}
func (majorController *MajorController) GetAllMajorController(ctx *gin.Context) {
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
	majors, total, err := majorController.MajorService.GetAllMajor(page, limit)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Không thể lấy thông tin!", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.NewSuccessResponse(http.StatusOK, "Lấy danh sách ngành thành công", majors, total, page, limit))
}
func (majorController *MajorController) DeleteMajorController(ctx *gin.Context) {
	majorId := ctx.Param("id")
	res, err := majorController.MajorService.DeleteMajor(utils.ConvertStringToObjectId(majorId))
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Xóa ngành không thành công!", err.Error())
		return
	}
	if res < 1 {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Ngành không tồn tại!", "")
		return
	}
	message := fmt.Sprintf("Xóa ngành thành công")
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, message, nil))
}
func (majorController *MajorController) GetMajorDetailsController(ctx *gin.Context) {
	majorId := ctx.Param("id")
	major, err := majorController.MajorService.GetMajorDetails(utils.ConvertStringToObjectId(majorId))
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Ngành không tồn tại!", err.Error())
		return
	}
	ctx.JSON(
		http.StatusOK,
		common.SimpleSuccessResponse(http.StatusOK, "Lấy thông tin ngành thành công!", major),
	)
}
func (majorController *MajorController) RegisterMajorRoutes(rg *gin.RouterGroup) {
	majorroute := rg.Group("/major") // Client
	{
		majorroute.Use(Middlewares.AuthValidationBearerMiddleware)
		{
			majorroute.GET("/details/:id", majorController.GetMajorDetailsController)
			majorroute.GET("/all", majorController.GetAllMajorController)
		}
	}
	majoradminroute := rg.Group("/admin/major")
	{
		majoradminroute.Use(Middlewares.AuthValidationBearerMiddleware)
		majoradminroute.Use(Middlewares.RoleMiddleware("admin"))
		{
			majoradminroute.POST("/add", majorController.CreateMajorController)
			majoradminroute.DELETE("/:id", majorController.DeleteMajorController)
			majoradminroute.PATCH("/:id", majorController.UpdateMajorController)
		}
	}
}
