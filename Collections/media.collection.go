package Collections

import (
	"context"
	"errors"
	"fmt"
	"gin-gonic-gom/Models"
	"gin-gonic-gom/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type MediaModel struct {
	Id        primitive.ObjectID `bson:"_id"`
	Url       string             `json:"url" bson:"url" binding:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func (m *MediaModel) GetCollectionName() string {
	return "Medias"
}
func (m *MediaModel) Upload(DB *mongo.Database, URL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	// Parse Query
	m.Id = primitive.NewObjectID()
	m.Url = URL
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	if _, err := DB.Collection(m.GetCollectionName()).InsertOne(ctx, m); err != nil {
		return err
	} else {
		return nil
	}
}

func (m *MediaModel) InsertManyUser(userList []Models.UserModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	if len(userList) == 0 {
		return errors.New("Danh sách rỗng")
	}
	var usersInsertInterface []interface{}
	var duplicateEmails []string
	duplicateCount := 0
	for _, user := range userList {
		var existsUser Models.UserModel
		filter := bson.D{
			{"$or",
				bson.A{
					bson.D{{"email", user.Email}},
					bson.D{{"phone", user.Phone}},
				}},
		}
		err := config.GetMongoDB().Collection("Users").FindOne(ctx, filter).Decode(&existsUser)
		if err == mongo.ErrNoDocuments {
			usersInsertInterface = append(usersInsertInterface, user)
		} else if err != nil {
			return err
		} else {
			duplicateMessage := fmt.Sprintf("Email %s đã tồn tại \\n", existsUser.Email) // client replace thành <br>
			duplicateMessage = fmt.Sprintf("Số điện thoại %s đã tồn tại \\n", existsUser.Phone)
			duplicateEmails = append(duplicateEmails, duplicateMessage)
			duplicateCount++
		}
	}
	if duplicateCount == 0 {
		_, err := config.GetMongoDB().Collection("Users").InsertMany(ctx, usersInsertInterface)
		if err != nil {
			return err
		}
	} else {
		//for _, msg := range duplicateEmails {
		//	return errors.New(msg)
		//}
		return errors.New(fmt.Sprintf("%v", duplicateEmails))
	}
	return nil
}
