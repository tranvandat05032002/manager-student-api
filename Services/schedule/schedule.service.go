package schedule

import "gin-gonic-gom/Models"

type ScheduleService interface {
	CreateSchedule(model Models.ScheduleModel) error
}
