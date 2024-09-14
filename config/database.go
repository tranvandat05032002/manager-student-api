package config

import (
	"context"
	"fmt"
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
	MongoURL := os.Getenv("MONGO_URL")
	if MongoURL == "" {
		MongoURL = "mongodb://admin:tranlybuu@mongo:1236/?authSource=admin"
	}
	mongoconn := options.Client().ApplyURI(MongoURL)
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
