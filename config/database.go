package config

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

var (
	err         error
	mongoclient *mongo.Client
	mongoDB     *mongo.Database
)

func Connect(ctx context.Context) (*mongo.Database, error) {

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
