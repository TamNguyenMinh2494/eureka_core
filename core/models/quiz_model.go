package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Quiz struct {
	Id            primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	CourseId      string             `json:"course_id,required"`
	Question      string             `json:"question,required"`
	CorrectAnswer string             `json:"answer,required"`
	LastUpdate    int64              `json:"last_update"`
}
