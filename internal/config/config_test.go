package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Set TEST_ENV agar LoadConfig pakai .env.test
	os.Setenv("TEST_ENV", "true")

	// Cek apakah file .env.test ada
	if _, err := os.Stat(".env.test"); os.IsNotExist(err) {
		// Buat file .env.test dengan konten default jika tidak ada
		file, err := os.Create(".env.test")
		if err != nil {
			t.Fatal("Gagal membuat file .env.test:", err)
		}
		defer file.Close()

		// Tulis konten default ke dalam file
		file.WriteString("DB_HOST=localhost\n")
		file.WriteString("DB_PORT=5432\n")
		file.WriteString("DB_USER=testuser\n")
		file.WriteString("DB_PASSWORD=testpassword\n")
		file.WriteString("DB_NAME=testdb\n")
		file.WriteString("DB_SSL_MODE=disable\n")
		file.WriteString("JWT_SECRET=supersecret\n")
		file.WriteString("JWT_EXPIRATION=3600\n")
	}

	// Panggil LoadConfig()
	config, err := LoadConfig()

	// Cek tidak ada error
	assert.NoError(t, err, "LoadConfig seharusnya tidak error")

	// Validasi nilai dari .env.test
	assert.Equal(t, "localhost", config.DBHost)
	assert.Equal(t, 5432, config.DBPort)
	assert.Equal(t, "testuser", config.DBUser)
	assert.Equal(t, "testpassword", config.DBPassword)
	assert.Equal(t, "testdb", config.DBName)
	assert.Equal(t, "disable", config.DBSSLMode)
	assert.Equal(t, "supersecret", config.JWTSecret)
	assert.Equal(t, "3600", config.JWTExpiration)
}

func TestLoadConfigInvalidDBPort(t *testing.T) {
	// Set invalid DB_PORT
	os.Setenv("DB_PORT", "invalid_number")

	// Call LoadConfig
	_, err := LoadConfig()

	// Expect error because DB_PORT is invalid
	assert.Error(t, err, "LoadConfig should return an error for invalid DB_PORT")
}

func TestLoadConfigMissingEnvFile(t *testing.T) {
	// Simulate missing .env file by unsetting variables
	os.Clearenv()

	// Call LoadConfig
	_, err := LoadConfig()

	// Expect an error because .env file is missing
	assert.Error(t, err, "LoadConfig should return an error if .env is missing")
}
