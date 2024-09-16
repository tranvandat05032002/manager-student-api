package Controllers

import (
	"gin-gonic-gom/Middlewares"
	"gin-gonic-gom/Models"
	"gin-gonic-gom/Services/schedule"
	"gin-gonic-gom/common"
	"gin-gonic-gom/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
)

type ScheduleController struct {
	ScheduleService schedule.ScheduleService
}

func NewSchedule(scheduleService schedule.ScheduleService) ScheduleController {
	return ScheduleController{
		ScheduleService: scheduleService,
	}
}
func (scheduleController *ScheduleController) CreateScheduleController(ctx *gin.Context) {
	scheduleInput := Models.ScheduleModel{}
	if err := ctx.ShouldBindBodyWith(&scheduleInput, binding.JSON); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	err := scheduleController.ScheduleService.CreateSchedule(scheduleInput)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorAddDataMessage, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, common.SuccessAddDataMessage, nil))
}

func (scheduleController *ScheduleController) UpdateScheduleController(ctx *gin.Context) {
	var scheduleInput Models.ScheduleModel
	id := ctx.Param("id")
	if err := ctx.ShouldBindJSON(&scheduleInput); err != nil {
		errorMessages := utils.GetErrorMessagesResponse(err)
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	err := scheduleController.ScheduleService.UpdateSchedule(utils.ConvertStringToObjectId(id), scheduleInput)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Cập nhật lịch học thất bại!", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "Cập nhật lịch học thành công!", nil))
}

func (scheduleController *ScheduleController) GetAllScheduleController(ctx *gin.Context) {
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
	schedules, total, err := scheduleController.ScheduleService.GetAllSchedule(page, limit)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Không thể lấy thông tin!", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.NewSuccessResponse(http.StatusOK, "Lấy danh sách lịch học thành công", schedules, total, page, limit))
}

func (scheduleController *ScheduleController) DeleteScheduleController(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := scheduleController.ScheduleService.DeleteSchedule(utils.ConvertStringToObjectId(id))
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Xóa lịch học không thành công!", err.Error())
		return
	}
	if res < 1 {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Lịch học không tồn tại!", "")
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "Xóa lịch học thành công", nil))
}

func (scheduleController *ScheduleController) GetScheduleDetailsController(ctx *gin.Context) {
	id := ctx.Param("id")
	schedule, err := scheduleController.ScheduleService.GetScheduleDetails(utils.ConvertStringToObjectId(id))
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Lịch học không tồn tại!", err.Error())
		return
	}
	ctx.JSON(
		http.StatusOK,
		common.SimpleSuccessResponse(http.StatusOK, "Lấy thông tin lịch học thành công!", schedule),
	)
}

func (scheduleController *ScheduleController) RegisterScheduleRoutes(rg *gin.RouterGroup) {
	scheduleroute := rg.Group("/schedule")
	{
		scheduleroute.Use(Middlewares.AuthValidationBearerMiddleware)
		{
			//scheduleroute.POST("/term/export/excel", statisticalController.ExportFileToExcel)
		}
	}
	adminscheduleroute := rg.Group("/admin/schedule")
	{
		adminscheduleroute.Use(Middlewares.AuthValidationBearerMiddleware)
		adminscheduleroute.Use(Middlewares.RoleMiddleware("admin"))
		{
			adminscheduleroute.GET("/all", scheduleController.GetAllScheduleController)
			adminscheduleroute.GET("/details/:id", scheduleController.GetScheduleDetailsController)
			adminscheduleroute.POST("/add", scheduleController.CreateScheduleController)
			adminscheduleroute.PATCH("/:id", scheduleController.UpdateScheduleController)
			adminscheduleroute.DELETE("/:id", scheduleController.DeleteScheduleController)
		}
	}
}
