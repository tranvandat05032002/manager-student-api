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
