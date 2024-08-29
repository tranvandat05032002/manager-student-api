package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type StatisticalOfTermRes struct {
	ID               primitive.ObjectID `json:"_id" bson:"_id"`
	EndDate          time.Time          `json:"end_date" bson:"end_date"`
	StartDate        time.Time          `json:"start_date" bson:"start_date"`
	Subjects         []SubjectModel     `json:"subjects" bson:"subjects"`
	TermAcademicYear string             `json:"term_academic_year" bson:"term_academic_year"`
	TermSemester     string             `json:"term_semester" bson:"term_semester"`
	TotalCredits     int                `json:"total_credits" bson:"total_credits"`
	TotalSubjects    int                `json:"total_subjects" bson:"total_subjects"`
}
type StatisticalExportInput struct {
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
	TermAcademicYear string `json:"term_academic_year"`
	TermSemester     string `json:"term_semester"`
	TotalCredits     int    `json:"total_credits"`
	TotalSubjects    int    `json:"total_subjects"`
}
