package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/N3moAhead/bomberman/website/internal/cfg"
	"github.com/N3moAhead/bomberman/website/internal/templates/dashboard"
	"github.com/N3moAhead/bomberman/website/internal/templates/home"
	"github.com/N3moAhead/bomberman/website/internal/templates/scoreboard"
	"github.com/N3moAhead/bomberman/website/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

var log = logger.New("[Router]")

// store wird eine package-level Variable, damit die Middleware darauf zugreifen kann.
var store sessions.Store

const appSessionName = "bomberman-session"

func Start(cfg *cfg.Config) {

	// --- AUTH ---
	// WICHTIG: Der `key` sollte in Produktion aus einer Umgebungsvariable kommen,
	// damit die Sessions einen Server-Neustart überleben.
	key := "ein-sehr-geheimer-key-der-mindestens-32-bytes-lang-ist-dev-only" // 64 bytes
	maxAge := 86400 * 30                                                     // 30 Tage
	isProd := false                                                          // In Produktion auf true setzen

	// Wir stellen zurück auf CookieStore.
	cookieStore := sessions.NewCookieStore([]byte(key))
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   isProd,
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

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	FileServer(r, "/static", http.Dir("./static"))

	r.Get("/auth/{provider}", func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")
		r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
		gothic.BeginAuthHandler(w, r)
	})

	r.Get("/auth/{provider}/callback", func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")
		r = r.WithContext(context.WithValue(r.Context(), "provider", provider))

		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintf(w, "Fehler bei der Authentifizierung: %v", err)
			return
		}

		appSession, err := store.Get(r, appSessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		appSession.Values["user_id"] = user.UserID
		appSession.Values["nickname"] = user.NickName
		appSession.Values["avatar_url"] = user.AvatarURL
		if err := appSession.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("User '%s' logged in successfully", user.NickName)
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	})

	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		appSession, err := store.Get(r, appSessionName)
		if err == nil {
			appSession.Options.MaxAge = -1
			_ = appSession.Save(r, w)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		h := home.Home("Lukas")
		h.Render(context.Background(), w)
	})

	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)

		r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
			appSession, _ := store.Get(r, appSessionName)
			nickname, ok := appSession.Values["nickname"].(string)
			if !ok {
				nickname = "User"
			}
			d := dashboard.Dashboard(nickname)
			d.Render(context.Background(), w)
		})
	})

	r.Get("/scoreboard", func(w http.ResponseWriter, r *http.Request) {
		s := scoreboard.Scoreboard()
		s.Render(context.Background(), w)
	})

	log.Info("Starting website on port %s", cfg.Port)
	err := http.ListenAndServe(cfg.Port, r)
	if err != nil {
		log.Error("failed to start website: %v", err)
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, appSessionName)

		if err != nil || session.IsNew || session.Values["user_id"] == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not allow url parameters")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
