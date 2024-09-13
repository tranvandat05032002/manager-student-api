package media

import (
	"context"
	"errors"
	"fmt"
	"gin-gonic-gom/Models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MediaImplementService struct {
	mediacollection *mongo.Collection
	usercollection  *mongo.Collection
	ctx             context.Context
}

func NewMediaService(mediacollection *mongo.Collection, usercollection *mongo.Collection, ctx context.Context) MediaService {
	return &MediaImplementService{
		mediacollection: mediacollection,
		usercollection:  usercollection,
		ctx:             ctx,
	}
}

func (a *MediaImplementService) Upload(path string) error {
	MediaData := Models.MediaModel{
		Id:  primitive.NewObjectID(),
		Url: path,
	}
	_, err := a.mediacollection.InsertOne(a.ctx, &MediaData)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}
func (a *MediaImplementService) UploadExcelDataUser(userList []Models.UserModel) error {
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
		err := a.usercollection.FindOne(a.ctx, filter).Decode(&existsUser)
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
		_, err := a.usercollection.InsertMany(a.ctx, usersInsertInterface)
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
