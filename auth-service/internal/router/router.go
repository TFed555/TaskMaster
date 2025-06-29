package router

import (
	"auth-service/internal/controllers"
	"auth-service/internal/middleware"
	// "log"
	// "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Router struct {
	ChiRouter *chi.Mux
	Port string
}

func InitNewRouter(authController *controllers.AuthController) (*Router) {

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3001", "http://127.0.0.1:3001"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    }))

	// router.Post("/", middleware.AuthMiddleware(authController.Test))
	router.Post("/api/registration", authController.Register)
	router.Post("/api/login", authController.Authorize)
	router.Post("/api/refresh", authController.Refresh)

	router.Get("/api/testSecure", middleware.AuthMiddleware(authController.Test))
	// router.Post("/testJWT", authController.TestJWT)

	return &Router{
		ChiRouter: router,
		Port: ":8080",
	}
}