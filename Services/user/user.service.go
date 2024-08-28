package user

import (
	"gin-gonic-gom/Models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	CreateUser(*Models.CreateUserInput) error
	CheckExistEmail(string) (bool, error)
	DeleteTokenExp()
	DeleteOTPExp()
	LoginUser(*Models.AuthInput, *gin.Context)
	GetMe(string, string) (*Models.UserModel, error)
	UpdateMe(string, *Models.UserUpdate) (*Models.UserModel, error)
	GetAccount(string) (*Models.UserModel, error)
	GetAll(int, int) ([]*Models.UserModel, int, error)
	UpdateUser(*Models.AccountUpdate, primitive.ObjectID) (*Models.UserModel, error)
	DeleteUser(string) (int, error)
	ChangePassword(string, *Models.ChangePasswordInput) error
	FindEmail(string) (*Models.UserModel, error)
	ForgotPasswordByOTP(*Models.ForgotPasswordInput) (bool, error)
	SaveOTPForUser(primitive.ObjectID, string) error
	VerifyOTP(string, string) (bool, error)
	ResendOTP(primitive.ObjectID, string) (bool, error)

	// Student
	//GetStudentDetails(primitive.ObjectID, primitive.ObjectID) (*Models.StudentDetail, error)
}
