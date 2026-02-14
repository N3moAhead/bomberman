package router

import (
	"encoding/json"
	"net/http"

	"github.com/N3moAhead/bomberman/website/internal/db"
	"github.com/N3moAhead/bomberman/website/internal/models"
	"github.com/N3moAhead/bomberman/website/internal/templates/matches"
	"github.com/N3moAhead/bomberman/website/internal/viewmodels"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
)

func MatchRoutes() chi.Router {
	matchRouter := chi.NewRouter()

	matchRouter.Get("/{matchID}", handleGetMatch)

	matchRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		user, _ := r.Context().Value(UserContextKey).(*models.User)
		botMatches, _ := db.GetMatches(0, 50)
		s := matches.Matches(csrf.Token(r), user, botMatches)
		s.Render(r.Context(), w)
	})

	return matchRouter
}

func handleGetMatch(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value(UserContextKey).(*models.User)
	matchID := chi.URLParam(r, "matchID")
	if matchID == "" {
		// Or handle error appropriately
		http.NotFound(w, r)
		return
	}

	match, err := db.GetMatchByMatchID(matchID)
	if err != nil {
		// handle error
		http.Error(w, "failed to get match", http.StatusInternalServerError)
		return
	}

	if len(match.History) == 0 {
		// Or render a page indicating history is not available
		http.Error(w, "match history not available", http.StatusNotFound)
		return
	}

	vm, err := viewmodels.NewMatchDetail(match)
	if err != nil {
		http.Error(w, "failed to create view model", http.StatusInternalServerError)
		return
	}

	historyJson, err := json.Marshal(vm.History)
	if err != nil {
		http.Error(w, "failed to marshal history", http.StatusInternalServerError)
		return
	}

	matches.Details(csrf.Token(r), user, vm, string(historyJson)).Render(r.Context(), w)
}
