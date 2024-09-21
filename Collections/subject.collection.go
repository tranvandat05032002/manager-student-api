package Collections

import (
	"context"
	"fmt"
	"gin-gonic-gom/config"
	"gin-gonic-gom/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
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
type Subjects []SubjectModel

func (s SubjectModel) GetCollectionName() string {
	return "Subjects"
}
func (s *SubjectModel) CheckExist(DB *mongo.Database, subjectCode string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	// parse query
	filter := bson.M{"subject_code": subjectCode}
	if result := DB.Collection(s.GetCollectionName()).FindOne(ctx, filter); result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			fmt.Println("Running Err --> ", result.Err())
			return false, nil
		}
		return false, result.Err()
	} else {
		return true, result.Decode(&s)
	}
}
func (s *SubjectModel) Create(DB *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	// Parse Query
	s.ID = primitive.NewObjectID()
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	if _, err := DB.Collection(s.GetCollectionName()).InsertOne(ctx, s); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *SubjectModel) Update(DB *mongo.Database, id primitive.ObjectID, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	filter := bson.M{
		"_id": id,
	}
	update := utils.BuildUpdateQuery(data)
	update["updated_at"] = time.Now()
	dataUpdate := bson.M{"$set": update}
	_, err := DB.Collection(s.GetCollectionName()).UpdateOne(ctx, filter, dataUpdate)
	if err != nil {
		return err
	}
	return nil
}
func (s *SubjectModel) Count(DB *mongo.Database, filter interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	if total, err := DB.Collection(s.GetCollectionName()).CountDocuments(ctx, filter, options.Count()); err != nil {
		return 0, err
	} else {
		return total, nil
	}
}
func (s *SubjectModel) Find(DB *mongo.Database, limit, skip int) (Subjects, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cur, err := DB.Collection(s.GetCollectionName()).Find(ctx, bson.M{}, opts)
	defer cur.Close(ctx)
	var subjects Subjects
	for cur.Next(ctx) {
		var subject SubjectModel
		err := cur.Decode(&subject)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, subject)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return subjects, err
}

func (s *SubjectModel) Delete(DB *mongo.Database, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	filter := bson.M{"_id": id}
	_, err := DB.Collection(s.GetCollectionName()).DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("Err ---> ", err)
		return err
	}
	return nil
}
func (s *SubjectModel) FindByID(DB *mongo.Database, id primitive.ObjectID) (*SubjectModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	query := bson.M{"_id": id}
	err := DB.Collection(s.GetCollectionName()).FindOne(ctx, query).Decode(&s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
func (s *SubjectModel) Search(DB *mongo.Database, subjectName string, skip, limit int) (Subjects, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()

	totalCount, err := DB.Collection(s.GetCollectionName()).CountDocuments(ctx, bson.D{{"$text", bson.D{{"$search", subjectName}}}})
	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"$text", bson.D{{"$search", subjectName}}}}}},
		bson.D{{"$skip", skip}},
		bson.D{{"$limit", limit}},
	}
	var subjects Subjects
	cursor, err := DB.Collection(s.GetCollectionName()).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &subjects); err != nil {
		return nil, 0, err
	}
	return subjects, int(totalCount), nil

}
