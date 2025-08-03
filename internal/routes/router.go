package routes

import (
	"monad-indexer/internal/handlers"
	"monad-indexer/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Monad Dev Portfolio API"))
	})

	r.Route("/projects", func(r chi.Router) {
		r.Get("/", handlers.GetAllProjects)
		r.With(middleware.IsAdmin).Post("/", handlers.CreateProject)
	})

	// r.Route("/devs", func(r chi.Router) {
	// 	r.With(middleware.IsAdmin).Post("/", handlers.CreateDev)
		
	// })

	r.Get("/devs", handlers.GetAllDevs)

	r.Get("/dev", handlers.GetDev)
	
	return r
}