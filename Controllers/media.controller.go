package Controllers

import (
	"fmt"
	"gin-gonic-gom/Middlewares"
	"gin-gonic-gom/Services/media"
	"gin-gonic-gom/common"
	"gin-gonic-gom/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
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
func (mediaController *MediaController) RegisterMediaRoutes(rg *gin.RouterGroup) {
	mediaroute := rg.Group("/upload")
	{
		mediaroute.Use(Middlewares.AuthValidationBearerMiddleware)
		{
			mediaroute.POST("/image", mediaController.UploadImage)
		}
	}
}
