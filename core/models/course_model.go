package models

type Course struct {
	Id               string `json:"id" validate:"required"` // User-friendly ID
	Name             string `json:"name" validate:"required"`
	AuthorName       string `json:"author_name" validate:"required"`
	AuthorEmail      string `json:"author_email" validate:"required"`
	Fee              int64  `json:"fee" validate:"required"`
	StartDate        int64  `json:"start_date"`
	AllowEnroll      bool   `json:"allow_enroll" validate:"required"`
	IsPublic         bool   `json:"is_public" validate:"required"`
	PhotoUrl         string `json:"photo_url"`
	MarketingContent string `json:"marketing_content"`
}

type CourseSection struct {
	SectionId string `json:"section_id" validate:"required"`
	CourseId  string `json:"course_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Parent    string `json:"parent"`
	PhotoUrl  string `json:"photo_url"`
	Content   string `json:"content"`
}
