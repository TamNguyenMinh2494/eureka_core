package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type QuizSelector struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	SectionId string             `json:"sectionId" validate:"required"`
	List      []string           `json:"list" validate:"required"`
	Shuffle   bool               `json:"shuffle" binding:"required"`
	Filter    string             `json:"filter"`
}
