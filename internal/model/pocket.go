package model

import "gorm.io/gorm"

type Pocket struct {
	gorm.Model
	CustomerID uint    `json:"customer_id"`
	Name       string  `json:"name"`
	Balance    float64 `json:"balance"`
}
