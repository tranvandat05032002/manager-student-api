package major

import (
	"context"
	"errors"
	"fmt"
	"gin-gonic-gom/Models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MajorImplementService struct {
	majorcollection *mongo.Collection
	ctx             context.Context
}

func NewMajorService(majorcollection *mongo.Collection, ctx context.Context) MajorService {
	return &MajorImplementService{
		majorcollection: majorcollection,
		ctx:             ctx,
	}
}
func (a *MajorImplementService) MajorExist(majorId, majorName string) (bool, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"major_id": majorId},
			{"major_name": majorName},
		},
	}
	var major Models.MajorModel
	err := a.majorcollection.FindOne(a.ctx, filter).Decode(&major)
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
func (a *MajorImplementService) CreateMajor(major *Models.MajorModel) error {
	exists, err := a.MajorExist(major.MajorId, major.MajorName)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("Mã ngành hoặc tên ngành đã tồn tại trong hệ thống!")
	}
	major.Id = primitive.NewObjectID()
	major.CreatedAt = time.Now().UTC()
	major.UpdatedAt = time.Now().UTC()
	_, err = a.majorcollection.InsertOne(a.ctx, major)
	if err != nil {
		return errors.New("Tạo ngành thất bại")
	}
	return nil
}
func (a *MajorImplementService) UpdateMajor(id primitive.ObjectID, majorUpdate *Models.MajorUpdateReq) (*Models.MajorModel, error) {
	filter := bson.M{
		"_id": id,
	}
	updateFields := bson.D{}
	// fix lỗi nếu update thiếu trường thì sẽ không set giá trị đó là chuỗi rỗng
	if majorUpdate.MajorId != "" {
		updateFields = append(updateFields, bson.E{"major_id", majorUpdate.MajorId})
	}
	if majorUpdate.MajorName != "" {
		updateFields = append(updateFields, bson.E{"major_name", majorUpdate.MajorName})
	}
	updateFields = append(updateFields, bson.E{"updated_at", time.Now().UTC()})

	if len(updateFields) > 0 {
		majorDataUpdate := bson.D{
			{"$set", updateFields},
		}

		var majorRes *Models.MajorModel
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

		err := a.majorcollection.FindOneAndUpdate(a.ctx, filter, majorDataUpdate, opts).Decode(&majorRes)
		if err != nil {
			return nil, err
		}
		return majorRes, nil
	}
	return nil, fmt.Errorf("Hủy update!")
}
func (a *MajorImplementService) GetAllMajor(page, limit int) ([]*Models.MajorModel, int, error) {
	skip := limit * (page - 1)
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))
	cur, err := a.majorcollection.Find(a.ctx, bson.M{}, opts)
	total, err := a.majorcollection.CountDocuments(a.ctx, bson.M{})
	defer cur.Close(a.ctx)
	var majors []*Models.MajorModel
	for cur.Next(a.ctx) {
		var major *Models.MajorModel
		err := cur.Decode(&major)
		if err != nil {
			return nil, 0, err
		}
		majors = append(majors, major)
	}
	if err := cur.Err(); err != nil {
		return nil, 0, err
	}
	return majors, int(total), err
}
func (a *MajorImplementService) DeleteMajor(id primitive.ObjectID) (int, error) {
	filter := bson.M{"_id": id}
	res, err := a.majorcollection.DeleteOne(a.ctx, filter)
	return int(res.DeletedCount), err
}
func (a *MajorImplementService) GetMajorDetails(id primitive.ObjectID) (*Models.MajorModel, error) {
	var major *Models.MajorModel
	query := bson.M{"_id": id}
	err := a.majorcollection.FindOne(a.ctx, query).Decode(&major)
	return major, err
}
