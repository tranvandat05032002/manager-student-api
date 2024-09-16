package schedule

import (
	"context"
	"errors"
	"fmt"
	"gin-gonic-gom/Models"
	"gin-gonic-gom/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ScheduleImplementService struct {
	schedulecollection *mongo.Collection
	ctx                context.Context
}

func NewScheduleService(schedulecollection *mongo.Collection, ctx context.Context) ScheduleService {
	return &ScheduleImplementService{
		schedulecollection: schedulecollection,
		ctx:                ctx,
	}
}
func (a *ScheduleImplementService) ScheduleExist(room, dayOfWeek, startPeriod, endPeriod string) (bool, error) {
	filter := bson.M{"$and": bson.A{
		bson.M{"room": room},
		bson.M{"day_of_week": dayOfWeek},
		bson.M{"start_period": startPeriod},
		bson.M{"end_period": endPeriod},
	}}
	var schedule Models.ScheduleModel
	err := a.schedulecollection.FindOne(a.ctx, filter).Decode(&schedule)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
func (a *ScheduleImplementService) CreateSchedule(scheduleInput Models.ScheduleModel) error {
	timeHoChiMinhLocal, _ := utils.GetCurrentTimeInLocal("Asia/Ho_Chi_Minh")
	exists, err := a.ScheduleExist(scheduleInput.Room, scheduleInput.DayOfWeek, scheduleInput.StartPeriod, scheduleInput.EndPeriod)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("Lịch học đã tồn tại trong hệ thống!")
	}
	scheduleData := Models.ScheduleModel{
		ID:          primitive.NewObjectID(),
		TermID:      scheduleInput.TermID,
		DayOfWeek:   scheduleInput.DayOfWeek,
		Room:        scheduleInput.Room,
		StartPeriod: scheduleInput.StartPeriod,
		EndPeriod:   scheduleInput.EndPeriod,
		CreatedAt:   timeHoChiMinhLocal,
		UpdatedAt:   timeHoChiMinhLocal,
	}
	_, err = a.schedulecollection.InsertOne(a.ctx, scheduleData)
	if err != nil {
		fmt.Println("Error --> ", err)
		return errors.New("Tạo lịch học thất bại")
	}
	return nil
}
func (a *ScheduleImplementService) GetScheduleDetails(scheduleId primitive.ObjectID) (*Models.ScheduleModel, error) {
	var schedule *Models.ScheduleModel
	query := bson.M{"_id": scheduleId}
	err := a.schedulecollection.FindOne(a.ctx, query).Decode(&schedule)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no document found")
		}
		return nil, err
	}
	return schedule, err
}
func (a *ScheduleImplementService) UpdateSchedule(scheduleId primitive.ObjectID, scheduleUpdate Models.ScheduleModel) error {
	timeHoChiMinhLocal, _ := utils.GetCurrentTimeInLocal("Asia/Ho_Chi_Minh")
	filter := bson.M{
		"_id": scheduleId,
	}
	scheduleData := utils.BuildUpdateQuery(scheduleUpdate)
	scheduleData["updated_at"] = timeHoChiMinhLocal
	scheduleDataUpdate := bson.M{"$set": scheduleData}
	_, err := a.schedulecollection.UpdateOne(a.ctx, filter, scheduleDataUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (a *ScheduleImplementService) GetAllSchedule(page, limit int) ([]Models.ScheduleModel, int, error) {
	skip := limit * (page - 1)
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cur, err := a.schedulecollection.Find(a.ctx, bson.M{}, opts)
	total, err := a.schedulecollection.CountDocuments(a.ctx, bson.M{})
	defer cur.Close(a.ctx)
	var schedules Models.Schedules
	for cur.Next(a.ctx) {
		var schedule Models.ScheduleModel
		err := cur.Decode(&schedule)
		if err != nil {
			return nil, 0, err
		}
		schedules = append(schedules, schedule)
	}
	if err := cur.Err(); err != nil {
		return nil, 0, err
	}
	return schedules, int(total), err
}
func (a *ScheduleImplementService) DeleteSchedule(scheduleId primitive.ObjectID) (int, error) {
	filter := bson.M{"_id": scheduleId}
	res, err := a.schedulecollection.DeleteOne(a.ctx, filter)
	if err != nil {
		return 0, err
	}
	return int(res.DeletedCount), err
}
