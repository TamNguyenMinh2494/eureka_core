package models

import (
	"main/core/question"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Quiz struct {
	Id            primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	CourseId      string             `json:"course_id" validate:"required"`
	Question      string             `json:"question" validate:"required"`
	CorrectAnswer string             `json:"answer" validate:"required"`
	Type          string             `json:"type" validate:"required"`
	LastUpdate    int64              `json:"last_update"`
}

type StructuredQuiz struct {
	Raw      Quiz
	Question question.Question
}
