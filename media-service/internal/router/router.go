package router

import (
	"media-service/internal/controllers"
	"net/http"
	"shared/middleware"
	_ "shared/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Router struct {
	ChiRouter *chi.Mux
	Port      string
}

func InitNewRouter(mediaController controllers.MediaController, authMiddleware middleware.AuthMiddleware,
					authCon *grpc.ClientConn) Router {

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
	router.Get("/api/images/avatar", mediaController.Generate)
	// router.Post("/api/images/avatar", authMiddleware.SetAuthMiddleware(mediaController.SaveAvatar))
	router.Post("/api/images/avatar", mediaController.SaveAvatar)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			if authCon.GetState() != connectivity.Ready {
				w.WriteHeader(503)
				return
			}

			w.WriteHeader(200)
	})

	return Router{
		ChiRouter: router,
		Port:      ":8084",
	}
}
