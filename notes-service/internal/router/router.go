package router

import (
	"notes-service/internal/controllers"
	"shared/middleware"

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

	// router.Use(authMiddleware.SetAuthMiddleware)
	router.Get("/api/test", notesController.Test)

	router.Get("/api/todos", notesController.GetTodos)
	router.Get("/api/todos/archived", notesController.GetArchivedTodos)
	router.Delete("/api/todos/archived/{id}", notesController.DeleteTodo)
	router.Put("/api/todos/archived/{id}", notesController.RestoreTodo)
	router.Patch("/api/todos/{id}", notesController.UpdateTodo)
	router.Delete("/api/todos/{id}", notesController.ArchiveTodo)
	router.Get("/api/todos/{id}", notesController.GetOneTodo)
	router.Post("/api/todos", notesController.CreateTodo)
	router.Get("/api/todos/history", notesController.GetAuditTrail)

	router.Post("/api/tags", notesController.CreateTag)
	router.Get("/api/tags", notesController.GetTags)
	router.Patch("/api/tags/{id}", notesController.UpdateTag)
	router.Delete("/api/tags/{id}", notesController.DeleteTag)
	router.Post("/api/todos/{id}/tags", notesController.AddTagToTodo)
	router.Delete("/api/todos/{id}/tags/{tag_id}", notesController.ReduceTag)

	router.Post("/api/todos/plan", notesController.CreatePlanTodo)

	return Router{
		ChiRouter: router,
		Port:      ":8082",
	}
}
