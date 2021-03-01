package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Exams struct {
	Id           primitive.ObjectID `bson:"_id" json:"id"`
	Title        string             `json:"title" validate:"required"`
	SectionId    string             `json:"sectionid" validate:"required"`
	QuizSelector []QuizSelector     `json:"quiz_selector"`
	Duration     int64              `json:"duration" validate:"required"`
}

type TakenExams struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	Examinee    string             `json:"examinee"`
	CreatedDate int64              `json:"createddate"`
	Quizzes     []Quiz             `json:"quizzes"`
	Duration    int64              `json:"duration"`
}

type SubmittedExams struct {
	Id         primitive.ObjectID `bson:"_id" json:"id"`
	Examinee   string             `json:"examinee"`
	Quizzes    []Quiz             `json:"quizzes"`
	Score      float32            `json:"score"`
	TotalScore float32            `json:"totalscore"`
	SubmitDate int64              `json:"submitdate"`
}
