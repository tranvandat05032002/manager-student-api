package Controllers

import (
	"fmt"
	"gin-gonic-gom/Collections"
	"gin-gonic-gom/Common"
	"gin-gonic-gom/config"
	"gin-gonic-gom/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func CreateSchedule(c *gin.Context) {
	entry := Collections.ScheduleModel{}
	var (
		DB  = config.GetMongoDB()
		err error
		//Other config
		//.......
	)
	if err = c.ShouldBindBodyWith(&entry, binding.JSON); err != nil {
		//Logger
		//Response
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	fmt.Println("Data --> ", entry)
	res, errCheckSchedule := entry.CheckExist(DB, entry.Room, entry.DayOfWeek, entry.StartPeriod, entry.EndPeriod)
	if errCheckSchedule != nil {
		Common.NewErrorResponse(c, http.StatusInternalServerError, "Lỗi hệ thống! ", nil)
		return
	}
	if res {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Lịch học đã tồn tại!", nil)
		return
	}
	err = entry.Create(DB)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Lỗi hệ thống!", nil)
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Thêm dữ liệu thành công!", nil))
}

func GetAllSchedules(c *gin.Context) {
	var (
		schedule Collections.ScheduleModel
		DB       = config.GetMongoDB()
		err      error
		//Other config
		//.......
	)
	page, limit, skip := utils.Pagination(c)
	total, _ := schedule.Count(DB, bson.M{})
	res, err := schedule.Find(DB, limit, skip)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Không thể lấy thông tin!", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Lấy danh sách lịch học thành công", res, int(total), page, limit))
}
func UpdateSchedule(c *gin.Context) {
	id := c.Param("id")
	var (
		schedule Collections.ScheduleModel
		DB       = config.GetMongoDB()
		err      error
		//Other config
		//.......
	)
	if err = c.ShouldBindBodyWith(&schedule, binding.JSON); err != nil {
		//Logger
		//Response
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	err = schedule.Update(DB, utils.ConvertStringToObjectId(id), schedule)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Cập nhật lịch học thất bại!", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Cập nhật lịch học thành công!", nil))
}

func GetDetailSchedule(c *gin.Context) {
	id := c.Param("id")
	var (
		schedule Collections.ScheduleModel
		DB       = config.GetMongoDB()
		err      error
		//Other config
		//.......
	)
	res, err := schedule.FindByID(DB, utils.ConvertStringToObjectId(id))
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Lịch học không tồn tại!", err.Error())
		return
	}
	c.JSON(
		http.StatusOK,
		Common.SimpleSuccessResponse(http.StatusOK, "Lấy thông tin lịch học thành công!", res),
	)
}
func DeleteSchedule(c *gin.Context) {
	id := c.Param("id")
	var (
		schedule Collections.ScheduleModel
		DB       = config.GetMongoDB()
		err      error
		//Other config
		//.......
	)
	_, err = schedule.FindByID(DB, utils.ConvertStringToObjectId(id))
	if err != nil && err.Error() == mongo.ErrNoDocuments.Error() {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Không tìm thấy ngành!", nil)
		return
	}
	err = schedule.Delete(DB, utils.ConvertStringToObjectId(id))
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Xóa lịch học không thành công!", nil)
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Xóa lịch học thành công", nil))
}
