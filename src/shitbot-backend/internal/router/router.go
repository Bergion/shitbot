package router

import (
	"net/http"
	"shitbot/internal/auth"
	"shitbot/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func NewRouter(userHandler *handlers.UserHandler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/health"))
	r.Use(middleware.NoCache)
	r.Use(auth.TelegramAuth)

	r.Get("/users", userHandler.GetUser)
	r.Post("/users", userHandler.CreateUser)
	return cors.AllowAll().Handler(r)
}
