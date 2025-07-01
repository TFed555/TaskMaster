package router

import (
	// "log"
	// "net/http"

	"notes-service/internal/controllers"
	// "notes-service/internal/middleware"
	_"shared/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Router struct {
	ChiRouter *chi.Mux
	Port string
}

func InitNewRouter(notesController *controllers.NotesController) (*Router) {

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


	// router.Post("/api/registration", authController.Register)
	// router.Post("/api/login", authController.Authorize)
	// router.Post("/api/refresh", authController.Refresh)
	router.Get("/api/test", notesController.Test)
	router.Get("/api/todos?createdAt=(date YYYY-MM-DD)&filter=(after | before)&offset=(int)&limit=(int)", notesController.GetTodos)

	return &Router{
		ChiRouter: router,
		Port: ":8082",
	}
}