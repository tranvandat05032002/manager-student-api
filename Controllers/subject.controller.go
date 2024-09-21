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

func CreateSubject(c *gin.Context) {
	entry := Collections.SubjectModel{}
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
	res, errCheckTerm := entry.CheckExist(DB, entry.SubjectCode)
	if errCheckTerm != nil {
		Common.NewErrorResponse(c, http.StatusInternalServerError, "Lỗi hệ thống! ", nil)
		return
	}
	if res {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Môn học đã tồn tại!", nil)
		return
	}
	err = entry.Create(DB)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Lỗi hệ thống!", nil)
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, Common.SuccessAddDataMessage, nil))
}
func UpdateSubject(c *gin.Context) {
	entry := Collections.SubjectModel{}
	var (
		DB  = config.GetMongoDB()
		err error
		//Other config
		//.......
	)
	id := c.Param("id")
	if err = c.ShouldBindBodyWith(&entry, binding.JSON); err != nil {
		//Logger
		//Response
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	err = entry.Update(DB, utils.ConvertStringToObjectId(id), entry)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Cập nhật môn học thất bại!", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Cập nhật môn học thành công!", nil))
}
func DeleteSubject(c *gin.Context) {
	id := c.Param("id")
	subID := utils.ConvertStringToObjectId(id)
	var (
		subject Collections.SubjectModel
		DB      = config.GetMongoDB()
		err     error
		//Other config
		//.......
	)
	_, err = subject.FindByID(DB, subID)
	if err != nil && err.Error() == mongo.ErrNoDocuments.Error() {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Không tìm thấy môn học!", nil)
		return
	}
	err = subject.Delete(DB, subID)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Xóa môn học không thành công!", nil)
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Xóa môn học thành công", nil))
}
func GetSubjectDetail(c *gin.Context) {
	id := c.Param("id")
	var (
		subject Collections.SubjectModel
		DB      = config.GetMongoDB()
		err     error
		//Other config
		//.......
	)
	res, err := subject.FindByID(DB, utils.ConvertStringToObjectId(id))
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Môn học không tồn tại!", err.Error())
		return
	}
	c.JSON(
		http.StatusOK,
		Common.SimpleSuccessResponse(http.StatusOK, "Lấy thông tin môn học thành công!", res),
	)
}
func GetAllSubjects(c *gin.Context) {
	var (
		subject Collections.SubjectModel
		DB      = config.GetMongoDB()
		err     error
		//Other config
		//.......
	)
	page, limit, skip := utils.Pagination(c)
	res, err := subject.Find(DB, limit, skip)
	total, _ := subject.Count(DB, bson.M{})
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Không thể lấy thông tin!", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Lấy danh sách môn học thành công", res, int(total), page, limit))
}
func SearchSubject(c *gin.Context) {
	var (
		subject Collections.SubjectModel
		err     error
		DB      = config.GetMongoDB()
		//Other config
		//.......
	)
	query := c.Query("subject_name")
	page, limit, skip := utils.Pagination(c)
	res, total, err := subject.Search(DB, query, skip, limit)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Tìm kiếm xảy ra lỗi!", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Tìm kiếm môn học thành công!!", res, total, page, limit))
}
