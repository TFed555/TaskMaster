package router

import (
	// "log"
	// "net/http"

	"notes-service/internal/controllers"
	// "notes-service/internal/middleware"
	"shared/middleware"
	_ "shared/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Router struct {
	ChiRouter *chi.Mux
	Port      string
}

func InitNewRouter(notesController controllers.NotesController, authMiddleware middleware.AuthMiddleware) Router {

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

	router.Get("/api/test", authMiddleware.SetAuthMiddleware(notesController.Test))
	router.Get("/api/todos", authMiddleware.SetAuthMiddleware(notesController.GetTodos))
	router.Get("/api/todos/archived", authMiddleware.SetAuthMiddleware(notesController.GetArchivedTodos))
	router.Delete("/api/todos/archived/{id}", authMiddleware.SetAuthMiddleware(notesController.DeleteTodo))
	router.Put("/api/todos/archived/{id}", authMiddleware.SetAuthMiddleware(notesController.RestoreTodo))
	router.Get("/api/todos/{id}", authMiddleware.SetAuthMiddleware(notesController.GetOneTodo))
	router.Delete("/api/todos/{id}", authMiddleware.SetAuthMiddleware(notesController.ArchiveTodo))
	router.Post("/api/todos", authMiddleware.SetAuthMiddleware(notesController.CreateTodo))
	router.Patch("/api/todos", authMiddleware.SetAuthMiddleware(notesController.UpdateTodo))

	router.Post("/api/tags", authMiddleware.SetAuthMiddleware(notesController.CreateTag))
	router.Get("/api/tags", authMiddleware.SetAuthMiddleware(notesController.GetTags))
	router.Patch("/api/tags/{id}", authMiddleware.SetAuthMiddleware(notesController.UpdateTag))
	router.Delete("/api/tags/{id}", authMiddleware.SetAuthMiddleware(notesController.DeleteTag))
	router.Post("/api/todos/{id}/tags", authMiddleware.SetAuthMiddleware(notesController.AddTagToTodo))
	router.Delete("/api/todos/{id}/tags/{tag_id}", authMiddleware.SetAuthMiddleware(notesController.ReduceTag))

	return Router{
		ChiRouter: router,
		Port:      ":8082",
	}
}
