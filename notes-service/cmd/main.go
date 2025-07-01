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
	"github.com/jmoiron/sqlx"
	_"shared/middleware"
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

	notesRepo := repository.NewNotesRepository(db)
	notesService := services.NewNotesService(notesRepo)
	notesController := controllers.NewNotesController(notesService)
	// authMiddleware := middleware.NewAuthMiddleware()
	router := router.InitNewRouter(notesController)

	log.Printf("Starting server on %s \n", router.Port)
	if err := http.ListenAndServe(router.Port, router.ChiRouter); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}