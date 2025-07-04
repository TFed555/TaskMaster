package main

import (
	"shared/config/dbconfig"
	"auth-service/internal/controllers"
	"auth-service/internal/grpc"

	"auth-service/internal/migrations"
	"auth-service/internal/pkg/jwt"
	"auth-service/internal/repository"
	"auth-service/internal/router"
	"auth-service/internal/services"
	"log"
	"net/http"
	"shared/middleware"
)


func main() {
	dbconf := dbconfig.NewDBConfig()
	conn := dbconf.CreateConnection()
	err := migrations.RunMigration(conn, "internal/migrations")
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}
	db, err := dbconf.CreateDSN()
	if err != nil {
		log.Fatal("DB connection error: ", err)
	}

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	jwtFunc := jwt.NewJWTFunctional()
	authService := services.NewAuthService(userRepo, tokenRepo, jwtFunc)
	authController := controllers.NewAuthController(authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	go func() {
		if err:= grpc.StartGRPCServer(authService, ":50051"); err != nil {
			log.Fatalf("GrpcServer otkisaet:%v", err)
		}
	}()

	router := router.InitNewRouter(authController, authMiddleware)

	log.Printf("Starting server on %s \n", router.Port)
	if err := http.ListenAndServe(router.Port, router.ChiRouter); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}