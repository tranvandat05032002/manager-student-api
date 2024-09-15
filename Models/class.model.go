package Models

import "go.mongodb.org/mongo-driver/bson/primitive"

// // Class represents a class for a subject
type ClassModel struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ClassCode   string             `json:"class_code" bson:"class_code"`
	Teacher     string             `json:"teacher" bson:"teacher"`
	Schedule    string             `json:"schedule" bson:"schedule"`
	MaxStudents int                `json:"maxStudents" bson:"max_students"` // so luong sinh vien dang ky toi da
	Enrolled    int                `json:"enrolled" bson:"enrolled"`        // so luong sinh vien da dang ky
}
