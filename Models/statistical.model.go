package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatisticalOfTermRes struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Subjects      []SubjectModel     `json:"subjects" bson:"subjects"`
	TermSemester  int                `json:"term_semester" bson:"term_semester"`
	TermFromYear  int                `json:"term_from_year" bson:"term_from_year"`
	TermToYear    int                `json:"term_to_year" bson:"term_to_year"`
	TotalCredits  int                `json:"total_credits" bson:"total_credits"`
	TotalSubjects int                `json:"total_subjects" bson:"total_subjects"`
}
type StatisticalExportInput struct {
	TermSemester  int `json:"term_semester"`
	TermFromYear  int `json:"term_from_year"`
	TermToYear    int `json:"term_to_year"`
	TotalCredits  int `json:"total_credits"`
	TotalSubjects int `json:"total_subjects"`
}
