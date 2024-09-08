package Controllers

import (
	"fmt"
	"gin-gonic-gom/Middlewares"
	"gin-gonic-gom/Models"
	"gin-gonic-gom/Services/media"
	"gin-gonic-gom/common"
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

type MediaController struct {
	MediaService media.MediaService
}

func NewMedia(mediaService media.MediaService) MediaController {
	return MediaController{
		MediaService: mediaService,
	}
}
func GeneratorURLImage(ctx *gin.Context, fileNameExt string) string {
	imageURL := "http://" + ctx.Request.Host + "/static/images/" + fileNameExt
	return imageURL
}
func (mediaController *MediaController) UploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	ext := filepath.Ext(file.Filename)
	if len(ctx.Request.MultipartForm.File["file"]) > 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Chỉ cho phép upload một file duy nhất",
		})
		return
	} else if len(ctx.Request.MultipartForm.File["file"]) < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Không được để trống",
		})
		return
	}
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintln("Get form error: %s", err.Error()))
		return
	}
	if file.Size > maxFileSize {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Kích thước file cho phép tối đa là 5MB",
		})
		return
	}

	if !utils.IsAllowedImageExt(ext) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format"})
		return
	}
	// custom path
	newFileName := uuid.New().String()
	newFileNameWithExt := newFileName + ext
	path := filepath.Join("uploads/images", newFileNameWithExt)
	// lưu file vào trong uploads/images
	if ctx.SaveUploadedFile(file, path); err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintln("Upload image error: %s", err.Error()))
		return
	}
	// Tạo URL cho hình ảnh đã upload
	imageURL := GeneratorURLImage(ctx, newFileNameWithExt)
	err = mediaController.MediaService.Upload(imageURL)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, common.ErrorShouldBindDataMessage, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "Upload Success!", imageURL))
}
func (mediaController *MediaController) UploadUserExcel(ctx *gin.Context) {
	timeLocalHoChiMinh, _ := utils.GetCurrentTimeInLocal("Asia/Ho_Chi_Minh")
	file, err := ctx.FormFile("file")
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Không được để trống", err.Error())
		return
	}
	// kiểm tra định dạng file có phải là xlsm và xlsx không
	if !strings.HasSuffix(file.Filename, ".xlsx") && !strings.HasSuffix(file.Filename, ".xlsm") {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Chỉ chấp nhận định dạng .xlsm và .xlsx", err.Error())
		return
	}
	// Lưu file tạm thời vào local
	filePath := filepath.Join(os.TempDir(), file.Filename)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Đã xảy ra lỗi trong quá trình lưu file", err.Error())
		return
	}
	// Đọc file excel
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Đã xảy ra lỗi trong quá trình đọc file", err.Error())
		return
	}
	defer os.Remove(filePath)
	// Lấy sheet đầu tiên
	sheetNames := f.GetSheetList()
	if len(sheetNames) == 0 {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Không tìm thấy sheet trong file", nil)
		return
	}
	sheetName := sheetNames[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Đã xảy ra lỗi khi lấy dữ liệu từ sheet", err.Error())
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
	err = mediaController.MediaService.UploadExcelDataUser(usersList)
	if err != nil {
		common.NewErrorResponse(ctx, http.StatusBadRequest, "Đã xảy ra lỗi thêm dữ liệu, vui lòng sửa rồi thêm lại", err.Error())
		return
	}
	ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(http.StatusOK, "Đọc file excel thành công!", nil))
}
func (mediaController *MediaController) RegisterMediaRoutes(rg *gin.RouterGroup) {
	mediaroute := rg.Group("/upload")
	{
		mediaroute.Use(Middlewares.AuthValidationBearerMiddleware)
		{
			mediaroute.POST("/image", mediaController.UploadImage)
		}
	}
	adminmediaroute := rg.Group("/admin/excel/upload")
	{
		adminmediaroute.Use(Middlewares.AuthValidationBearerMiddleware)
		adminmediaroute.Use(Middlewares.RoleMiddleware("admin"))
		{
			adminmediaroute.POST("/user", mediaController.UploadUserExcel)
		}
	}
}
