package Models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubjectModel struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	SubjectCode string             `json:"subject_code" bson:"subject_code" binding:"required,min=4,max=20"` // Mã học phần
	SubjectName string             `json:"subject_name" bson:"subject_name" binding:"required,min=2,max=100"`
	Credits     int                `json:"credits" bson:"credits" binding:"required,gt=0"`
	IsMandatory bool               `json:"is_mandatory" bson:"is_mandatory"` // Học phần bắt buộc
	TermID      primitive.ObjectID `json:"term_id" bson:"term_id" binding:"required"`
	Department  string             `json:"department" bson:"department" binding:"required,min=2,max=100"` // Khoa
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type SubjectInput struct {
	TermID 		primitive.ObjectID `json:"term_id" binding:"required"`
	SubjectCode  string `json:"subject_code"` // Mã học phần
	SubjectName  string `json:"subject_name"`
	Credits      int    `json:"credits"`
	IsMandatory  bool   `json:"is_mandatory"` // Học phần bắt buộc
	Department   string `json:"department"` // Khoa
}
type TermModel struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	TermSemester int                `json:"term_semester" bson:"term_semester" binding:"required,oneof=1 2 3"`
	TermFromYear int                `json:"term_from_year" bson:"term_from_year" binding:"required,gte=1900,lte=2100"`
	TermToYear   int                `json:"term_to_year" bson:"term_to_year" binding:"required,gte=1900,lte=2100"`
	TotalCredits int                `json:"total_credits" bson:"total_credits"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}
type TermInput struct {
	TermSemester int `json:"term_semester"`
	TermFromYear int `json:"term_from_year"`
	TermToYear   int `json:"term_to_year"`
}
