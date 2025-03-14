package database

import (
	"testing"

	"github.com/danisasmita/customer-search/internal/config"
	"github.com/danisasmita/customer-search/internal/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	// Menggunakan SQLite dalam memori untuk testing
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestConnectDB(t *testing.T) {
	cfg := &config.Config{
		DBDriver: "sqlite",
		DBSource: "file::memory:?cache=shared",
	}

	db, err := ConnectDB(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func TestAutoMigrate(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	err = AutoMigrate(db)
	assert.NoError(t, err)

	// Cek apakah tabel-tabel sudah dibuat
	assert.True(t, db.Migrator().HasTable(&model.Customer{}))
	assert.True(t, db.Migrator().HasTable(&model.BankAccount{}))
	assert.True(t, db.Migrator().HasTable(&model.Pocket{}))
	assert.True(t, db.Migrator().HasTable(&model.TermDeposit{}))
	assert.True(t, db.Migrator().HasTable(&model.User{}))
}

func TestSeedData(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	err = AutoMigrate(db)
	assert.NoError(t, err)

	err = SeedData(db)
	assert.NoError(t, err)

	// Cek apakah data sudah di-seed
	var count int64
	db.Model(&model.Customer{}).Count(&count)
	assert.Equal(t, int64(11), count) // Sesuaikan dengan jumlah data yang di-seed
}

func TestSeedDataNoDuplicate(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	err = AutoMigrate(db)
	assert.NoError(t, err)

	// Seed data pertama kali
	err = SeedData(db)
	assert.NoError(t, err)

	// Seed data kedua kali
	err = SeedData(db)
	assert.NoError(t, err)

	// Cek apakah tidak ada data duplikat
	var count int64
	db.Model(&model.Customer{}).Count(&count)
	assert.Equal(t, int64(11), count) // Sesuaikan dengan jumlah data yang di-seed
}
