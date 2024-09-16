package Middlewares

import (
	"context"
	"gin-gonic-gom/Models"
	"gin-gonic-gom/common"
	"gin-gonic-gom/config"
	"gin-gonic-gom/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"os"
	"strings"
	"time"
)

var ctx context.Context

func AuthValidationBearerMiddleware(c *gin.Context) {
	var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
	tokenCollection := config.GetMongoDB().Collection("Tokens")
	//authHeader := c.GetHeader("Authorization")
	var token = c.Request.Header.Get("Authorization")
	deviced := c.Request.Header.Get("User-Agent")
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
	accessSecretKeyByte := []byte(accessSecretKey)
	claims, err := utils.DecodedToken(tokenString, accessSecretKeyByte)
	filterToken := bson.M{"access_token": tokenString}
	var tokenRes Models.TokenModel
	err = tokenCollection.FindOne(ctx, filterToken).Decode(&tokenRes)
	if err != nil {
		common.NewErrorResponse(c, http.StatusUnauthorized, "Decoded thất bại!", err.Error())
		c.Abort()
		return
	}
	c.Set("deviced", deviced)
	c.Set("userId", claims["userID"].(string))
	c.Set("role", claims["role"].(string))
	c.Next()
}

func RoleMiddleware(RoleRequired string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			common.NewErrorResponse(c, http.StatusNotFound, common.ErrorFindUser, "")
			c.Abort()
			return
		}
		if role != RoleRequired {
			common.NewErrorResponse(c, http.StatusForbidden, "Không đủ quyền truy cập routes!", "")
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
//			common.NewErrorResponse(c, status, "An error occurred", err.Error())
//			return
//		}
//
//		// Nếu không có lỗi, nhưng HTTP status khác 200 (ví dụ như 404 Not Found)
//		if c.Writer.Status() != http.StatusOK {
//			status := c.Writer.Status()
//			common.NewErrorResponse(c, status, http.StatusText(status), "Request failed")
//		}
//	}
//}
