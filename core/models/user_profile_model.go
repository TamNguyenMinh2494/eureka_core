package models

type UserProfile struct {
	DisplayName string `json:"display_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Gender      string `json:"gender" validate:"required"`
	DOB         uint64 `json:"dob" validate:"required"`
	Role        string `json:"role"`
}

type UpdatedUserProfile struct {
	DisplayName string `json:"display_name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Role        string `json:"role,omitempty"`
}

type UserAuth struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
