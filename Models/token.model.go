package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type TokenModel struct {
	Id           primitive.ObjectID `bson:"_id"`
	UserId       primitive.ObjectID `bson:"user_id"`
	RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
	Deviced      string             `json:"deviced" bson:"deviced"`
	IpAddress    string             `json:"ip_address" bson:"ip_address"`
	Exp          time.Time          `json:"exp" bson:"exp"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}
