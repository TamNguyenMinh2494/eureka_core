package models

type Enrollment struct {
	Id       string `json:"_id,omitempty"`
	CourseID string `json:"course_id" validation:"required"`
	Email    string `json:"email" validation:"required"`
	Date     int64  `json:"date"`
}
