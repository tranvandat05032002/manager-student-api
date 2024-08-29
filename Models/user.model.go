package Models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserModel struct {
	Id             primitive.ObjectID  `bson:"_id"`
	MajorId        *primitive.ObjectID `bson:"major_id"`
	Name           string              `json:"name" bson:"name" binding:"required,min=2,max=100"`
	Email          string              `json:"email" bson:"email" binding:"required,email"`
	Password       string              `json:"password,omitempty" bson:"password,omitempty" binding:"required,min=6,max=100"`
	Avatar         string              `json:"avatar" bson:"avatar" binding:"required,url"`
	Phone          string              `json:"phone" bson:"phone"`
	Gender         *bool               `bson:"gender" json:"gender"`
	Department     string              `json:"department" bson:"department"`
	DateOfBirth    time.Time           `json:"date_of_birth" bson:"date_of_birth"`
	EnrollmentDate time.Time           `json:"enrollment_date" bson:"enrollment_date"`
	HireDate       time.Time           `json:"hire_date" bson:"hire_date"`
	Address        string              `json:"address" bson:"address"`
	Role           string              `json:"role_type,omitempty" bson:"role_type,omitempty" binding:"required,eq=student|eq=teacher|eq=admin"`
	CreatedAt      time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at" bson:"updated_at"`
}

type AuthInput struct {
	Email    string `json:"email" bson:"email",binding:"required,email"`
	Password string `json:"password" bson:"password" binding:"required,min=6,max=100"`
}

type UserUpdate struct {
	MajorName      string    `json:"major_name"`
	Email          string    `json:"email,omitempty"`
	Name           string    `json:"name,omitempty"`
	Avatar         string    `json:"avatar,omitempty"`
	Phone          string    `json:"phone,omitempty"`
	Gender         bool      `json:"gender"`
	Department     string    `json:"department"`
	DateOfBirth    time.Time `json:"date_of_birth,omitempty"`
	EnrollmentDate time.Time `json:"enrollment_date,omitempty"`
	HireDate       time.Time `json:"hire_date,omitempty"`
	Address        string    `json:"address,omitempty"`
}
type CreateUserInput struct {
	Email    string `json:"email" `
	Password string `json:"password"`
	Role     string `json:"role_type"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
}
type AccountUpdate struct {
	Id string `json:"id"`
	UserUpdate
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	Role     string `json:"role_type"`
} //admin manager
type ChangePasswordInput struct {
	OlderPassword      string `json:"older_password"`
	NewPassword        string `json:"new_password,omitempty" bson:"new_password,omitempty" binding:"required,min=6,max=100"`
	ConfirmNewPassword string `json:"confirm_new_password,omitempty" bson:"confirm_new_password,omitempty" binding:"required,min=6,max=100"`
}
type FindEmailInput struct {
	Email string `json:"email" binding:"required,email"`
}
type ForgotPasswordInput struct {
	Email           string `json:"email" bson:"email" binding:"required,email"`
	OtpToken        string `json:"otp_token" bson:"otp_token" binding:"required"`
	NewPassword     string `json:"new_password,omitempty" bson:"new_password,omitempty" binding:"required,min=6,max=100"`
	ConfirmPassword string `json:"confirm_password,omitempty" bson:"confirm_password,omitempty" binding:"required,min=6,max=100"`
}
