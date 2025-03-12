package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBPort        int
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	JWTSecret     string
	JWTExpiration string
	DBDriver      string
	DBSource      string
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig() (*Config, error) {
	if os.Getenv("TEST_ENV") == "true" {
		if err := godotenv.Load(".env.test"); err != nil {
			return nil, err
		}
	} else {
		if err := godotenv.Load(); err != nil {
			return nil, err
		}
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	return &Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        dbPort,
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBSSLMode:     os.Getenv("DB_SSL_MODE"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiration: os.Getenv("JWT_EXPIRATION"),
		DBDriver:      os.Getenv("DB_DRIVER"),
		DBSource:      os.Getenv("DB_SOURCE"),
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
	}, nil
}
