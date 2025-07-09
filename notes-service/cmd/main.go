package main

import (
	"notes-service/pkg/app"
	// "log"
	// "net/http"
	// "notes-service/internal/controllers"
	// "notes-service/internal/migrations"
	// "notes-service/internal/repository"
	// "notes-service/internal/router"
	// "notes-service/internal/services"
	// "shared/config/dbconfig"
	// "shared/middleware"

	// "google.golang.org/grpc"
)

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
