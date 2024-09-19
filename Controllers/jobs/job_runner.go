package jobs

import (
	"gin-gonic-gom/Services/user"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

func JobRunner(us user.UserService) {
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	c := cron.New(cron.WithLocation(loc))
	// 4 giờ sáng mỗi ngày thì cron job sẽ hoạt động để xóa token/otp hết hạn
	if _, err := c.AddFunc("0 4 * * *", us.DeleteTokenExp); err != nil {
		log.Fatalf("Error adding cron job delete Token: %v", err)
	}
	if _, errOTP := c.AddFunc("0 4 * * *", us.DeleteOTPExp); errOTP != nil {
		log.Fatalf("Error adding cron job delete OTP: %v", errOTP)
	}
	//if _, errDelUser := c.AddFunc("@every 1m", us.CheckAndDeleteUsers); errDelUser != nil {
	//	log.Fatalf("Error adding cron job delete User: %v", errDelUser)
	//}
	c.Start()
}
