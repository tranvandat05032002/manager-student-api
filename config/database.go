package config

import (
	"context"
	"fmt"
	"gin-gonic-gom/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

const CTimeOut = 10 * time.Second

var (
	err         error
	mongoclient *mongo.Client
	mongoDB     *mongo.Database
)

func Connect(ctx context.Context) (*mongo.Database, error) {
	MongoURL := os.Getenv("MONGODB_URL")
	if MongoURL == "" {
		MongoURL = "mongodb://admin:tranlybuu@103.214.9.124:1236/?authSource=admin"
	}
	mongoconn := options.Client().ApplyURI(os.Getenv("MONGODB_URL"))
	mongoclient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal("error while connecting with mongo", err)
	}
	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("error while trying to ping mongo", err)
	}

	fmt.Println("mongo connection established")
	mongoDB = mongoclient.Database("MeteorDB")
	return mongoclient.Database("MeteorDB"), err
}
func GetMongoDB() *mongo.Database {
	return mongoDB
}
func CloseMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), CTimeOut)
	defer cancel()
	_ = mongoDB.Client().Disconnect(ctx)
}
func createIndex(indexName string, indexType interface{}) mongo.IndexModel {
	if indexType == "text" {
		indexModelText := mongo.IndexModel{Keys: bson.D{{indexName, indexType}}, Options: options.Index().SetDefaultLanguage("none")}
		return indexModelText
	}
	indexModelNotText := mongo.IndexModel{Keys: bson.D{{indexName, indexType}}}
	return indexModelNotText
}
func InitIndex() {
	ctx, cancel := context.WithTimeout(context.Background(), CTimeOut)
	defer cancel()
	//Major index
	indexMajorName := "major_name_text"
	majorCollection := GetMongoDB().Collection("Majors")
	indexMajorExists, errIndex := utils.CheckIndexExists(ctx, majorCollection, indexMajorName)
	if errIndex != nil {
		fmt.Println("Lỗi trong quá trình kiểm tra tồn tại index")
	}
	if !indexMajorExists {
		indexMajorModel := createIndex("major_name", "text")
		_, err := majorCollection.Indexes().CreateOne(context.TODO(), indexMajorModel)
		if err != nil {
			fmt.Println("Lỗi trong quá trình tạo index collection Major")
		}
	} else {
		fmt.Println("Index already exists:", indexMajorName)
	}
	//Subject Index
	indexSubjectName := "subject_name_text"
	subjectCollection := GetMongoDB().Collection("Subjects")
	indexSubjectExists, errIndex := utils.CheckIndexExists(ctx, subjectCollection, indexSubjectName)
	if errIndex != nil {
		fmt.Println("Lỗi trong quá trình kiểm tra tồn tại index")
	}
	if !indexSubjectExists {
		indexSubjectModel := createIndex("subject_name", "text")
		_, err := subjectCollection.Indexes().CreateOne(context.TODO(), indexSubjectModel)
		if err != nil {
			fmt.Println("Lỗi trong quá trình tạo index collection Subjects")
		}
	} else {
		fmt.Println("Index already exists:", indexSubjectName)
	}
}
