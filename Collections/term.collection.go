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

type TermModel struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	TermSemester int                `json:"term_semester" bson:"term_semester" binding:"required,oneof=1 2 3"`
	TermFromYear int                `json:"term_from_year" bson:"term_from_year" binding:"required,gte=1900,lte=2100"`
	TermToYear   int                `json:"term_to_year" bson:"term_to_year" binding:"required,gte=1900,lte=2100"`
	TotalCredits int                `json:"total_credits" bson:"total_credits"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}
type Terms []TermModel

func (t TermModel) GetCollectionName() string {
	return "Terms"
}
func (t *TermModel) CheckExist(DB *mongo.Database, termSemester, fromYear, toYear int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	// parse query
	filter := bson.M{
		"$and": []bson.M{
			{"term_semester": termSemester},
			{"term_from_year": fromYear},
			{"term_to_year": toYear},
		},
	}
	if result := DB.Collection(t.GetCollectionName()).FindOne(ctx, filter); result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			fmt.Println("Running Err --> ", result.Err())
			return false, nil
		}
		return false, result.Err()
	} else {
		return true, result.Decode(&t)
	}
}
func (t *TermModel) Create(DB *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	// Parse Query
	t.ID = primitive.NewObjectID()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	if _, err := DB.Collection(t.GetCollectionName()).InsertOne(ctx, t); err != nil {
		return err
	} else {
		return nil
	}
}

func (t *TermModel) Update(DB *mongo.Database, id primitive.ObjectID, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	filter := bson.M{
		"_id": id,
	}
	update := utils.BuildUpdateQuery(data)
	update["updated_at"] = time.Now()
	dataUpdate := bson.M{"$set": update}
	_, err := DB.Collection(t.GetCollectionName()).UpdateOne(ctx, filter, dataUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (t *TermModel) Find(DB *mongo.Database, limit, skip int) (Terms, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	count := 0
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cur, err := DB.Collection(t.GetCollectionName()).Find(ctx, bson.M{}, opts)
	defer cur.Close(ctx)
	var terms Terms
	for cur.Next(ctx) {
		var term TermModel
		err := cur.Decode(&term)
		if err != nil {
			return nil, 0, err
		}
		count++
		terms = append(terms, term)
	}
	if err := cur.Err(); err != nil {
		return nil, 0, err
	}
	return terms, count, err
}

func (t *TermModel) Delete(DB *mongo.Database, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	filter := bson.M{"_id": id}
	_, err := DB.Collection(t.GetCollectionName()).DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println("Err ---> ", err)
		return err
	}
	return nil
}
func (t *TermModel) FindByID(DB *mongo.Database, id primitive.ObjectID) (*TermModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	query := bson.M{"_id": id}
	err := DB.Collection(t.GetCollectionName()).FindOne(ctx, query).Decode(&t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
