package media

import (
	"context"
	"fmt"
	"gin-gonic-gom/Models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MediaImplementService struct {
	mediacollection *mongo.Collection
	ctx             context.Context
}

func NewMediaService(mediacollection *mongo.Collection, ctx context.Context) MediaService {
	return &MediaImplementService{
		mediacollection: mediacollection,
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
