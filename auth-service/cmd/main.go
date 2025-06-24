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
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_"gorm.io/driver/postgres"
	_ "gorm.io/gorm"
	"auth-service/internal/migrations"
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
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	//перенести в отдельный файл
	router := chi.NewRouter()
	router.Post("/reg", authController.Register)
	router.Post("/auth", authController.Authorize)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

