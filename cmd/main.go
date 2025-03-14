package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/danisasmita/customer-search/internal/config"
	"github.com/danisasmita/customer-search/internal/handler"
	"github.com/danisasmita/customer-search/internal/repository"
	"github.com/danisasmita/customer-search/internal/service"
	"github.com/danisasmita/customer-search/pkg/database"
	"github.com/danisasmita/customer-search/pkg/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	migrate := flag.Bool("migrate", false, "Run database migrations")
	seed := flag.Bool("seed", false, "Seed database with initial data")
	flag.Parse()

	// Ambil environment variable
	cfg := config.Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "user"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "customer_search"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
	}

	// Print konfigurasi untuk debugging
	fmt.Println("Database Configuration:")
	fmt.Printf("Host: %s\n", cfg.DBHost)
	fmt.Printf("User: %s\n", cfg.DBUser)
	fmt.Printf("Database: %s\n", cfg.DBName)
	fmt.Printf("Port: %d\n", cfg.DBPort)

	// Buat DSN dari environment variables
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if *migrate {
		err = database.AutoMigrate(db)
		if err != nil {
			log.Fatalf("failed to migrate database: %v", err)
		}
		log.Println("Migration completed successfully!")
	}

	if *seed {
		err = database.SeedData(db)
		if err != nil {
			log.Fatalf("failed to seed database: %v", err)
		}
		log.Println("Seeding completed successfully!")
	}

	customerRepo := repository.NewCustomerRepository(db)
	customerService := service.NewCustomerService(customerRepo)
	customerHandler := handler.NewCustomerHandler(customerService)

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
	}))

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	authorized := r.Group("/")
	authorized.Use(middleware.JWTAuth())
	{
		authorized.GET("/customers", customerHandler.SearchByName)
	}

	serverAddress := getEnv("SERVER_ADDRESS", ":8080")
	log.Printf("Starting server on %s\n", serverAddress)
	r.Run(serverAddress)
}

// Fungsi untuk mendapatkan environment variable dengan nilai default
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Fungsi untuk mendapatkan environment variable sebagai integer dengan nilai default
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("failed to convert %s to int: %v", key, err)
	}
	return value
}
