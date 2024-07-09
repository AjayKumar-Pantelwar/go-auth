package routes

import (
	userhandler "go-authentication/src/internal/interfaces/input/api/rest/handler"
	"go-authentication/src/internal/interfaces/input/api/rest/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitRoutes() http.Handler {
	router := chi.NewRouter()

	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", userhandler.Register)
		r.Post("/login", userhandler.Login)
		// r.Post("/logout", Logout)

	})

	router.Route("/user", func(r chi.Router) {
		r.Use(middleware.Authenticate)
		r.Get("/profile", userhandler.Profile)
	})

	return router
}
