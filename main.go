package main

import (
	"context"
	"fmt"
	"gin-gonic-gom/Routes"
	"gin-gonic-gom/config"
	_ "gin-gonic-gom/docs"
	"gin-gonic-gom/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lpernett/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

// @Title  Manager Student Service API
// @Version 1.0
// @description Manager Student service API in Go using Gin Framework

// @Host localhost:4000
// @BasePath /v1
var (
	server      *gin.Engine
	ctx         context.Context
	mongoClient *mongo.Client
	validate    *validator.Validate
)

func InitializeConfig() {

	env := os.Getenv("ENV")
	if env == "production" {
		err := godotenv.Load(".env.production")
		if err != nil {
			log.Fatalf("Error loading .env.production file")
		}
	} else {
		err := godotenv.Load(".env.development")
		fmt.Println("Data --> ", os.Getenv("ENV"))
		if err != nil {
			log.Fatalf("Error loading .env.development file")
		}
	}
	if _, errFile := os.Stat("uploads/images"); os.IsNotExist(errFile) {
		os.Mkdir("uploads/images", os.ModePerm)
	}
	server = gin.Default()
	validate = validator.New()
	//serve file
	server.Static("/static", "./uploads")
}
func InitializeDatabase() {
	ctx = context.TODO()
	config.Connect(ctx)
}
func main() {
	InitializeConfig()
	InitializeDatabase()
	defer func(mongoClient *mongo.Client, ctx context.Context) {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Println("Error disconnecting MongoDB client: ---> ", err)
		}
	}(mongoClient, ctx)
	config.InitIndex()
	err := utils.InitCache()
	if err != nil {
		fmt.Println(err)
		return
	}
	//Document
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//REST API
	basepath := server.Group("/v1/api")
	//middleware cors
	Routes.Router(basepath)
	//jobs.JobRunner()
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000" // Giá trị mặc định nếu không có biến môi trường PORT
	}

	// Xác định giá trị của biến môi trường HOST, mặc định là "localhost" nếu không có
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0" // Giá trị mặc định nếu không có biến môi trường HOST
	}

	go func() {
		address := host + ":" + port
		if err := server.Run(address); err != nil {
			log.Fatalf("Error running Gin server: %v", err)
		}
	}()

	select {} // Giữ chương trình chạy 1
}
