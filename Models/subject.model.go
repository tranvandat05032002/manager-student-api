package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SubjectModel struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	SubjectCode  string             `json:"subject_code" bson:"subject_code"` // Mã học phần
	SubjectName  string             `json:"subject_name" bson:"subject_name"`
	Credits      int                `json:"credits" bson:"credits"`
	IsMandatory  bool               `json:"is_mandatory" bson:"is_mandatory"` // Học phần bắt buộc
	TermID       primitive.ObjectID `json:"term_id" bson:"term_id"`
	AcademicYear string             `json:"academic_year" bson:"academic_year"` // Học phần được dạy trong năm nào
	Department   string             `json:"department" bson:"department"`       // Khoa
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

type SubjectInput struct {
	SubjectCode  string `json:"subject_code"` // Mã học phần
	SubjectName  string `json:"subject_name"`
	Credits      int    `json:"credits"`
	IsMandatory  bool   `json:"is_mandatory"` // Học phần bắt buộc
	Term         string `json:"term"`
	AcademicYear string `json:"academic_year"`
	Department   string `json:"department"` // Khoa
}
type TermModel struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	TermSemester     string             `json:"term_semester" bson:"term_semester"`
	TermAcademicYear string             `json:"term_academic_year" bson:"term_academic_year"`
	StartDate        time.Time          `json:"start_date" bson:"start_date"`
	EndDate          time.Time          `json:"end_date" bson:"end_date"`
	TotalCredits     int                `json:"total_credits" bson:"total_credits"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}
type TermInput struct {
	TermSemester     string    `json:"term_semester"`
	TermAcademicYear string    `json:"term_academic_year"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
}
