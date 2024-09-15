package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// format: Thứ 2 [9-12, E401 - Lab]
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
