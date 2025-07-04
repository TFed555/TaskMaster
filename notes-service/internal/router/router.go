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
	Port string
}

func InitNewRouter(notesController controllers.NotesController, authMiddleware middleware.AuthMiddleware) (Router) {

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


	
	router.Get("/api/test", authMiddleware.SetAuthMiddleware(notesController.Test))
	// router.Get("/api/todos", authMiddleware.SetAuthMiddleware(notesController.GetTodos))
	// router.Post("/api/createTodo", authMiddleware.SetAuthMiddleware(notesController.CreateTodo))
	// router.Patch("/api/updateTodo", authMiddleware.SetAuthMiddleware(notesController.UpdateTodo))
	// router.Get("/api/archivedtodos", authMiddleware.SetAuthMiddleware(notesController.GetArchivedTodos))
	// router.Post("/api/archiveTodo", authMiddleware.SetAuthMiddleware(notesController.ArchiveTodo))

	router.Get("/api/todos", notesController.GetTodos)
	router.Get("/api/archivedtodos", notesController.GetArchivedTodos)
	router.Post("/api/createTodo", notesController.CreateTodo)
	router.Post("/api/archiveTodo", notesController.ArchiveTodo)
	router.Patch("/api/updateTodo", notesController.UpdateTodo)	
	router.Get("/api/todos/{id}", notesController.GetOneTodo)

	return Router{
		ChiRouter: router,
		Port: ":8082",
	}
}