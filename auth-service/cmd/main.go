package main

import (
	"auth-service/internal/controllers"
	"auth-service/internal/repository"
	"auth-service/internal/services"
	"fmt"
	"log"
	_ "net"
	"net/http"
	_ "os"
	_ "time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_"gorm.io/driver/postgres"
	_ "gorm.io/gorm"
	"auth-service/internal/migrations"
	"auth-service/internal/middleware"
)

func main() {
	//параметры будут подаваться из env или cfg (local.yaml)
	// connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	 	os.Getenv("DB_HOST"),
	// 		os.Getenv("DB_PORT"),
	// 		os.Getenv("DB_USER"),
	// 		os.Getenv("DB_PASSWORD"),
	// 		os.Getenv("DB_NAME"))
	conn := "postgres://user:password@localhost:5432/tasksdb?sslmode=disable"
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
	router.Post("/api/register", authController.Register)
	router.Post("/api/login", authController.Authorize)
	router.Post("/api/refresh", authController.Refresh)

	router.Get("/api/testSecure", middleware.AuthMiddleware(authController.Test))
	// router.Post("/testJWT", authController.TestJWT)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

