package models

type StudentInvitation struct {
	CourseId string `json:"course_id" validate:"required"`
	Email    string `json:"email" validate:"required"`
}
