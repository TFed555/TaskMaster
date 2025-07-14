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

	// router.Post("/", middleware.AuthMiddleware(authController.Test))
	router.Post("/api/registration", authController.Register)
	router.Post("/api/login", authController.Authorize)
	//router.Get("/api/refresh", authController.Refresh)

	// router.Get("/api/testSecure", authController.Test)
	router.Get("/api/test/cookie", authMiddleware.SetAuthMiddleware(authController.TestCookie))
	router.Get("/api/test/middleware", authMiddleware.SetAuthMiddleware(authController.Test))
	router.Delete("/api/login", authMiddleware.SetAuthMiddleware(authController.Logout))
	// router.Post("/testJWT", authController.TestJWT)
	router.Delete("/api/user", authMiddleware.SetAuthMiddleware(authController.DeleteUser))
	router.Patch("/api/user", authMiddleware.SetAuthMiddleware(authController.UpdateUser))

	return Router{
		ChiRouter: router,
		Port: ":8080",
	}
}