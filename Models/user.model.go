package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserModel struct {
	Id             primitive.ObjectID  `bson:"_id"`
	MajorId        *primitive.ObjectID `bson:"major_id"`
	Email          string              `json:"email" bson:"email" binding:"required,email"`
	Password       string              `json:"password" bson:"password" binding:"required,min=6,max=30"`
	Role           string              `json:"role_type" bson:"role_type" binding:"required,eq=student|eq=teacher|eq=admin"`
	Phone          string              `json:"phone" bson:"phone" binding:"required,len=10""`
	Name           string              `json:"name" bson:"name" binding:"required,min=2,max=100"`
	Avatar         string              `json:"avatar" bson:"avatar" binding:"required,url"`
	Gender         bool                `json:"gender" bson:"gender"`
	Department     string              `json:"department" bson:"department"`
	DateOfBirth    time.Time           `json:"date_of_birth" bson:"date_of_birth"`
	EnrollmentDate time.Time           `json:"enrollment_date" bson:"enrollment_date"`
	HireDate       time.Time           `json:"hire_date" bson:"hire_date"`
	Address        string              `json:"address" bson:"address"`
	CreatedAt      time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at" bson:"updated_at"`
}

type AuthInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=30"`
}

type UserUpdate struct {
	MajorName      string    `json:"major_name"` // Chưa validate
	Email          string    `json:"email" binding:"required,email"`
	Password       string    `json:"password" binding:"required,min=6,max=30"`
	Role           string    `json:"role_type" binding:"required,eq=student|eq=teacher|eq=admin"`
	Phone          string    `json:"phone" binding:"required,len=10""`
	Name           string    `json:"name" binding:"required,min=2,max=100"`
	Avatar         string    `json:"avatar" binding:"required,url"`
	Gender         bool      `json:"gender"`
	Department     string    `json:"department"`
	DateOfBirth    time.Time `json:"date_of_birth"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	HireDate       time.Time `json:"hire_date"`
	Address        string    `json:"address"`
}
type CreateUserInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=30"`
	Role     string `json:"role_type" binding:"required,eq=student|eq=teacher|eq=admin"`
	Phone    string `json:"phone" binding:"required,len=10""`
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Avatar   string `json:"avatar" binding:"required,url"`
}
type AccountUpdate struct {
	Id string `json:"id"`
	UserUpdate
} //admin manager
type ChangePasswordInput struct {
	OlderPassword      string `json:"older_password" binding:"required,min=6,max=30"`
	NewPassword        string `json:"new_password,omitempty" binding:"required,min=6,max=30"`
	ConfirmNewPassword string `json:"confirm_new_password,omitempty" binding:"required,min=6,max=30"`
}
type FindEmailInput struct {
	Email string `json:"email" binding:"required,email"`
}
type ForgotPasswordInput struct {
	Email           string `json:"email" bson:"-" binding:"required,email"`
	OtpToken        string `json:"otp_token" bson:"-" binding:"required"`
	NewPassword     string `json:"new_password,omitempty" bson:"-" binding:"required,min=6,max=30"`
	ConfirmPassword string `json:"confirm_password,omitempty" bson:"-" binding:"required,min=6,max=30"`
}
