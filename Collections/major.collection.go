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

type MajorModel struct {
	Id        primitive.ObjectID `bson:"_id"`
	MajorId   string             `json:"major_id" bson:"major_id" binding:"required,min=4,max=20"`
	MajorName string             `json:"major_name" bson:"major_name" binding:"required,min=2,max=100"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
type MajorUpdateReq struct {
	MajorId   string `json:"major_id" bson:"-"`
	MajorName string `json:"major_name" bson:"-"`
}
type Majors []MajorModel

func (m *MajorModel) GetCollectionName() string {
	return "Majors"
}
func (m *MajorModel) CheckExist(DB *mongo.Database, majorId, majorName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	// parse query
	filter := bson.M{
		"$or": []bson.M{
			{"major_id": majorId},
			{"major_name": majorName},
		},
	}
	if result := DB.Collection(m.GetCollectionName()).FindOne(ctx, filter); result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			fmt.Println("Err")
			return false, nil
		}
		return false, result.Err()
	} else {
		return true, result.Decode(&m)
	}
}
func (m *MajorModel) Create(DB *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	// Parse Query
	m.Id = primitive.NewObjectID()
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	if _, err := DB.Collection(m.GetCollectionName()).InsertOne(ctx, m); err != nil {
		return err
	} else {
		return nil
	}
}

func (m *MajorModel) Update(DB *mongo.Database, id primitive.ObjectID, data interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	filter := bson.M{
		"_id": id,
	}
	update := utils.BuildUpdateQuery(data)
	update["updated_at"] = time.Now()
	dataUpdate := bson.M{"$set": update}
	_, err := DB.Collection(m.GetCollectionName()).UpdateOne(ctx, filter, dataUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (m *MajorModel) Find(DB *mongo.Database, limit, skip int) (Majors, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	count := 0
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cur, err := DB.Collection(m.GetCollectionName()).Find(ctx, bson.M{}, opts)
	defer cur.Close(ctx)
	var majors Majors
	for cur.Next(ctx) {
		var major MajorModel
		err := cur.Decode(&major)
		if err != nil {
			return nil, 0, err
		}
		count++
		majors = append(majors, major)
	}
	if err := cur.Err(); err != nil {
		return nil, 0, err
	}
	return majors, count, err
}

func (m *MajorModel) Delete(DB *mongo.Database, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	filter := bson.M{"_id": id}
	_, err := DB.Collection(m.GetCollectionName()).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (m *MajorModel) FindByID(DB *mongo.Database, id primitive.ObjectID) (*MajorModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	query := bson.M{"_id": id}
	err := DB.Collection(m.GetCollectionName()).FindOne(ctx, query).Decode(&m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *MajorModel) Search(DB *mongo.Database, majorName string, skip, limit int) (Majors, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()

	totalCount, err := DB.Collection(m.GetCollectionName()).CountDocuments(ctx, bson.D{{"$text", bson.D{{"$search", majorName}}}})
	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"$text", bson.D{{"$search", majorName}}}}}},
		bson.D{{"$skip", skip}},
		bson.D{{"$limit", limit}},
	}
	var majorList Majors
	cursor, err := DB.Collection(m.GetCollectionName()).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &majorList); err != nil {
		return nil, 0, err
	}
	return majorList, int(totalCount), nil

}
