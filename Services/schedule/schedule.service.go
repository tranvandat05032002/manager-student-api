package schedule

import (
	"gin-gonic-gom/Models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScheduleService interface {
	CreateSchedule(Models.ScheduleModel) error
	GetScheduleDetails(primitive.ObjectID) (*Models.ScheduleModel, error)
	UpdateSchedule(primitive.ObjectID, Models.ScheduleModel) error
	GetAllSchedule(int, int) ([]Models.ScheduleModel, int, error)
	DeleteSchedule(primitive.ObjectID) (int, error)
}
