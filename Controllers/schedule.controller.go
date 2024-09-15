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

//	func (scheduleController *ScheduleController) UpdateScheduleController(ctx *gin.Context) {
//		var scheduleInput Models.ScheduleModel
//		id := ctx.Param("id")
//		if err := ctx.ShouldBindJSON(&subjectInput); err != nil {
//			errorMessages := utils.GetErrorMessagesResponse(err)
//			common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, errorMessages)
//			return
//		}
//		res, err := scheduleController.SubjectService.UpdateSubject(utils.ConvertStringToObjectId(id), subjectInput)
//		if err != nil {
//			common.NewErrorResponse(ctx, http.StatusBadRequest, "Cập nhật môn học thất bại!", err.Error())
//			return
//		}
//		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "Cập nhật môn học thành công!", res))
//	}
//
//	func (scheduleController *ScheduleController) GetAllscheduleController(ctx *gin.Context) {
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
//		subjects, total, err := scheduleController.ScheduleService.GetAllSubject(page, limit)
//		if err != nil {
//			common.NewErrorResponse(ctx, http.StatusBadRequest, "Không thể lấy thông tin!", err.Error())
//			return
//		}
//		ctx.JSON(http.StatusOK, common.NewSuccessResponse(http.StatusOK, "Lấy danh sách môn học thành công", subjects, total, page, limit))
//	}
//
//	func (scheduleController *ScheduleController) DeletescheduleController(ctx *gin.Context) {
//		majorId := ctx.Param("id")
//		res, err := scheduleController.ScheduleService.DeleteSubject(utils.ConvertStringToObjectId(majorId))
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
//	func (scheduleController *ScheduleController) GetSubjectDetailsController(ctx *gin.Context) {
//		subjectId := ctx.Param("id")
//		subject, err := scheduleController.ScheduleService.GetSubjectDetails(utils.ConvertStringToObjectId(subjectId))
//		if err != nil {
//			common.NewErrorResponse(ctx, http.StatusBadRequest, "Môn học không tồn tại!", err.Error())
//			return
//		}
//		ctx.JSON(
//			http.StatusOK,
//			common.SimpleSuccessResponse(http.StatusOK, "Lấy thông tin môn học thành công!", subject),
//		)
//	}
//
//	func (scheduleController *ScheduleController) SearchscheduleController(ctx *gin.Context) {
//		subjectNameQuery := ctx.Query("subject_name")
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
//		nameSubject := string(subjectNameQuery)
//		subjects, total, _ := scheduleController.ScheduleService.SearchSubject(nameSubject, page, limit)
//		ctx.JSON(http.StatusOK, common.NewSuccessResponse(http.StatusOK, "Tìm kiếm môn học thành công!!", subjects, total, page, limit))
//	}
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
			adminscheduleroute.POST("/add", scheduleController.CreateScheduleController)
		}
	}
}
