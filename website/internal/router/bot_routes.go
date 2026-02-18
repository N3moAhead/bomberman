package router

import (
	"net/http"
	"strconv"

	"github.com/N3moAhead/bomberman/website/internal/db"
	"github.com/N3moAhead/bomberman/website/internal/models"
	"github.com/N3moAhead/bomberman/website/internal/templates/bots"
	"github.com/N3moAhead/bomberman/website/internal/viewmodels"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func botRoutes(botRouter chi.Router) {
	botRouter.Get("/{botID}", func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value(userContextKey).(*models.User)
		botIDParam := chi.URLParam(r, "botID")

		if botIDParam == "" {
			http.NotFound(w, r)
			return
		}

		botID, err := strconv.Atoi(botIDParam)
		if err != nil {
			log.Error("Error converting ID '%s' to integer: %+v\n", botIDParam, err)
			http.Error(w, "Invalid ID format. Must be an integer.", http.StatusBadRequest)
			return
		}

		bot, err := db.GetBotByID(uint(botID))
		if err != nil {
			http.Error(w, "failed to get bot", http.StatusInternalServerError)
			return
		}
		matches, err := db.GetMatchesForBot(bot, 1, 50)
		if err != nil {
			http.Error(w, "failed to get bot details", http.StatusInternalServerError)
		}

		vm, err := viewmodels.NewBotDetail(bot, matches)

		botDetailTemplate := bots.Detail(csrf.Token(r), user, vm)
		err = botDetailTemplate.Render(r.Context(), w)
		renderError(err, w)
	})

	// --- Secured Routes ---
	botRouter.Group(func(authRouter chi.Router) {
		authRouter.Use(authMiddleware)
		authRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			user, _ := r.Context().Value(userContextKey).(*models.User)
			userBots, _ := db.GetBotsForUser(user)

			b := bots.Overview(user, csrf.Token(r), userBots)
			err := b.Render(r.Context(), w)
			renderError(err, w)
		})

		authRouter.Get("/new", func(w http.ResponseWriter, r *http.Request) {
			user, _ := r.Context().Value(userContextKey).(*models.User)
			form := &viewmodels.NewBotForm{}
			b := bots.NewBot(user, csrf.Token(r), form)
			err := b.Render(r.Context(), w)
			renderError(err, w)
		})

		authRouter.Post("/new", func(w http.ResponseWriter, r *http.Request) {
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

			user, _ := r.Context().Value(userContextKey).(*models.User)
			if !form.Validate() {
				w.WriteHeader(http.StatusBadRequest)
				b := bots.NewBot(user, csrf.Token(r), &form)
				err := b.Render(r.Context(), w)
				renderError(err, w)
				return
			}

			err := db.CreateBot(form.ToDbModel(user.ID))
			if err != nil {
				form.Errors["db_error"] = "Error while saving to the database please try again later."
				w.WriteHeader(http.StatusInternalServerError)
				b := bots.NewBot(user, csrf.Token(r), &form)
				err := b.Render(r.Context(), w)
				renderError(err, w)
				return
			}

			http.Redirect(w, r, "/bots", http.StatusFound)
		})

	})
}
