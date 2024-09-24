package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func GenerateAccessToken(userId primitive.ObjectID, role string, secretJWTKey string) (string, error) {
	now := time.Now().UTC()
	//expirationTimeUTCPlus7 := ConvertDurationToTimeUTC(exp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userId,
		"exp":    time.Now().Add(2 * time.Minute).Unix(),
		"iat":    now.Unix(),
		"role":   role,
	})
	return token.SignedString([]byte(secretJWTKey))
}
func GenerateRefreshToken(userId primitive.ObjectID, role string, secretJWTKey string) (string, error) {
	now := time.Now().UTC()
	//expirationTimeUTCPlus7 := ConvertDurationToTimeUTC(exp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userId,
		"exp":    time.Now().Add(72 * time.Hour).Unix(),
		"iat":    now.Unix(),
		"role":   role,
	})
	return token.SignedString([]byte(secretJWTKey))
}

type ClaimsType struct {
	UserId string           `json:"userID"`
	Role   string           `json:"role"`
	Exp    *jwt.NumericDate `json:"exp"`
	Iat    *jwt.NumericDate `json:"iat"`
	jwt.RegisteredClaims
}

func DecodedToken(tokenString string, secretkey string) (*ClaimsType, error) {
	claims := &ClaimsType{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretkey), nil
	})
	if token != nil {
		if claims, ok := token.Claims.(*ClaimsType); ok && token.Valid {
			return claims, nil
		}
	}
	return claims, err
}
