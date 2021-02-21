package models

type Account struct {
	Email   string `json:"email" validate:"required"`
	Balance int64  `json:"balance" validate:"required"`
}
