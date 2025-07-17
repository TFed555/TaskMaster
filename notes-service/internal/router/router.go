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

	router.Group(func(r chi.Router) {
			r.Use(authMiddleware.SetAuthMiddleware)
			r.Get("/api/test", notesController.Test)

			r.Get("/api/todos", notesController.GetTodos)
			r.Get("/api/todos/archived", notesController.GetArchivedTodos)
			r.Delete("/api/todos/archived/{id}", notesController.DeleteTodo)
			r.Put("/api/todos/archived/{id}", notesController.RestoreTodo)
			r.Patch("/api/todos/{id}", notesController.UpdateTodo)
			r.Delete("/api/todos/{id}", notesController.ArchiveTodo)
			r.Get("/api/todos/{id}", notesController.GetOneTodo)
			r.Post("/api/todos", notesController.CreateTodo)
			r.Get("/api/todos/history", notesController.GetAuditTrail)

			r.Post("/api/tags", notesController.CreateTag)
			r.Get("/api/tags", notesController.GetTags)
			r.Patch("/api/tags/{id}", notesController.UpdateTag)
			r.Delete("/api/tags/{id}", notesController.DeleteTag)
			r.Post("/api/todos/{id}/tags", notesController.AddTagToTodo)
			r.Delete("/api/todos/{id}/tags/{tag_id}", notesController.ReduceTag)
	})
	// r.Post("/api/todos/plan", notesController.CreatePlanTodo)

	router.Post("/api/todos/plan", notesController.CreatePlanTodo)

	return Router{
		ChiRouter: router,
		Port:      ":8082",
	}
}
