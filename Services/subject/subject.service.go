package subject

import (
	"gin-gonic-gom/Models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubjectService interface {
	CreateSubject(Models.SubjectInput) error
	GetSubjectDetails(primitive.ObjectID) (*Models.SubjectModel, error)
	UpdateSubject(primitive.ObjectID, Models.SubjectInput) (Models.SubjectModel, error)
	GetAllSubject(int, int) ([]Models.SubjectModel, int, error)
	DeleteSubject(primitive.ObjectID) (int, error)
}
