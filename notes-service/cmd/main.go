package main

import (
	"log"
	"net/http"
	"notes-service/internal/controllers"
	"notes-service/internal/migrations"
	"notes-service/internal/repository"
	"notes-service/internal/router"
	"notes-service/internal/services"
	"shared/config/dbconfig"
	"shared/middleware"

	"google.golang.org/grpc"
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

	authCon, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to grpc server: %v", err)
	}
	defer authCon.Close()

	grpcAuthService := middleware.NewGRPCAuthService(authCon)
	authMiddleware := middleware.NewAuthMiddleware(grpcAuthService)
	notesRepo := repository.NewNotesRepository(db)
	notesService := services.NewNotesService(notesRepo)
	notesController := controllers.NewNotesController(notesService)

	router := router.InitNewRouter(notesController, authMiddleware)

	log.Printf("Starting server on %s \n", router.Port)
	if err := http.ListenAndServe(router.Port, router.ChiRouter); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
