package database

import (
	"fmt"

	"github.com/danisasmita/customer-search/internal/config"
	"github.com/danisasmita/customer-search/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database interface {
	ConnectDB(cfg interface{}) (*gorm.DB, error)
	AutoMigrate(db *gorm.DB) error
	SeedData(db *gorm.DB) error
}

func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	if cfg.DBDriver == "sqlite" {
		db, err = gorm.Open(sqlite.Open(cfg.DBSource), &gorm.Config{})
	} else {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	return db, err
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Customer{},
		&model.BankAccount{},
		&model.Pocket{},
		&model.TermDeposit{},
		&model.User{},
	)
}

func SeedData(db *gorm.DB) error {
	var count int64
	db.Model(&model.Customer{}).Count(&count)
	if count > 0 {
		return nil
	}

	customers := []model.Customer{
		{
			Name:  "John Doe",
			Email: "john@example.com",
			BankAccounts: []model.BankAccount{
				{AccountNumber: "1234567890", Balance: 1000.0},
			},
			Pockets: []model.Pocket{
				{Name: "Savings", Balance: 500.0},
			},
			TermDeposits: []model.TermDeposit{
				{Amount: 2000.0, Duration: 12},
			},
		},
		{
			Name:  "Jane Smith",
			Email: "jane@example.com",
			BankAccounts: []model.BankAccount{
				{AccountNumber: "2345678901", Balance: 2500.0},
			},
			Pockets: []model.Pocket{
				{Name: "Emergency", Balance: 800.0},
			},
			TermDeposits: []model.TermDeposit{
				{Amount: 5000.0, Duration: 24},
			},
		},
		{
			Name:  "Robert Johnson",
			Email: "robert@example.com",
			BankAccounts: []model.BankAccount{
				{AccountNumber: "3456789012", Balance: 3200.0},
			},
			Pockets: []model.Pocket{
				{Name: "Vacation", Balance: 1200.0},
			},
			TermDeposits: []model.TermDeposit{
				{Amount: 10000.0, Duration: 36},
			},
		},
		{
			Name:  "Emily Davis",
			Email: "emily@example.com",
			BankAccounts: []model.BankAccount{
				{AccountNumber: "4567890123", Balance: 4300.0},
			},
			Pockets: []model.Pocket{
				{Name: "Education", Balance: 2000.0},
			},
			TermDeposits: []model.TermDeposit{
				{Amount: 3500.0, Duration: 6},
			},
		},
		{
			Name:  "Michael Wilson",
			Email: "michael@example.com",
			BankAccounts: []model.BankAccount{
				{AccountNumber: "5678901234", Balance: 7500.0},
			},
			Pockets: []model.Pocket{
				{Name: "Car", Balance: 3000.0},
			},
			TermDeposits: []model.TermDeposit{
				{Amount: 15000.0, Duration: 48},
			},
		},
		{
			Name:  "Sarah Brown",
			Email: "sarah@example.com",
			BankAccounts: []model.BankAccount{
				{AccountNumber: "6789012345", Balance: 1800.0},
			},
			Pockets: []model.Pocket{
				{Name: "House", Balance: 5000.0},
			},
			TermDeposits: []model.TermDeposit{
				{Amount: 8000.0, Duration: 18},
			},
		},
		{
			Name:  "David Lee",
			Email: "david@example.com",
			BankAccounts: []model.BankAccount{
				{AccountNumber: "7890123456", Balance: 9200.0},
			},
			Pockets: []model.Pocket{
				{Name: "Gadgets", Balance: 700.0},
			},
			TermDeposits: []model.TermDeposit{
				{Amount: 6000.0, Duration: 9},
			},
		},
		{
			Name:  "Jennifer Taylor",
			Email: "jennifer@example.com",
			BankAccounts: []model.BankAccount{
				{AccountNumber: "8901234567", Balance: 4100.0},
			},
			Pockets: []model.Pocket{
				{Name: "Travel", Balance: 1500.0},
			},
			TermDeposits: []model.TermDeposit{
				{Amount: 12000.0, Duration: 30},
			},
		},
		{
			Name:  "Kevin Martinez",
			Email: "kevin@example.com",
			BankAccounts: []model.BankAccount{
				{AccountNumber: "9012345678", Balance: 6700.0},
			},
			Pockets: []model.Pocket{
				{Name: "Business", Balance: 4500.0},
			},
			TermDeposits: []model.TermDeposit{
				{Amount: 25000.0, Duration: 60},
			},
		},
		{
			Name:  "Lisa Anderson",
			Email: "lisa@example.com",
			BankAccounts: []model.BankAccount{
				{AccountNumber: "0123456789", Balance: 3400.0},
			},
			Pockets: []model.Pocket{
				{Name: "Wedding", Balance: 7000.0},
			},
			TermDeposits: []model.TermDeposit{
				{Amount: 9500.0, Duration: 15},
			},
		},
		{
			Name:  "Thomas Wright",
			Email: "thomas@example.com",
			BankAccounts: []model.BankAccount{
				{AccountNumber: "1122334455", Balance: 5600.0},
			},
			Pockets: []model.Pocket{
				{Name: "Retirement", Balance: 10000.0},
			},
			TermDeposits: []model.TermDeposit{
				{Amount: 30000.0, Duration: 72},
			},
		},
	}
	return db.Create(&customers).Error
}
