package main

import (
	"context"
	"fmt"
	"gin-gonic-gom/Controllers"
	"gin-gonic-gom/Controllers/jobs"
	"gin-gonic-gom/Routes"
	statiscal "gin-gonic-gom/Services/statistical"
	"gin-gonic-gom/Services/user"
	"gin-gonic-gom/config"
	_ "gin-gonic-gom/docs"
	"gin-gonic-gom/utils"
	"github.com/lpernett/godotenv"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Title  Manager Student Service API
// @Version 1.0
// @description Manager Student service API in Go using Gin Framework

// @Host localhost:4000
// @BasePath /v1
var (
	server       *gin.Engine
	us           user.UserService
	uc           Controllers.UserController
	statisticals statiscal.StatisticalService
	statisticalc Controllers.StatisticalController
	ctx          context.Context
	userco       *mongo.Collection
	tokenco      *mongo.Collection
	otpco        *mongo.Collection
	mediaco      *mongo.Collection
	majorco      *mongo.Collection
	subjectco    *mongo.Collection
	termco       *mongo.Collection
	scheduleco   *mongo.Collection
	mongoClient  *mongo.Client
	validate     *validator.Validate
)

func createIndex(indexName string, indexType interface{}) mongo.IndexModel {
	if indexType == "text" {
		indexModelText := mongo.IndexModel{Keys: bson.D{{indexName, indexType}}, Options: options.Index().SetDefaultLanguage("none")}
		return indexModelText
	}
	indexModelNotText := mongo.IndexModel{Keys: bson.D{{indexName, indexType}}}
	return indexModelNotText
}
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
	//config cors
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, //client
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	validate = validator.New()
	//serve file
	server.Static("/static", "./uploads")
}
func InitializeDatabase() {
	ctx = context.TODO()
	mongoCon, _ := config.Connect(ctx)
	// Users collection
	userco = mongoCon.Collection("Users")
	indexNameUser := "name_text"
	indexUserExists, errIndex := utils.CheckIndexExists(ctx, userco, indexNameUser)
	if errIndex != nil {
		fmt.Println("Lỗi trong quá trình kiểm tra tồn tại index")
	}
	if !indexUserExists {
		indexUserModel := createIndex("name", "text")
		_, err := userco.Indexes().CreateOne(context.TODO(), indexUserModel)
		if err != nil {
			fmt.Println("Lỗi trong quá trình tạo index collection Users")
		}
	} else {
		fmt.Println("Index already exists:", indexNameUser)
	}
	// Token collection
	tokenco = mongoCon.Collection("Tokens")
	// OTP collection
	otpco = mongoCon.Collection("OTPs")
	//Schedule collection
	scheduleco = mongoCon.Collection("Schedules")

	us = user.NewUserService(userco, majorco, tokenco, otpco, ctx)
	uc = Controllers.New(us)

	statisticals = statiscal.NewStatisticalService(termco, ctx)
	statisticalc = Controllers.NewStatistical(statisticals)
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
	//err := utils.InitCache()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//Document
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//REST API
	basepath := server.Group("/v1/api")
	Routes.Router(basepath)
	//uc.RegisterAuthRoutes(basepath)
	//statisticalc.RegisterStatisticalRoutes(basepath)
	//schedulec.RegisterScheduleRoutes(basepath)
	jobs.JobRunner(us)
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
