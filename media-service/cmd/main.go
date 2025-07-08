package main

import (
	"log"
	"media-service/internal/config"
	"media-service/internal/controllers"
	"media-service/internal/repository"
	"media-service/internal/router"
	"media-service/internal/services"
	"net/http"
	"shared/config/dbconfig"
	"shared/middleware"
	"media-service/internal/grpc"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	dbconf := dbconfig.NewDBConfig()
	conn := dbconf.CreateConnection()
	if conn == "" {
		log.Fatal("Failed to create connection")
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

	minioClient, err := config.NewMinioClient()
	if err != nil {
		log.Panic("Failed to connect to minio")
	}
	grpcAuthService := middleware.NewGRPCAuthService(authCon)
	authMiddleware := middleware.NewAuthMiddleware(grpcAuthService)
	mediaRepo := repository.NewMediaRepository(db, minioClient)
	mediaService := services.NewMediaService(mediaRepo)
	mediaController := controllers.NewMediaController(mediaService)

	go func() {
		if err:= grpc_media.StartGRPCServer(mediaService, ":50050"); err != nil {
			log.Fatalf("GrpcServer otkisaet:%v", err)
		}
	}()

	router := router.InitNewRouter(mediaController, authMiddleware)

	log.Printf("Starting server on %s \n", router.Port)
	if err := http.ListenAndServe(router.Port, router.ChiRouter); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	
}