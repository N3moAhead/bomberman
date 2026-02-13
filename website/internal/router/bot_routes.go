package router

import (
	"net/http"

	"github.com/N3moAhead/bomberman/website/internal/db"
	"github.com/N3moAhead/bomberman/website/internal/models"
	"github.com/N3moAhead/bomberman/website/internal/templates/bots"
	"github.com/N3moAhead/bomberman/website/internal/viewmodels"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func botRoutes(botRouter chi.Router) {
	botRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value(UserContextKey).(*models.User)
		userBots, _ := db.GetBotsForUser(user)

		b := bots.Overview(user, csrf.Token(r), userBots)
		b.Render(r.Context(), w)
	})

	botRouter.Get("/new", func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value(UserContextKey).(*models.User)
		form := &viewmodels.NewBotForm{}
		b := bots.NewBot(user, csrf.Token(r), form)
		b.Render(r.Context(), w)
	})

	botRouter.Post("/new", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		form := viewmodels.NewBotForm{
			Name:          r.FormValue("name"),
			Description:   r.FormValue("description"),
			DockerHubUrl:  r.FormValue("dockerHubUrl"),
			CreatedWithAi: r.FormValue("ai") == "on",
		}

		user, _ := r.Context().Value(UserContextKey).(*models.User)
		if !form.Validate() {
			w.WriteHeader(http.StatusBadRequest)
			b := bots.NewBot(user, csrf.Token(r), &form)
			b.Render(r.Context(), w)
			return
		}

		err := db.CreateBot(form.ToDbModel(user.ID))
		if err != nil {
			form.Errors["db_error"] = "Error while saving to the database please try again later."
			w.WriteHeader(http.StatusInternalServerError)
			b := bots.NewBot(user, csrf.Token(r), &form)
			b.Render(r.Context(), w)
			return
		}

		http.Redirect(w, r, "/bots", http.StatusFound)
	})
}
