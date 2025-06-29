package main

import (
	"auth-service/internal/config/dbconfig"
	"auth-service/internal/controllers"
	"auth-service/internal/middleware"
	"auth-service/internal/migrations"
	"auth-service/internal/repository"
	"auth-service/internal/router"
	"auth-service/internal/services"
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)


func main() {
	dbconf := dbconfig.NewDBConfig()

	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			dbconf.User,
			dbconf.Password,
			dbconf.Host,
			dbconf.Port,
			dbconf.Name,
			dbconf.SSLMode)

	err := migrations.RunMigration(conn, "internal/migrations")
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		dbconf.Host,
		dbconf.Port,
		dbconf.User,
		dbconf.Name,
		dbconf.Password,
		dbconf.SSLMode))
	if err != nil {
		log.Fatal("DB connection error: ", err)
	}

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	authService := services.NewAuthService(userRepo, tokenRepo)
	authController := controllers.NewAuthController(authService)

	authMiddleware := middleware.NewAuthMiddleware(authService)

	router := router.InitNewRouter(authController, authMiddleware)

	log.Printf("Starting server on %s \n", router.Port)
	if err := http.ListenAndServe(router.Port, router.ChiRouter); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}