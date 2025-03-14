package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	customerName      = "John Doe"
	customerEmail     = "john@example.com"
	customerEmailJane = "jane@example.com"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
	}

	return gormDB, mock
}

func TestCustomerRepositoryFindByName(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	repo := NewCustomerRepository(gormDB)

	t.Run("success - search by name", func(t *testing.T) {
		// Mock data untuk customers
		rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, customerName, customerEmail, time.Now(), time.Now(), nil)

		// Mock query untuk customers
		mock.ExpectQuery(`SELECT \* FROM "customers" WHERE name LIKE \$1 AND "customers"."deleted_at" IS NULL`).
			WithArgs("%John%").
			WillReturnRows(rows)

		// Mock query untuk bank_accounts
		mock.ExpectQuery(`SELECT customer_id, account_number FROM "bank_accounts" WHERE "bank_accounts"."customer_id" = \$1 AND "bank_accounts"."deleted_at" IS NULL`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"customer_id", "account_number"}))

		// Mock query untuk pockets
		mock.ExpectQuery(`SELECT customer_id, name, balance FROM "pockets" WHERE "pockets"."customer_id" = \$1 AND "pockets"."deleted_at" IS NULL`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"customer_id", "name", "balance"}))

		// Mock query untuk term_deposits
		mock.ExpectQuery(`SELECT customer_id, amount, duration FROM "term_deposits" WHERE "term_deposits"."customer_id" = \$1 AND "term_deposits"."deleted_at" IS NULL`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"customer_id", "amount", "duration"}))

		customers, err := repo.FindByName("John", "", "")

		assert.NoError(t, err)
		assert.Equal(t, 1, len(customers))
		assert.Equal(t, customerName, customers[0].Name)
		assert.Equal(t, customerEmail, customers[0].Email)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("success - search by email", func(t *testing.T) {
		// Mock data untuk customers
		rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "Jane Doe", customerEmailJane, time.Now(), time.Now(), nil)

		// Mock query untuk customers
		mock.ExpectQuery(`SELECT \* FROM "customers" WHERE email LIKE \$1 AND "customers"."deleted_at" IS NULL`).
			WithArgs("%jane@example.com%").
			WillReturnRows(rows)

		// Mock query untuk bank_accounts
		mock.ExpectQuery(`SELECT customer_id, account_number FROM "bank_accounts" WHERE "bank_accounts"."customer_id" = \$1 AND "bank_accounts"."deleted_at" IS NULL`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"customer_id", "account_number"}))

		// Mock query untuk pockets
		mock.ExpectQuery(`SELECT customer_id, name, balance FROM "pockets" WHERE "pockets"."customer_id" = \$1 AND "pockets"."deleted_at" IS NULL`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"customer_id", "name", "balance"}))

		// Mock query untuk term_deposits
		mock.ExpectQuery(`SELECT customer_id, amount, duration FROM "term_deposits" WHERE "term_deposits"."customer_id" = \$1 AND "term_deposits"."deleted_at" IS NULL`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"customer_id", "amount", "duration"}))

		customers, err := repo.FindByName("", customerEmailJane, "")

		assert.NoError(t, err)
		assert.Equal(t, 1, len(customers))
		assert.Equal(t, "Jane Doe", customers[0].Name)
		assert.Equal(t, "jane@example.com", customers[0].Email)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
	// Test case: Success - search by account number
	t.Run("success - search by account number", func(t *testing.T) {
		// Mock data untuk customers
		rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "name", "email"}).
			AddRow(1, time.Now(), time.Now(), nil, customerName, customerEmail)

		// Mock query untuk customers dengan JOIN bank_accounts
		mock.ExpectQuery(`SELECT "customers"."id","customers"."created_at","customers"."updated_at","customers"."deleted_at","customers"."name","customers"."email" FROM "customers" JOIN bank_accounts ON bank_accounts.customer_id = customers.id WHERE bank_accounts.account_number = \$1 AND "customers"."deleted_at" IS NULL`).
			WithArgs("123456").
			WillReturnRows(rows)

		// Mock query untuk bank_accounts
		mock.ExpectQuery(`SELECT customer_id, account_number FROM "bank_accounts" WHERE "bank_accounts"."customer_id" = \$1 AND "bank_accounts"."deleted_at" IS NULL`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"customer_id", "account_number"}).
				AddRow(1, "123456"))

		// Mock query untuk pockets
		mock.ExpectQuery(`SELECT customer_id, name, balance FROM "pockets" WHERE "pockets"."customer_id" = \$1 AND "pockets"."deleted_at" IS NULL`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"customer_id", "name", "balance"}))

		// Mock query untuk term_deposits
		mock.ExpectQuery(`SELECT customer_id, amount, duration FROM "term_deposits" WHERE "term_deposits"."customer_id" = \$1 AND "term_deposits"."deleted_at" IS NULL`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"customer_id", "amount", "duration"}))

		// Panggil fungsi repository
		customers, err := repo.FindByName("", "", "123456")

		// Assertions
		assert.NoError(t, err)                                  // Pastikan tidak ada error
		assert.Equal(t, 1, len(customers))                      // Pastikan jumlah data customer adalah 1
		assert.Equal(t, "John Doe", customers[0].Name)          // Pastikan nama customer sesuai
		assert.Equal(t, "john@example.com", customers[0].Email) // Pastikan email customer sesuai

		// Pastikan semua ekspektasi mock terpenuhi
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("error - no customers found", func(t *testing.T) {
		// Mock query untuk customers (tidak ada data)
		mock.ExpectQuery(`SELECT \* FROM "customers" WHERE name LIKE \$1 AND "customers"."deleted_at" IS NULL`).
			WithArgs("%Unknown%").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "created_at", "updated_at", "deleted_at"}))

		customers, err := repo.FindByName("Unknown", "", "")

		assert.NoError(t, err)
		assert.Equal(t, 0, len(customers))

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error - database error", func(t *testing.T) {
		// Mock query untuk customers (error database)
		mock.ExpectQuery(`SELECT \* FROM "customers" WHERE name LIKE \$1 AND "customers"."deleted_at" IS NULL`).
			WithArgs("%John%").
			WillReturnError(gorm.ErrInvalidDB) // Simulasikan error database

		// Panggil fungsi repository
		customers, err := repo.FindByName("John", "", "")

		// Assertions
		assert.Error(t, err)     // Pastikan error tidak nil
		assert.Nil(t, customers) // Pastikan data customer adalah nil

		// Pastikan semua ekspektasi mock terpenuhi
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
