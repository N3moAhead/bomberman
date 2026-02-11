package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/N3moAhead/bomberman/website/internal/db"
	"github.com/N3moAhead/bomberman/website/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
)

func githubLogin(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
	gothic.BeginAuthHandler(w, r)
}

func githubLoginCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintf(w, "Error during Authentication: %v", err)
		return
	}

	user, err := db.GetOrCreateUser(&models.User{
		Username:  gothUser.NickName,
		AvatarURL: gothUser.AvatarURL,
		GithubID:  gothUser.UserID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appSession, err := store.Get(r, appSessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appSession.Values["user_id"] = user.ID
	appSession.Values["nickname"] = user.Username
	appSession.Values["avatar_url"] = user.AvatarURL
	if err := appSession.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("User '%s' logged in successfully", user.Username)
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	appSession, err := store.Get(r, appSessionName)
	if err == nil {
		appSession.Options.MaxAge = -1
		_ = appSession.Save(r, w)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, appSessionName)

		if err != nil || session.IsNew || session.Values["user_id"] == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
