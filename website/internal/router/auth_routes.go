package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/N3moAhead/bombahead/website/internal/cfg"
	"github.com/N3moAhead/bombahead/website/internal/db"
	"github.com/N3moAhead/bombahead/website/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

func authSetup(cfg *cfg.Config) {
	// --- AUTH Setup ---
	maxAge := 86400 * 30
	cookieStore := sessions.NewCookieStore([]byte(cfg.SessionSecret))
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   cfg.IsProduction,
	}
	store = cookieStore
	gothic.Store = store

	goth.UseProviders(
		github.New(
			cfg.GithubCLientId,
			cfg.GithubClientSecret,
			cfg.NextAuthUrl,
			"read:user",
		),
	)
}

func githubLogin(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), providerContextKey, provider))
	gothic.BeginAuthHandler(w, r)
}

func githubLoginCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), providerContextKey, provider))

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
	http.Redirect(w, r, "/bots", http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	appSession, err := store.Get(r, appSessionName)
	if err == nil {
		appSession.Options.MaxAge = -1
		_ = appSession.Save(r, w)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func isAuthenticated(r *http.Request) bool {
	session, err := store.Get(r, appSessionName)
	if err != nil || session.IsNew || session.Values["user_id"] == nil {
		return false
	}
	return true
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isAuthenticated(r) {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
