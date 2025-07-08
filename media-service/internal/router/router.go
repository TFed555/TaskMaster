package router

import (
	"media-service/internal/controllers"
	"shared/middleware"
	_ "shared/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Router struct {
	ChiRouter *chi.Mux
	Port      string
}

func InitNewRouter(mediaController controllers.MediaController, authMiddleware middleware.AuthMiddleware) Router {

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:3001"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Control-Allow-Origin"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/api/test", authMiddleware.SetAuthMiddleware(mediaController.Test))
	router.Get("/api/save", mediaController.Generate)
	router.Post("/api/save/avatar", authMiddleware.SetAuthMiddleware(mediaController.SaveAvatar))

	return Router{
		ChiRouter: router,
		Port:      ":8084",
	}
}
