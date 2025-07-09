package main

import (
	"media-service/pkg/app"
	// "log"
	// "media-service/internal/config"
	// "media-service/internal/controllers"
	// "media-service/internal/repository"
	// "media-service/internal/router"
	// "media-service/internal/services"
	// "net/http"
	// "shared/config/dbconfig"
	// "shared/middleware"
	// "media-service/internal/grpc"
	// _ "github.com/lib/pq"
	// "google.golang.org/grpc"
)

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}