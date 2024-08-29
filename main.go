package main

import (
	"context"
	"gin-gonic-gom/Controllers"
	"gin-gonic-gom/Services/major"
	"gin-gonic-gom/Services/media"
	"gin-gonic-gom/Services/subject"
	"gin-gonic-gom/Services/term"
	"gin-gonic-gom/Services/user"
	"gin-gonic-gom/config"
	"github.com/gin-gonic/gin"
	"github.com/lpernett/godotenv"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"time"
)

var (
	server      *gin.Engine
	us          user.UserService
	uc          Controllers.UserController
	ms          media.MediaService
	mc          Controllers.MediaController
	mjs         major.MajorService
	mjc         Controllers.MajorController
	subs        subject.SubjectService
	subc        Controllers.SubjectController
	terms       term.TermService
	termc       Controllers.TermController
	ctx         context.Context
	userco      *mongo.Collection
	tokenco     *mongo.Collection
	otpco       *mongo.Collection
	mediaco     *mongo.Collection
	majorco     *mongo.Collection
	subjectco   *mongo.Collection
	termco      *mongo.Collection
	mongoClient *mongo.Client
)

func InitializeConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if _, errFile := os.Stat("uploads/images"); os.IsNotExist(errFile) {
		os.Mkdir("uploads/images", os.ModePerm)
	}
	server = gin.Default()
	//serve file
	server.Static("/static", "./uploads")
}
func InitializeDatabase() {
	ctx = context.TODO()
	mongoCon, _ := config.Connect(ctx)
	userco = mongoCon.Collection("Users")
	tokenco = mongoCon.Collection("Tokens")
	otpco = mongoCon.Collection("OTPs")
	mediaco = mongoCon.Collection("Medias")
	majorco = mongoCon.Collection("Majors")
	subjectco = mongoCon.Collection("Subjects")
	termco = mongoCon.Collection("Terms")

	us = user.NewUserService(userco, majorco, tokenco, otpco, ctx)
	uc = Controllers.New(us)

	ms = media.NewMediaService(mediaco, ctx)
	mc = Controllers.NewMedia(ms)

	mjs = major.NewMajorService(majorco, ctx)
	mjc = Controllers.NewMajor(mjs)

	subs = subject.NewMajorService(subjectco, termco, ctx)
	subc = Controllers.NewSubject(subs)

	terms = term.NewTermService(termco, ctx)
	termc = Controllers.NewTerm(terms)
}
func main() {
	InitializeConfig()
	InitializeDatabase()
	defer func(mongoClient *mongo.Client, ctx context.Context) {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Println("Error disconnecting MongoDB client: ---> ", err)
		}
	}(mongoClient, ctx)
	basepath := server.Group("/v1")
	uc.RegisterUserRoutes(basepath)
	mc.RegisterMediaRoutes(basepath)
	mjc.RegisterMajorRoutes(basepath)
	subc.RegisterSubjectRoutes(basepath)
	termc.RegisterTermRoutes(basepath)
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	c := cron.New(cron.WithLocation(loc))
	// 4 giờ sáng mỗi ngày thì cron job sẽ hoạt động để xóa token/otp hết hạn
	_, err := c.AddFunc("0 4 * * *", us.DeleteTokenExp)
	_, errOTP := c.AddFunc("0 4 * * *", us.DeleteOTPExp)
	if err != nil {
		log.Fatalf("Error adding cron job: %v", err)
	}
	if errOTP != nil {
		log.Fatalf("Error adding cron job: %v", err)
	}

	c.Start()
	go func() {
		if err := server.Run(os.Getenv("HOST") + ":" + os.Getenv("PORT")); err != nil {
			log.Fatalf("Error running Gin server: %v", err)
		}
	}()

	select {} // Giữ chương trình chạy 1
}
