package major

import (
	"gin-gonic-gom/Models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MajorService interface {
	CreateMajor(*Models.MajorModel) error
	UpdateMajor(primitive.ObjectID, *Models.MajorUpdateReq) (*Models.MajorModel, error)
	GetAllMajor(int, int) ([]*Models.MajorModel, int, error)
	GetMajorDetails(primitive.ObjectID) (*Models.MajorModel, error)
	DeleteMajor(primitive.ObjectID) (int, error)
}
