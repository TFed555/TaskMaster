package main

import (
	// "auth-service/internal/controllers"
	// "auth-service/internal/grpc"
	// "shared/config/dbconfig"

	// "auth-service/internal/migrations"
	// "auth-service/internal/pkg/jwt"
	// "auth-service/internal/repository"
	// "auth-service/internal/router"
	// "auth-service/internal/services"
	// "log"
	// "net/http"
	// "shared/middleware"
	"auth-service/pkg/app"
)


func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
	// dbconf := dbconfig.NewDBConfig()
	// conn := dbconf.CreateConnection()
	// err := migrations.RunMigration(conn, "internal/migrations")
	// if err != nil {
	// 	log.Fatalf("Migration error: %v", err)
	// }
	// db, err := dbconf.CreateDSN()
	// if err != nil {
	// 	log.Fatal("DB connection error: ", err)
	// }

	// userRepo := repository.NewUserRepository(db)
	// tokenRepo := repository.NewTokenRepository(db)
	// jwtFunc := jwt.NewJWTFunctional()
	// authService := services.NewAuthService(userRepo, tokenRepo, jwtFunc, nil)
	// authController := controllers.NewAuthController(authService)
	// authMiddleware := middleware.NewAuthMiddleware(authService)

	// go func() {
	// 	if err:= grpc_auth.StartGRPCServer(authService, ":50051"); err != nil {
	// 		log.Fatalf("GrpcServer otkisaet:%v", err)
	// 	}
	// }()
	
	// router := router.InitNewRouter(authController, authMiddleware)

	// log.Printf("Starting server on %s \n", router.Port)
	// if err := http.ListenAndServe(router.Port, router.ChiRouter); err != nil {
	// 	log.Fatalf("Failed to start server: %v", err)
	// }
}