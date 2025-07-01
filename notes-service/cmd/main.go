package main

import (
	"fmt"
	"log"
	"net/http"
	"notes-service/internal/config/dbconfig"
	"notes-service/internal/controllers"
	"notes-service/internal/migrations"
	"notes-service/internal/repository"
	"notes-service/internal/router"
	"notes-service/internal/services"
	"shared/middleware"
	_ "shared/middleware"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)


func main() {

	//переделать инициализацию подключения к бд
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