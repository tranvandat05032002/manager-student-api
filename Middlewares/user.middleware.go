package Middlewares

import (
	"context"
	"gin-gonic-gom/Collections"
	"gin-gonic-gom/Common"
	"gin-gonic-gom/config"
	"gin-gonic-gom/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
	"time"
)

var ctx context.Context

func AuthValidationBearerMiddleware(c *gin.Context) {
	var (
		DB         = config.GetMongoDB()
		tokenEntry Collections.TokenModel
		err        error
	)
	var token = c.Request.Header.Get("Authorization")
	device := c.Request.Header.Get("User-Agent")
	if token == "" {
		cookie, _ := c.Request.Cookie("access_token")
		token = cookie.String()
		if len(token) > 6 {
			token = token[6:]
		}
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Hết phiên đăng nhập",
			})
			return
		}
	}
	if len(token) > 6 {
		token = token[6:]
	}
	tokenString := strings.TrimSpace(token)
	accessSecretKey := os.Getenv("ACCESS_TOKEN_SECRET")
	userToken, err := utils.DecodedToken(tokenString, accessSecretKey)
	if err != nil {
		Common.NewErrorResponse(c, http.StatusUnauthorized, "Decoded thất bại!", err.Error())
		c.Abort()
		return
	}
	err = tokenEntry.CheckExistToken(DB, tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Hết phiên đăng nhập!",
		})
		return
	}
	if userToken.Exp.Unix() < time.Now().Unix() {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Hết phiên đăng nhập!",
		})
		return
	}
	c.Set("device", device)
	c.Set("userId", userToken.UserId)
	c.Set("role", userToken.Role)
	c.Next()
}

func RoleMiddleware(RoleRequired string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			Common.NewErrorResponse(c, http.StatusNotFound, Common.ErrorFindUser, "")
			c.Abort()
			return
		}
		if role != RoleRequired {
			Common.NewErrorResponse(c, http.StatusForbidden, "Không đủ quyền truy cập routes!", "")
			c.Abort()
			return
		}
		c.Next()
	}
}

//func ErrorHandlerMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		c.Next()
//
//		// Kiểm tra xem có lỗi nào được lưu trong Gin context không
//		if len(c.Errors) > 0 {
//			// Lấy lỗi cuối cùng
//			err := c.Errors.Last()
//
//			// Xác định mã lỗi HTTP
//			status := c.Writer.Status()
//			fmt.Println("StatusCode ---> ", status)
//			if status == http.StatusOK {
//				status = http.StatusInternalServerError // Gán mặc định là 500 nếu chưa có lỗi HTTP cụ thể
//			}
//
//			// Gửi phản hồi lỗi
//			Common.NewErrorResponse(c, status, "An error occurred", err.Error())
//			return
//		}
//
//		// Nếu không có lỗi, nhưng HTTP status khác 200 (ví dụ như 404 Not Found)
//		if c.Writer.Status() != http.StatusOK {
//			status := c.Writer.Status()
//			Common.NewErrorResponse(c, status, http.StatusText(status), "Request failed")
//		}
//	}
//}
