// backend/api/routes/routes.go
package routes

import (
	"net/http"

	"jazz/backend/handlers"
	"jazz/backend/pkg/middlewares"

	"github.com/go-chi/chi/v5"
)

// SetupRoutes sets up the application's routes.
func SetupRoutes() http.Handler {
	r := chi.NewRouter()

	// Public routes
	r.Post("/register", handlers.RegisterUserHandler)
	r.Post("/login", handlers.LoginHandler)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)
		r.Get("/profile", handlers.UserProfileHandler) // Exemplo de rota protegida
	})

	return r
}
