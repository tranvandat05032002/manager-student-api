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

func CreateMajor(c *gin.Context) {
	entry := Collections.MajorModel{}
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
	res, errCheckMajor := entry.CheckExist(DB, entry.MajorId, entry.MajorName)
	if errCheckMajor != nil {
		Common.NewErrorResponse(c, http.StatusInternalServerError, "Lỗi hệ thống! ", nil)
		return
	}
	if res {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Ngành đã tồn tại!", nil)
		return
	}
	err = entry.Create(DB)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Lỗi hệ thống!", nil)
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Thêm dữ liệu thành công!", nil))
}

func GetAllMajors(c *gin.Context) {
	var (
		major Collections.MajorModel
		DB    = config.GetMongoDB()
		err   error
		//Other config
		//.......
	)
	page, limit, skip := utils.Pagination(c)
	total, _ := major.Count(DB, bson.M{})
	res, err := major.Find(DB, limit, skip)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Không thể lấy thông tin!", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Lấy danh sách ngành thành công", res, int(total), page, limit))
}
func UpdateMajor(c *gin.Context) {
	request := Collections.MajorModel{}
	id := c.Param("id")
	var (
		major Collections.MajorModel
		DB    = config.GetMongoDB()
		err   error
		//Other config
		//.......
	)
	if err = c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		//Logger
		//Response
		errorMessages := utils.GetErrorMessagesResponse(err)
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, errorMessages)
		return
	}
	err = major.Update(DB, utils.ConvertStringToObjectId(id), request)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Cập nhật ngành thất bại!", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Cập nhật ngành thành công!", nil))
}

func GetDetailMajor(c *gin.Context) {
	id := c.Param("id")
	var (
		major Collections.MajorModel
		DB    = config.GetMongoDB()
		err   error
		//Other config
		//.......
	)
	res, err := major.FindByID(DB, utils.ConvertStringToObjectId(id))
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Ngành không tồn tại!", err.Error())
		return
	}
	c.JSON(
		http.StatusOK,
		Common.SimpleSuccessResponse(http.StatusOK, "Lấy thông tin ngành thành công!", res),
	)
}
func DeleteMajor(c *gin.Context) {
	id := c.Param("id")
	var (
		major Collections.MajorModel
		DB    = config.GetMongoDB()
		err   error
		//Other config
		//.......
	)
	_, err = major.FindByID(DB, utils.ConvertStringToObjectId(id))
	if err != nil && err.Error() == mongo.ErrNoDocuments.Error() {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Không tìm thấy ngành!", nil)
		return
	}
	err = major.Delete(DB, utils.ConvertStringToObjectId(id))
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Xóa ngành không thành công!", nil)
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Xóa ngành thành công", nil))
}
func SearchMajor(c *gin.Context) {
	var (
		major Collections.MajorModel
		err   error
		DB    = config.GetMongoDB()
		//Other config
		//.......
	)
	query := c.Query("major_name")
	page, limit, skip := utils.Pagination(c)
	res, total, err := major.Search(DB, query, skip, limit)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Tìm kiếm xảy ra lỗi!", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Tìm kiếm ngành thành công!!", res, total, page, limit))
}
