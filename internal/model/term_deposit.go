package model

import "gorm.io/gorm"

type TermDeposit struct {
	gorm.Model
	CustomerID uint    `json:"customer_id"`
	Amount     float64 `json:"amount"`
	Duration   int     `json:"duration"`
}
