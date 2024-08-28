package Models

import "go.mongodb.org/mongo-driver/bson/primitive"

type MediaModel struct {
	Id  primitive.ObjectID `bson:"_id"`
	Url string             `json:"url",bson:"url",binding:"required"`
}
