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

func CreateTerm(c *gin.Context) {
	entry := Collections.TermModel{}
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
	res, errCheckTerm := entry.CheckExist(DB, entry.TermSemester, entry.TermFromYear, entry.TermToYear)
	if errCheckTerm != nil {
		Common.NewErrorResponse(c, http.StatusInternalServerError, "Lỗi hệ thống! ", nil)
		return
	}
	if res {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Học kỳ đã tồn tại!", nil)
		return
	}
	err = entry.Create(DB)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Lỗi hệ thống!", nil)
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Thêm dữ liệu thành công!", nil))
}
func UpdateTerm(c *gin.Context) {
	entry := Collections.TermModel{}
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
		Common.NewErrorResponse(c, http.StatusBadRequest, "Cập nhật học kỳ trong năm thất bại!", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Cập nhật học kỳ thành công!", nil))

}
func DeleteTerm(c *gin.Context) {
	id := c.Param("id")
	termID := utils.ConvertStringToObjectId(id)
	var (
		term Collections.TermModel
		DB   = config.GetMongoDB()
		err  error
		//Other config
		//.......
	)
	_, err = term.FindByID(DB, termID)
	if err != nil && err.Error() == mongo.ErrNoDocuments.Error() {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Không tìm thấy học kỳ!", nil)
		return
	}
	err = term.Delete(DB, termID)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Xóa học kỳ không thành công!", nil)
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Xóa học kỳ thành công", nil))
}
func GetTermDetail(c *gin.Context) {
	id := c.Param("id")
	var (
		term Collections.TermModel
		DB   = config.GetMongoDB()
		err  error
		//Other config
		//.......
	)
	res, err := term.FindByID(DB, utils.ConvertStringToObjectId(id))
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Học kỳ không tồn tại!", err.Error())
		return
	}
	c.JSON(
		http.StatusOK,
		Common.SimpleSuccessResponse(http.StatusOK, "Lấy thông tin học kỳ thành công!", res),
	)
}
func GetAllTerms(c *gin.Context) {
	var (
		term Collections.TermModel
		DB   = config.GetMongoDB()
		err  error
		//Other config
		//.......
	)
	page, limit, skip := utils.Pagination(c)
	res, err := term.Find(DB, limit, skip)
	total, _ := term.Count(DB, bson.M{})
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Không thể lấy thông tin!", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.NewSuccessResponse(http.StatusOK, "Lấy danh sách học kỳ thành công", res, int(total), page, limit))
}
