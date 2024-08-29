package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OTPModel struct {
	Id        primitive.ObjectID `bson:"_id"`
	UserId    primitive.ObjectID `bson:"user_id"`
	OTPCode   string             `json:"otp_code" bson:"otp_code"`
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
type OTPReq struct {
	Email   string `json:"email" bson:"email" binding:"required,email"`
	OTPCode string `json:"otp_code" bson:"otp_code" binding:"required,len=6"`
}
type OTPRes struct {
	OTPCode   string    `json:"otp_code" bson:"otp_code"`
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"`
}
