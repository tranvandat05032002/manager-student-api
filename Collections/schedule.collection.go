package Collections

import (
	"context"
	"errors"
	"fmt"
	"gin-gonic-gom/config"
	"gin-gonic-gom/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ScheduleModel struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	TermID      primitive.ObjectID `json:"term_id" bson:"term_id" binding:"required"`
	DayOfWeek   string             `json:"day_of_week" bson:"day_of_week" binding:"required"`
	StartPeriod string             `json:"start_period" bson:"start_period"`
	EndPeriod   string             `json:"end_period" bson:"end_period"`
	Room        string             `json:"room" bson:"room" binding:"required"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
type Schedules []ScheduleModel

func (s ScheduleModel) GetCollectionName() string {
	return "Schedules"
}

func (s *ScheduleModel) CheckExist(DB *mongo.Database, room, dayOfWeek, startPeriod, endPeriod string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	// parse query
	filter := bson.M{"$and": bson.A{
		bson.M{"room": room},
		bson.M{"day_of_week": dayOfWeek},
		bson.M{"start_period": startPeriod},
		bson.M{"end_period": endPeriod},
	}}
	if result := DB.Collection(s.GetCollectionName()).FindOne(ctx, filter); result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			fmt.Println("Running Err --> ", result.Err())
			return false, nil
		}
		return false, result.Err()
	} else {
		return true, result.Decode(&s)
	}
}
func (s *ScheduleModel) Create(DB *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	// Parse Query
	s.ID = primitive.NewObjectID()
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	if _, err := DB.Collection(s.GetCollectionName()).InsertOne(ctx, s); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *ScheduleModel) Update(DB *mongo.Database, id primitive.ObjectID, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	filter := bson.M{
		"_id": id,
	}
	errFilter := DB.Collection(s.GetCollectionName()).FindOne(ctx, filter).Decode(&s)
	if errFilter != nil {
		if errFilter == mongo.ErrNoDocuments {
			return errors.New("Không tìm thấy!")
		}
		return errFilter
	}
	update := utils.BuildUpdateQuery(data)
	update["updated_at"] = time.Now()
	dataUpdate := bson.M{"$set": update}
	_, err := DB.Collection(s.GetCollectionName()).UpdateOne(ctx, filter, dataUpdate)
	if err != nil {
		return err
	}
	return nil
}
func (s *ScheduleModel) Count(DB *mongo.Database, filter interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	if total, err := DB.Collection(s.GetCollectionName()).CountDocuments(ctx, filter, options.Count()); err != nil {
		return 0, err
	} else {
		return total, nil
	}
}
func (s *ScheduleModel) Find(DB *mongo.Database, limit, skip int) (Schedules, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cur, err := DB.Collection(s.GetCollectionName()).Find(ctx, bson.M{}, opts)
	defer cur.Close(ctx)
	var schedules Schedules
	for cur.Next(ctx) {
		var schedule ScheduleModel
		err := cur.Decode(&schedule)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return schedules, err
}

func (s *ScheduleModel) Delete(DB *mongo.Database, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	filter := bson.M{"_id": id}
	_, err := DB.Collection(s.GetCollectionName()).DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("Err ---> ", err)
		return err
	}
	return nil
}
func (s *ScheduleModel) FindByID(DB *mongo.Database, id primitive.ObjectID) (*ScheduleModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	query := bson.M{"_id": id}
	err := DB.Collection(s.GetCollectionName()).FindOne(ctx, query).Decode(&s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
