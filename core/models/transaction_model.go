package models

type Transaction struct {
	Email     string `json:"email" validate:"required,email"`
	SKU       string `json:"sku" validate:"required"`
	Quantity  int64  `json:"quantity" validate:"required"`
	Amount    int64  `json:"amount" validate:"required"`
	timestamp string `json:"timestamp" validate:"required"`
}
