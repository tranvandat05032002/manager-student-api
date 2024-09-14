package main

import (
	"context"
	"fmt"
	"gin-gonic-gom/Controllers"
	"gin-gonic-gom/Services/major"
	"gin-gonic-gom/Services/media"
	statiscal "gin-gonic-gom/Services/statistical"
	"gin-gonic-gom/Services/subject"
	"gin-gonic-gom/Services/term"
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
	"github.com/robfig/cron/v3"
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
	ms           media.MediaService
	mc           Controllers.MediaController
	mjs          major.MajorService
	mjc          Controllers.MajorController
	subs         subject.SubjectService
	subc         Controllers.SubjectController
	terms        term.TermService
	termc        Controllers.TermController
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
	mongoClient  *mongo.Client
	validate     *validator.Validate
)

func createIndex(collection *mongo.Collection, indexName string, indexType interface{}) mongo.IndexModel {
	if indexType == "text" {
		indexModelText := mongo.IndexModel{Keys: bson.D{{indexName, indexType}}, Options: options.Index().SetDefaultLanguage("none")}
		return indexModelText
	}
	indexModelNotText := mongo.IndexModel{Keys: bson.D{{indexName, indexType}}}
	return indexModelNotText
}
func InitializeConfig() {
	env := os.Getenv("ENV")
	if env != "prod" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	if _, errFile := os.Stat("uploads/images"); os.IsNotExist(errFile) {
		os.Mkdir("uploads/images", os.ModePerm)
	}
	server = gin.Default()
	//config TrustedProxies and IPV6
	server.SetTrustedProxies([]string{
		os.Getenv("HOST"),
		"192.168.1.10",                          // Địa chỉ IPv4 proxy tin cậy
		"::1",                                   // Địa chỉ IPv6 localhost
		"2405:4802:6563:140:10b3:beb9:c910:16e", // Địa chỉ IPv6 proxy tin cậy
	})
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
	index_name_user := "name_text"
	indexUserExists, errIndex := utils.CheckIndexExists(ctx, userco, index_name_user)
	if errIndex != nil {
		fmt.Println("Lỗi trong quá trình kiểm tra tồn tại index")
	}
	if !indexUserExists {
		indexUserModel := createIndex(userco, "name", "text")
		userco.Indexes().CreateOne(context.TODO(), indexUserModel)
	} else {
		fmt.Println("Index already exists:", index_name_user)
	}
	// Token collection
	tokenco = mongoCon.Collection("Tokens")
	// OTP collection
	otpco = mongoCon.Collection("OTPs")
	//Media collection
	mediaco = mongoCon.Collection("Medias")
	// Major collection
	majorco = mongoCon.Collection("Majors")
	indexMajorName := "major_name_text"
	indexMajorExists, errIndex := utils.CheckIndexExists(ctx, majorco, indexMajorName)
	if errIndex != nil {
		fmt.Println("Lỗi trong quá trình kiểm tra tồn tại index")
	}
	if !indexMajorExists {
		indexMajorModel := createIndex(majorco, "major_name", "text")
		majorco.Indexes().CreateOne(context.TODO(), indexMajorModel)
	} else {
		fmt.Println("Index already exists:", indexMajorName)
	}
	// Subject collection
	subjectco = mongoCon.Collection("Subjects")
	indexSubjectName := "subject_name_text"
	indexSubjectExists, errIndex := utils.CheckIndexExists(ctx, subjectco, indexSubjectName)
	if errIndex != nil {
		fmt.Println("Lỗi trong quá trình kiểm tra tồn tại index")
	}
	if !indexSubjectExists {
		indexSubjectModel := createIndex(subjectco, "subject_name", "text")
		subjectco.Indexes().CreateOne(context.TODO(), indexSubjectModel)
	} else {
		fmt.Println("Index already exists:", indexSubjectName)
	}
	// Term collection
	termco = mongoCon.Collection("Terms")

	us = user.NewUserService(userco, majorco, tokenco, otpco, ctx)
	uc = Controllers.New(us)
	ms = media.NewMediaService(mediaco, userco, ctx)
	mc = Controllers.NewMedia(ms)

	mjs = major.NewMajorService(majorco, ctx)
	mjc = Controllers.NewMajor(mjs)

	subs = subject.NewMajorService(subjectco, termco, ctx)
	subc = Controllers.NewSubject(subs)

	terms = term.NewTermService(termco, ctx)
	termc = Controllers.NewTerm(terms)

	statisticals = statiscal.NewStatisticalService(termco, ctx)
	statisticalc = Controllers.NewStatistical(statisticals)
}
func main() {
	InitializeConfig()
	InitializeDatabase()
	defer func(mongoClient *mongo.Client, ctx context.Context) {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Println("Error disconnecting MongoDB client: ---> ", err)
			return
		}
	}(mongoClient, ctx)
	if err := utils.InitCache(); err != nil {
		fmt.Println(err)
		return
	}
	//Document
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//REST API
	basepath := server.Group("/v1")
	uc.RegisterUserRoutes(basepath)
	mc.RegisterMediaRoutes(basepath)
	mjc.RegisterMajorRoutes(basepath)
	subc.RegisterSubjectRoutes(basepath)
	termc.RegisterTermRoutes(basepath)
	statisticalc.RegisterStatisticalRoutes(basepath)
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	c := cron.New(cron.WithLocation(loc))
	// 4 giờ sáng mỗi ngày thì cron job sẽ hoạt động để xóa token/otp hết hạn
	if _, err := c.AddFunc("0 4 * * *", us.DeleteTokenExp); err != nil {
		log.Fatalf("Error adding cron job delete Token: %v", err)
	}
	if _, errOTP := c.AddFunc("0 4 * * *", us.DeleteOTPExp); errOTP != nil {
		log.Fatalf("Error adding cron job delete OTP: %v", errOTP)
	}
	if _, errDelUser := c.AddFunc("@every 1m", us.CheckAndDeleteUsers); errDelUser != nil {
		log.Fatalf("Error adding cron job delete User: %v", errDelUser)
	}
	c.Start()
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
