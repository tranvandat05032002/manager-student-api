package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MajorModel struct {
	Id        primitive.ObjectID `bson:"_id"`
	MajorId   string             `json:"major_id" bson:"major_id" binding:"required,min=4,max=20"`
	MajorName string             `json:"major_name" bson:"major_name" binding:"required,min=2,max=100"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
type MajorUpdateReq struct {
	MajorId   string `json:"major_id" bson:"-"`
	MajorName string `json:"major_name" bson:"-"`
}
