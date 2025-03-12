package model

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	Name         string        `json:"name"`
	Email        string        `json:"email"`
	BankAccounts []BankAccount `json:"bank_accounts"`
	Pockets      []Pocket      `json:"pockets"`
	TermDeposits []TermDeposit `json:"term_deposits"`
}

type CustomerResponse struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name"`
	Email        string        `json:"email"`
	BankAccounts []BankAccount `json:"bank_accounts,omitempty"`
	Pockets      []Pocket      `json:"pockets,omitempty"`
	TermDeposits []TermDeposit `json:"term_deposits,omitempty"`
}
