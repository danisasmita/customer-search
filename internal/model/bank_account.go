package model

import "gorm.io/gorm"

type BankAccount struct {
	gorm.Model
	CustomerID    uint    `json:"customer_id"`
	AccountNumber string  `json:"account_number"`
	Balance       float64 `json:"balance"`
}
