package utils

import (
	"fmt"
	"gin-gonic-gom/Models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateAccessToken(exp time.Duration, user *Models.UserModel, secretJWTKey string) (string, error) {
	now := time.Now().UTC()
	expirationTimeUTCPlus7 := ConvertDurationToTimeUTC(exp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.Id,
		"exp":    expirationTimeUTCPlus7.Unix(),
		"iat":    now.Unix(),
		"role":   user.Role,
	})
	return token.SignedString([]byte(secretJWTKey))
}
func GenerateRefreshToken(exp time.Duration, user *Models.UserModel, secretJWTKey string) (string, error) {
	now := time.Now().UTC()
	expirationTimeUTCPlus7 := ConvertDurationToTimeUTC(exp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.Id,
		"exp":    expirationTimeUTCPlus7.Unix(),
		"iat":    now.Unix(),
		"role":   user.Role,
	})
	return token.SignedString([]byte(secretJWTKey))
}

type TokenType struct {
	UserId string    `json:"userID"`
	Exp    time.Time `json:"exp"`
	Iat    time.Time `json:"iat"`
}

func DecodedToken(tokenString string, secretkey []byte) (jwt.MapClaims, error) {
	//claims := &TokenType{}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretkey, nil
	})
	//fmt.Println("UserId: ", claims.UserId)
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Errorf("Claims Error")
	}
	return claims, nil
	//return claims, nil

	//token, err := jwt.Parse("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0dGwiOjE3MzIzMzY1OTMsInVzZXJJRCI6IjY2YmFlMGQyMDdiOGI5NTUyOTY1ZGNkMyJ9.y1Yyhz3IBOyzrqZrOH-ERfBUNpFxul6igHZ3BNNixsI", func(token *jwt.Token) (interface{}, error) {
	//	// Don't forget to validate the alg is what you expect:
	//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	//		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	//	}
	//
	//	// TODO: Move this to env variable.
	//
	//	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
	//	return secretkey, nil
	//})
	//if err != nil {
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	//	c.Abort()
	//	return
	//}

	//claims, ok := token.Claims.(jwt.MapClaims)
	//if !ok {
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	//	c.Abort()
	//	return
	//}
}
