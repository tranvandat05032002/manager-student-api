package term

import (
	"gin-gonic-gom/Models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TermService interface {
	CreateTerm(Models.TermInput) error
	GetTermDetails(primitive.ObjectID) (Models.TermModel, error)
	UpdateTerm(primitive.ObjectID, Models.TermInput) (Models.TermModel, error)
	GetAllTerm(int, int) ([]Models.TermModel, int, error)
	DeleteTerm(primitive.ObjectID) (int, error)
}
