package models

type Course struct {
	Id               string `json:"_id,omitempty"`
	Name             string `json:"name" validate:"required"`
	AuthorName       string `json:"author_name" validate:"required"`
	AuthorEmail      string `json:"author_email" validate:"required"`
	Fee              int64  `json:"fee"`
	StartDate        int64  `json:"start_date"`
	AllowEnroll      bool   `json:"allow_enroll" validate:"required"`
	IsPublic         bool   `json:"is_public" validate:"required"`
	MarketingContent string `json:"marketing_content"`
}
