package Collections

import (
	"context"
	"errors"
	"gin-gonic-gom/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type TokenModel struct {
	Id           primitive.ObjectID `bson:"_id"`
	UserId       primitive.ObjectID `bson:"user_id"`
	RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
	AccessToken  string             `json:"access_token" bson:"access_token"`
	Device       string             `json:"device" bson:"device"`
	IpAddress    string             `json:"ip_address" bson:"ip_address"`
	Exp          time.Time          `json:"exp" bson:"exp"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}
type OTPModel struct {
	Id        primitive.ObjectID `bson:"_id"`
	UserId    primitive.ObjectID `bson:"user_id"`
	OTPCode   string             `json:"otp_code" bson:"otp_code"`
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func (t TokenModel) GetCollectionName() string {
	return "Tokens"
}
func (o *OTPModel) GetCollectionNameSecondary() string {
	return "OTPs"
}

type Tokens []TokenModel
type Otps []OTPModel

func (t *TokenModel) FindAndUpdate(DB *mongo.Database, filter interface{}, update interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	_, err := DB.Collection(t.GetCollectionName()).UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}
func (t *TokenModel) Count(DB *mongo.Database, filter interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	if total, err := DB.Collection(t.GetCollectionName()).CountDocuments(ctx, filter, options.Count()); err != nil {
		return 0, err
	} else {
		return total, nil
	}
}
func (t *TokenModel) CheckAndDeleteDevice(DB *mongo.Database, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	filter := bson.M{
		"user_id": id,
	}
	cursor, err := DB.Collection(t.GetCollectionName()).Find(ctx, filter)
	if err != nil {
		return err
	}
	var tokens Tokens
	if err := cursor.All(ctx, &tokens); err != nil {
		return err
	}
	numDevices := len(tokens)
	if numDevices >= 2 {
		devicesToDelete := tokens[:numDevices-2]
		if len(devicesToDelete) > 0 {
			var idsToDelete []primitive.ObjectID
			for _, token := range devicesToDelete {
				idsToDelete = append(idsToDelete, token.Id)
			}
			_, err = (DB.Collection(t.GetCollectionName())).DeleteMany(ctx, bson.M{"_id": bson.M{"$in": idsToDelete}})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (t *TokenModel) CheckExistToken(DB *mongo.Database, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	filter := bson.M{
		"access_token": token,
	}
	if result := DB.Collection(t.GetCollectionName()).FindOne(ctx, filter); result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return errors.New("không tìm thấy token")
		}
		return result.Err()
	} else {
		return result.Decode(&t)
	}
}
func (t *TokenModel) FindByToken(DB *mongo.Database, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	filter := bson.M{
		"access_token": token,
	}
	err := DB.Collection(t.GetCollectionName()).FindOne(ctx, filter).Decode(&t)
	if err != nil {
		return err
	}
	return nil
}
func (t *TokenModel) DeleteOne(DB *mongo.Database, filter interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	_, err := DB.Collection(t.GetCollectionName()).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
func (o *OTPModel) SaveOTP(DB *mongo.Database, userId primitive.ObjectID, otpHash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	updateData := OTPModel{
		Id:        primitive.NewObjectID(),
		UserId:    userId,
		OTPCode:   otpHash,
		ExpiresAt: time.Now().UTC().Add(5 * time.Minute),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	_, err := DB.Collection(o.GetCollectionNameSecondary()).InsertOne(ctx, updateData)
	if err != nil {
		return err
	}
	return nil
}
func (o *OTPModel) FindOne(DB *mongo.Database, filter interface{}) (*OTPModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	err := DB.Collection(o.GetCollectionNameSecondary()).FindOne(ctx, filter).Decode(&o)
	if err != nil {
		return nil, err
	}
	return o, nil
}
func (o *OTPModel) VerifyOTP(DB *mongo.Database, email, otpHash string) (bool, error) {
	var userEntry UserModel
	user, err := userEntry.FindByEmail(DB, email)
	if err != nil {
		return false, err
	}
	filter := bson.M{"$and": bson.A{bson.M{"user_id": user.Id}, bson.M{"otp_code": otpHash}}}
	res, err := o.FindOne(DB, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, errors.New("OTP không đúng!")
		}
		return false, errors.New("Email không tồn tại!")
	}
	timeNowUTC := time.Now().UTC()
	if otpHash != res.OTPCode {
		return false, errors.New("OTP không đúng!")
	}
	if timeNowUTC.After(res.ExpiresAt) {
		return false, errors.New("OTP đã hết hạn!")
	}
	return true, nil
}
func (o *OTPModel) UpdateOne(DB *mongo.Database, filter, update interface{}, opts *options.UpdateOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.CTimeOut)
	defer cancel()
	_, err := DB.Collection(o.GetCollectionNameSecondary()).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}
func (o *OTPModel) ResendOTP(DB *mongo.Database, id primitive.ObjectID, otpHash string) (bool, error) {
	filter := bson.M{"user_id": id}
	_, err := o.FindOne(DB, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, errors.New("Chưa thực hiện gửi OTP!")
		}
		return false, errors.New("Chưa thực hiện gửi OTP!")
	}
	// Bước 4: Cập nhật hoặc tạo mới OTP trong cơ sở dữ liệu
	updateOTPData := bson.D{{
		"$set", bson.D{
			{"user_id", id},
			{"otp_code", otpHash},
			{"expires_at", time.Now().UTC().Add(5 * time.Minute)},
			{"created_at", time.Now()},
			{"updated_at", time.Now()},
		},
	}}
	err = o.UpdateOne(DB, filter, updateOTPData, options.Update().SetUpsert(true))
	if err != nil {
		return false, errors.New("Không thể cập nhật OTP trong cơ sở dữ liệu!")
	}
	return true, nil
}
