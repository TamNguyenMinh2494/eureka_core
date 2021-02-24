package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Exams struct {
	Id           primitive.ObjectID `bson:"_id" json:"id"`
	Title        string             `json:"title" validate:"required"`
	SectionId    string             `json:"sectionid" validate:"required"`
	QuizSelector []QuizSelector     `json:"quiz_selector"`
}

type TakenExams struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Examinee    string             `json:"examinee"`
	CreatedDate int64              `json:"createddate"`
	Quizzes     []Quiz             `json:"quizzes"`
}
