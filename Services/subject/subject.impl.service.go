package subject

import (
	"context"
	"errors"
	"fmt"
	"gin-gonic-gom/Models"
	"gin-gonic-gom/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SubjectImplementService struct {
	subjectcollection *mongo.Collection
	ctx               context.Context
}

func NewMajorService(subjectcollection *mongo.Collection, ctx context.Context) SubjectService {
	return &SubjectImplementService{
		subjectcollection: subjectcollection,
		ctx:               ctx,
	}
}
func (a *SubjectImplementService) SubjectExist(subjectCode string) (bool, error) {
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
func (a *SubjectImplementService) CreateSubject(subjectInput Models.SubjectInput) error {
	timeHoChiMinhLocal, _ := utils.GetCurrentTimeInLocal("Asia/Ho_Chi_Minh")
	exists, err := a.SubjectExist(subjectInput.SubjectCode)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("Mã môn học đã tồn tại trong hệ thống!")
	}
	// check và trả về mã học kỳ và set vào trong TermID
	// create một subject
	subjectData := Models.SubjectModel{
		ID:           primitive.NewObjectID(),
		SubjectCode:  subjectInput.SubjectCode,
		SubjectName:  subjectInput.SubjectName,
		Credits:      subjectInput.Credits,
		IsMandatory:  subjectInput.IsMandatory,
		Department:   subjectInput.Department,
		AcademicYear: subjectInput.AcademicYear,
		TermID:       primitive.NewObjectID(),
		CreatedAt:    timeHoChiMinhLocal,
		UpdatedAt:    timeHoChiMinhLocal,
	}
	_, err = a.subjectcollection.InsertOne(a.ctx, subjectData)
	if err != nil {
		return errors.New("Tạo môn học thất bại")
	}
	return nil
}

func (a *SubjectImplementService) UpdateSubject(id primitive.ObjectID, subjectUpdate bson.M) (Models.SubjectModel, error) {
	timeHoChiMinhLocal, _ := utils.GetCurrentTimeInLocal("Asia/Ho_Chi_Minh")
	filter := bson.M{
		"_id": id,
	}
	subjectUpdate["updated_at"] = timeHoChiMinhLocal
	subjectDataUpdate := bson.M{"$set": subjectUpdate}
	var subjectRes Models.SubjectModel
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := a.subjectcollection.FindOneAndUpdate(a.ctx, filter, subjectDataUpdate, opts).Decode(&subjectRes)
	if err != nil {
		return subjectRes, err
	}
	return subjectRes, nil
}

func (a *SubjectImplementService) GetAllSubject(page, limit int) ([]Models.SubjectModel, int, error) {
	skip := limit * (page - 1)
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cur, err := a.subjectcollection.Find(a.ctx, bson.M{}, opts)
	total, err := a.subjectcollection.CountDocuments(a.ctx, bson.M{})
	defer cur.Close(a.ctx)
	var subjects []Models.SubjectModel
	for cur.Next(a.ctx) {
		var subject Models.SubjectModel
		err := cur.Decode(&subject)
		if err != nil {
			return nil, 0, err
		}
		subjects = append(subjects, subject)
	}
	if err := cur.Err(); err != nil {
		return nil, 0, err
	}
	return subjects, int(total), err
}
func (a *SubjectImplementService) DeleteSubject(id primitive.ObjectID) (int, error) {
	filter := bson.M{"_id": id}
	res, err := a.subjectcollection.DeleteOne(a.ctx, filter)
	if err != nil {
		return 0, err
	}
	return int(res.DeletedCount), err
}
func (a *SubjectImplementService) GetSubjectDetails(id primitive.ObjectID) (*Models.SubjectModel, error) {
	var subject *Models.SubjectModel
	query := bson.M{"_id": id}
	err := a.subjectcollection.FindOne(a.ctx, query).Decode(&subject)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no document found")
		}
		return nil, err
	}
	return subject, err
}
