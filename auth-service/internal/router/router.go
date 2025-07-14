package router

import (
	"auth-service/internal/controllers"
	"shared/middleware"
	// "auth-service/internal/middleware"
	// "log"
	// "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Router struct {
	ChiRouter *chi.Mux
	Port string
}

func InitNewRouter(authController controllers.AuthController, authMiddleware middleware.AuthMiddleware) (Router) {

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        // AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Control-Allow-Origin"},
		AllowedHeaders:   []string{"*"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    }))

	router.Route("/api", func(r chi.Router) {
		r.Post("/registration", authController.Register)
		r.Post("/login", authController.Authorize)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.SetAuthMiddleware)
			r.Get("/test/cookie", authController.TestCookie)
			r.Get("/test/middleware", authController.Test)
			r.Delete("/login", authController.Logout)
			r.Delete("/user", authController.DeleteUser)
			r.Patch("/user", authController.UpdateUser)
		})
	})


	return Router{
		ChiRouter: router,
		Port: ":8080",
	}
}