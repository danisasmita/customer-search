package main

import (
	"flag"
	"log"
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

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dsn := "host=" + cfg.DBHost + " user=" + cfg.DBUser + " password=" + cfg.DBPassword + " dbname=" + cfg.DBName + " port=" + strconv.Itoa(cfg.DBPort) + " sslmode=disable TimeZone=Asia/Jakarta"
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
	customerHandler := handler.NewCustomerHandler(*customerService)

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

	r.Run(cfg.ServerAddress)
}
