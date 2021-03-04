package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Course struct {
	Id               string `json:"id" validate:"required"` // User-friendly ID
	Name             string `json:"name" validate:"required"`
	AuthorName       string `json:"author_name" `
	AuthorEmail      string `json:"author_email"`
	Fee              int64  `json:"fee"`
	StartDate        int64  `json:"start_date"`
	AllowEnroll      bool   `json:"allow_enroll" binding:"required"`
	IsPublic         bool   `json:"is_public" binding:"required"`
	PhotoUrl         string `json:"photo_url"`
	MarketingContent string `json:"marketing_content"`
}

type CourseSection struct {
	Id       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	CourseId string             `json:"course_id" validate:"required"`
	Name     string             `json:"name" validate:"required"`
	Parent   string             `json:"parent"`
	PhotoUrl string             `json:"photo_url"`
	Content  string             `json:"content"`
}
