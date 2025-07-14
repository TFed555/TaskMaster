package router

import (
	"notes-service/internal/controllers"
	notesMiddleware "notes-service/internal/middleware"
	"shared/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Router struct {
	ChiRouter *chi.Mux
	Port      string
}

func InitNewRouter(notesController controllers.NotesController, authMiddleware middleware.AuthMiddleware, notesMiddleware notesMiddleware.NotesMiddleware) Router {

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

	router.Use(authMiddleware.SetAuthMiddleware)

	router.Route("/api/todos", func(r chi.Router) {
		// r.Use(notesMiddleware.SetNotesMIddleware)
		r.Get("/", notesController.GetTodos)
		r.Get("/archived", notesController.GetArchivedTodos)
		notesGroup := r.With(notesMiddleware.SetNotesMIddleware)
		{
			notesGroup.Delete("/archived/{id}", notesController.DeleteTodo)
			notesGroup.Put("/archived/{id}", notesController.RestoreTodo)
			notesGroup.Patch("/{id}", notesController.UpdateTodo)
			notesGroup.Delete("/{id}", notesController.ArchiveTodo)
		}

		r.Get("/{id}", notesController.GetOneTodo)
		r.Post("/", notesController.CreateTodo)
		r.Get("/history", notesController.GetAuditTrail)
	})

	router.Post("/api/tags", notesController.CreateTag)
	router.Get("/api/tags", notesController.GetTags)
	router.Patch("/api/tags/{id}", notesController.UpdateTag)
	router.Delete("/api/tags/{id}", notesController.DeleteTag)
	router.Post("/api/todos/{id}/tags", notesController.AddTagToTodo)
	router.Delete("/api/todos/{id}/tags/{tag_id}", notesController.ReduceTag)

	return Router{
		ChiRouter: router,
		Port:      ":8082",
	}
}
