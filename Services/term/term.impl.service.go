package term

import (
	"context"
	"errors"
	"gin-gonic-gom/Models"
	"gin-gonic-gom/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TermImplementService struct {
	termcollection *mongo.Collection
	ctx            context.Context
}

func NewTermService(termcollection *mongo.Collection, ctx context.Context) TermService {
	return &TermImplementService{
		termcollection: termcollection,
		ctx:            ctx,
	}
}
func (a *TermImplementService) TermExist(subjectCode string) (bool, error) {
	filter := bson.M{"subject_code": subjectCode}
	var subject Models.SubjectModel
	err := a.subjectcollection.FindOne(a.ctx, filter).Decode(&subject)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
func (a *TermImplementService) CreateTerm(termInput Models.TermInput) error {
	timeHoChiMinhLocal, _ := utils.GetCurrentTimeInLocal("Asia/Ho_Chi_Minh")
	//exists, err := a.SubjectExist(subjectInput.SubjectCode)
	//if err != nil {
	//	return err
	//}
	//if exists {
	//	return errors.New("Mã môn học đã tồn tại trong hệ thống!")
	//}
	// check và trả về mã học kỳ và set vào trong TermID
	// create một subject
	termData := Models.TermModel{
		ID:               primitive.NewObjectID(),
		TermSemester:     termInput.TermSemester,
		TermAcademicYear: termInput.TermAcademicYear,
		StartDate:        termInput.StartDate,
		EndDate:          termInput.EndDate,
		CreatedAt:        timeHoChiMinhLocal,
		UpdatedAt:        timeHoChiMinhLocal,
	}
	_, err = a.termcollection.InsertOne(a.ctx, termData)
	if err != nil {
		return errors.New("Tạo học kỳ thất bại")
	}
	return nil
}
func (a *TermImplementService) GetTermDetails(primitive.ObjectID) (*Models.TermInput, error) {
	return nil, nil
}
func (a *TermImplementService) UpdateTerm(primitive.ObjectID, bson.M) (Models.TermModel, error) {
	return Models.TermModel{}, nil
}

//	func (a *TermImplementService) GetAllTerm(int, int) ([]Models.TermModel, int, error) {
//		return , 0, nil
//	}
func (a *TermImplementService) DeleteTerm(primitive.ObjectID) (int, error) {
	return 0, nil
}
