package repository

import (
	"github.com/danisasmita/customer-search/internal/model"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	FindByName(name, email, accountNumber string) ([]model.Customer, error)
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) FindByName(name, email, accountNumber string) ([]model.Customer, error) {
	var customers []model.Customer

	query := r.db.Preload("BankAccounts", func(db *gorm.DB) *gorm.DB {
		return db.Select("customer_id, account_number")
	}).
		Preload("Pockets", func(db *gorm.DB) *gorm.DB {
			return db.Select("customer_id, name, balance")
		}).
		Preload("TermDeposits", func(db *gorm.DB) *gorm.DB {
			return db.Select("customer_id, amount, duration")
		})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}
	if accountNumber != "" {
		query = query.Joins("JOIN bank_accounts ON bank_accounts.customer_id = customers.id").
			Where("bank_accounts.account_number = ?", accountNumber)
	}

	err := query.Find(&customers).Error
	return customers, err
}
