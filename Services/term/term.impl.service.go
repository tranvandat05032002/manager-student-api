package term

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
func (a *TermImplementService) CheckTermExist(term_semester, fromYear, toYear int) (bool, error) {
	filter := bson.M{
		"$and": []bson.M{
			{"term_semester": term_semester},
			{"term_from_year": fromYear},
			{"term_to_year": toYear},
		},
	}
	var term Models.TermModel
	err := a.termcollection.FindOne(a.ctx, filter).Decode(&term)
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
	exists, err := a.CheckTermExist(termInput.TermSemester, termInput.TermFromYear, termInput.TermToYear)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("Học kỳ của năm đã tồn tại trong hệ thống!")
	}
	//check và trả về mã học kỳ và set vào trong TermID
	// create một subject
	termData := Models.TermModel{
		ID:           primitive.NewObjectID(),
		TermSemester: termInput.TermSemester,
		TermFromYear: termInput.TermFromYear,
		TermToYear:   termInput.TermToYear,
		CreatedAt:    timeHoChiMinhLocal,
		UpdatedAt:    timeHoChiMinhLocal,
	}
	_, err = a.termcollection.InsertOne(a.ctx, termData)
	if err != nil {
		return errors.New("Tạo học kỳ thất bại")
	}
	return nil
}
func (a *TermImplementService) GetTermDetails(id primitive.ObjectID) (Models.TermModel, error) {
	var term Models.TermModel
	query := bson.M{"_id": id}
	err := a.termcollection.FindOne(a.ctx, query).Decode(&term)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return term, fmt.Errorf("no document found")
		}
		return term, err
	}
	return term, err
}
func (a *TermImplementService) UpdateTerm(id primitive.ObjectID, termUpdate Models.TermInput) (Models.TermModel, error) {
	timeHoChiMinhLocal, _ := utils.GetCurrentTimeInLocal("Asia/Ho_Chi_Minh")
	filter := bson.M{
		"_id": id,
	}
	termData := utils.BuildUpdateQuery(termUpdate)
	termData["updated_at"] = timeHoChiMinhLocal
	termDataUpdate := bson.M{"$set": termData}
	var termRes Models.TermModel
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := a.termcollection.FindOneAndUpdate(a.ctx, filter, termDataUpdate, opts).Decode(&termRes)
	if err != nil {
		return termRes, err
	}
	return termRes, nil
}

func (a *TermImplementService) GetAllTerm(page, limit int) ([]Models.TermModel, int, error) {
	skip := limit * (page - 1)
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cur, err := a.termcollection.Find(a.ctx, bson.M{}, opts)
	total, err := a.termcollection.CountDocuments(a.ctx, bson.M{})
	defer cur.Close(a.ctx)
	var terms []Models.TermModel
	for cur.Next(a.ctx) {
		var term Models.TermModel
		err := cur.Decode(&term)
		if err != nil {
			return nil, 0, err
		}
		terms = append(terms, term)
	}
	if err := cur.Err(); err != nil {
		return nil, 0, err
	}
	return terms, int(total), err
}
func (a *TermImplementService) DeleteTerm(id primitive.ObjectID) (int, error) {
	filter := bson.M{"_id": id}
	res, err := a.termcollection.DeleteOne(a.ctx, filter)
	if err != nil {
		return 0, err
	}
	return int(res.DeletedCount), err
}
