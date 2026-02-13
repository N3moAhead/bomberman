package router

import (
	"net/http"

	"github.com/N3moAhead/bomberman/website/internal/models"
	"github.com/N3moAhead/bomberman/website/internal/templates/bots"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func botRoutes(botRouter chi.Router) {
	botRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value(UserContextKey).(*models.User)
		b := bots.Overview(user, csrf.Token(r))
		b.Render(r.Context(), w)
	})

	botRouter.Get("/new", func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value(UserContextKey).(*models.User)
		b := bots.NewBot(user, csrf.Token(r))
		b.Render(r.Context(), w)
	})

	botRouter.Post("/new", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/bots", http.StatusFound)
	})
}
