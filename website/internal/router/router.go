package router

import (
	"context"
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
	key := "ein-sehr-geheimer-key-der-mindestens-32-bytes-lang-ist-dev-only"
	maxAge := 86400 * 30
	isProd := false

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
	r.Use(middleware.Recoverer)
	FileServer(r, "/static", http.Dir("./static"))

	/// --- Auth Routes ---
	r.Get("/auth/{provider}", githubLogin)
	r.Get("/auth/{provider}/callback", githubLoginCallback)
	r.Get("/logout", logout)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		h := home.Home()
		h.Render(context.Background(), w)
	})

	r.Get("/scoreboard", func(w http.ResponseWriter, r *http.Request) {
		s := scoreboard.Scoreboard()
		s.Render(context.Background(), w)
	})

	// --- Secured Routes ---
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

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

	log.Info("Starting website on port %s", cfg.Port)
	err := http.ListenAndServe(cfg.Port, r)
	if err != nil {
		log.Error("failed to start website: %v", err)
	}
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
