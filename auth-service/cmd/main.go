package main

import (
	"auth-service/internal/controllers"
	"auth-service/internal/repository"
	"auth-service/internal/services"
	"fmt"
	"log"
	_ "net"
	"net/http"
	"os"
	_ "os"
	_ "time"

	"auth-service/internal/middleware"
	"auth-service/internal/migrations"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"
)

func initDb() (string, string, string, string, string) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Can't load .env file")
	}
	dbUser, exists := os.LookupEnv("DB_USER")
	if !exists {
		log.Printf("Can't lookup for DB_USER")
	}
	dbPass, exists := os.LookupEnv("DB_PASSWORD")
	if !exists {
		log.Printf("Can't lookup for DB_USER")
	}
	dbHost, exists := os.LookupEnv("DB_HOST")
	if !exists {
		log.Printf("Can't lookup for DB_USER")
	}
	dbPort, exists := os.LookupEnv("DB_PORT")
	if !exists {
		log.Printf("Can't lookup for DB_USER")
	}
	dbName, exists := os.LookupEnv("DB_NAME")
	if !exists {
		log.Printf("Can't lookup for DB_USER")
	}
	return dbUser, dbPass, dbHost, dbPort, dbName
}

func main() {
	//параметры будут подаваться из env или cfg (local.yaml)
	// conn := "postgres://user:password@localhost:5432/tasksdb?sslmode=disable"
	dbUser, dbPass, dbHost, dbPort, dbName := initDb()

	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			dbUser,
			dbPass,
			dbHost,
			dbPort,
			dbName)

	err := migrations.RunMigration(conn, "internal/migrations")
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=localhost port=5432 user=user dbname=tasksdb password=password sslmode=disable",))
	if err != nil {
		log.Fatal("DB connection error: ", err)
	}

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	authService := services.NewAuthService(userRepo, tokenRepo)
	authController := controllers.NewAuthController(authService)

	//перенести в отдельный файл
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3001", "http://127.0.0.1:3001"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    }))

	// router.Post("/", middleware.AuthMiddleware(authController.Test))
	router.Post("/api/registration", authController.Register)
	router.Post("/api/login", authController.Authorize)
	router.Post("/api/refresh", authController.Refresh)

	router.Get("/api/testSecure", middleware.AuthMiddleware(authController.Test))
	// router.Post("/testJWT", authController.TestJWT)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}