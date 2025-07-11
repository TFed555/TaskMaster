package app

import (
	"auth-service/internal/controllers"
	"auth-service/internal/grpc"
	"fmt"
	"path/filepath"
	"runtime"
	"shared/config/dbconfig"

	"auth-service/internal/migrations"
	"auth-service/internal/pkg/jwt"
	"auth-service/internal/repository"
	"auth-service/internal/router"
	"auth-service/internal/services"
	"log"
	"net/http"
	"shared/middleware"
)


func Run() error {
	dbconf := dbconfig.NewDBConfig()
	conn := dbconf.CreateConnection()
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	pathToMigrations := filepath.ToSlash(filepath.Join(basepath, "../../internal/migrations"))
	log.Print(pathToMigrations)
	// pathToMigrations = filepath.Join(".", "internal", "migrations")
	err := migrations.RunMigration(conn, "/app/internal/migrations")
	if err != nil {
		log.Fatalf("Migration error: %v", err)
		return fmt.Errorf("Migration error: %v", err)
	}

	db, err := dbconf.CreateDSN()
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
		return fmt.Errorf("DB connection error: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	jwtFunc := jwt.NewJWTFunctional()
	authService := services.NewAuthService(userRepo, tokenRepo, jwtFunc, nil)
	authController := controllers.NewAuthController(authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	go func() {
		if err:= grpc_auth.StartGRPCServer(authService, ":50051"); err != nil {
			log.Fatalf("GrpcServer otkisaet:%v", err)
		}
	}()
	
	router := router.InitNewRouter(authController, authMiddleware)

	log.Printf("Starting auth-service on %s \n", router.Port)
	if err := http.ListenAndServe(router.Port, router.ChiRouter); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return fmt.Errorf("Failed to start server: %v", err)
	}

	return nil
}