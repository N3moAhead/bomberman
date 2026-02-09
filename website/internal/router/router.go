package router

import (
	"context"
	"net/http"

	"github.com/N3moAhead/bomberman/website/internal/cfg"
	"github.com/N3moAhead/bomberman/website/internal/templates/home"
	"github.com/N3moAhead/bomberman/website/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var log = logger.New("[Router]")

func Start(cfg *cfg.Config) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		h := home.Home("Lukas")
		h.Render(context.Background(), w)
	})

	log.Info("Starting website on port %s", cfg.Port)
	err := http.ListenAndServe(cfg.Port, r)
	if err != nil {
		log.Error("failed to start website: %v", err)
	}
}
