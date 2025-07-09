package app

import (
	"log"
	"net/http"
	"notes-service/internal/controllers"
	"notes-service/internal/repository"
	"notes-service/internal/router"
	"notes-service/internal/services"
	"shared/config/dbconfig"
	"shared/middleware"
	"fmt"
	"google.golang.org/grpc"
	_ "github.com/lib/pq"
)

func Run() error {
	dbconf := dbconfig.NewDBConfig()
	conn := dbconf.CreateConnection()
	// err := migrations.RunMigration(conn, "internal/migrations")
	// if err != nil {
	// 	log.Fatalf("Migration error: %v", err)
	// 	return fmt.Errorf("Migration error: %v", err)
	// }
	if conn == "" {
		log.Fatal("Failed to create connection")
		return fmt.Errorf("Failed to create connection")
	}
	db, err := dbconf.CreateDSN()
	if err != nil {
		log.Fatal("DB connection error: ", err)
		return fmt.Errorf("DB connection error: %v", err)
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
		return fmt.Errorf("Failed to start server: %v", err)
	}

	return nil
}
