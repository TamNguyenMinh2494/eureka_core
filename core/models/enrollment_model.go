package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Enrollment struct {
	Id       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	CourseID string             `json:"course_id" validation:"required"`
	Email    string             `json:"email" validation:"required"`
	Date     int64              `json:"date"`
}
