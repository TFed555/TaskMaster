package main

import (
	"fmt"
	"log"
	_ "net"
	_"time"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
    "gorm.io/driver/postgres"
	"github.com/go-chi/chi/v5"
	"net/http"
	"auth-service/internal/controllers"
	"auth-service/internal/repository"	
	"auth-service/internal/services"
)

func main() {
	//параметры будут подаваться из env или cfg (local.yaml)
	connStr := fmt.Sprintf("host=127.0.0.1 port=5432 user=user dbname=mydb password=password sslmode=disable")
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		fmt.Print("DB error", err)
	}

	userRepo := repository.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	//перенести в отдельный файл
	router := chi.NewRouter()
	router.Post("/reg", authController.Register)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

