package Controllers

import (
	"gin-gonic-gom/Collections"
	"gin-gonic-gom/Common"
	"gin-gonic-gom/Models"
	"gin-gonic-gom/config"
	"gin-gonic-gom/constant"
	"gin-gonic-gom/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxFileSize = 5 * 1024 * 1024 // 5 MB
)

func UploadImage(c *gin.Context) {
	var (
		media Collections.MediaModel
		err   error
		DB    = config.GetMongoDB()
	)
	file, err := c.FormFile("file")
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "File không hợp lệ", nil)
		return
	}
	ext := filepath.Ext(file.Filename)
	if len(c.Request.MultipartForm.File["file"]) > 1 {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Chỉ cho phép upload một file duy nhất", nil)
		return
	} else if len(c.Request.MultipartForm.File["file"]) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Không được để trống",
		})
		return
	}
	if file.Size > maxFileSize {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Kích thước file cho phép tối đa là 5MB", nil)
		return
	}

	if !utils.IsAllowedImageExt(ext) {
		Common.NewErrorResponse(c, http.StatusBadRequest, "File không đúng định dạng", nil)
		return
	}
	// custom path
	newFileName := uuid.New().String()
	newFileNameWithExt := newFileName + ext
	path := filepath.Join("uploads/images", newFileNameWithExt)
	// lưu file vào trong uploads/images
	if c.SaveUploadedFile(file, path); err != nil {
		Common.NewErrorResponse(c, http.StatusInternalServerError, "Đã xảy ra lỗi hệ thông khi lưu file upload", nil)
		return
	}
	// Tạo URL cho hình ảnh đã upload
	imageURL := utils.GeneratorURLImage(c, newFileNameWithExt)
	err = media.Upload(DB, imageURL)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, Common.ErrorShouldBindDataMessage, err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Upload Success!", imageURL))
}
func UploadUserByExcel(c *gin.Context) {
	var (
		media Collections.MediaModel
		err   error
		//DB    = config.GetMongoDB()
	)
	timeLocalHoChiMinh, _ := utils.GetCurrentTimeInLocal("Asia/Ho_Chi_Minh")
	file, err := c.FormFile("file")
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Không được để trống", err.Error())
		return
	}
	// kiểm tra định dạng file có phải là xlsm và xlsx không
	if !strings.HasSuffix(file.Filename, ".xlsx") && !strings.HasSuffix(file.Filename, ".xlsm") {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Chỉ chấp nhận định dạng .xlsm và .xlsx", nil)
		return
	}
	// Lưu file tạm thời vào local
	filePath := filepath.Join(os.TempDir(), file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Đã xảy ra lỗi trong quá trình lưu file", nil)
		return
	}
	// Đọc file excel
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Đã xảy ra lỗi trong quá trình đọc file", err.Error())
		return
	}
	defer func() {
		// Xóa file và kiểm tra lỗi
		if err := os.Remove(filePath); err != nil {
			Common.NewErrorResponse(c, http.StatusInternalServerError, "Đã xảy ra lỗi hệ thống trong quá trình xóa file", nil)
			return
		}
	}()
	// Lấy sheet đầu tiên
	sheetNames := f.GetSheetList()
	if len(sheetNames) == 0 {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Không tìm thấy sheet trong file", nil)
		return
	}
	sheetName := sheetNames[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Đã xảy ra lỗi khi lấy dữ liệu từ sheet", err.Error())
		return
	}
	var usersList []Models.UserModel
	avatar := "https://images2.thanhnien.vn/528068263637045248/2024/1/25/c3c8177f2e6142e8c4885dbff89eb92a-65a11aeea03da880-1706156293184503262817.jpg"
	// Bỏ qua dòng đầu tiên (dòng tiêu đề) sau đó Gán dữ liệu của excel vào model
	for i, record := range rows {
		if i == 0 {
			continue
		}
		if len(record) >= 4 {
			newUser := Models.UserModel{
				Id:             primitive.NewObjectID(),
				MajorId:        nil,
				Role:           constant.STUDENT, // default là student khi tạo bằng excel
				Name:           record[0],
				Email:          record[1],
				Password:       record[2],
				Avatar:         avatar,
				Phone:          record[3],
				Address:        "",
				Department:     "",
				Gender:         nil,
				DateOfBirth:    timeLocalHoChiMinh,
				EnrollmentDate: timeLocalHoChiMinh,
				HireDate:       timeLocalHoChiMinh,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}
			usersList = append(usersList, newUser)
		}
	}
	err = media.InsertManyUser(usersList)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusBadRequest, "Đã xảy ra lỗi thêm dữ liệu, vui lòng sửa rồi thêm lại", err.Error())
		return
	}
	c.JSON(http.StatusOK, Common.SimpleSuccessResponse(http.StatusOK, "Đọc file excel thành công!", nil))
}
