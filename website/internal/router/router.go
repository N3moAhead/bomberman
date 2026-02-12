package router

import (
	"net/http"
	"strings"

	"github.com/N3moAhead/bomberman/website/internal/cfg"
	"github.com/N3moAhead/bomberman/website/internal/templates/dashboard"
	"github.com/N3moAhead/bomberman/website/internal/templates/home"
	"github.com/N3moAhead/bomberman/website/internal/templates/leaderboard"
	"github.com/N3moAhead/bomberman/website/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
)

var log = logger.New("[Router]")

// store wird eine package-level Variable, damit die Middleware darauf zugreifen kann.
var store sessions.Store

const appSessionName = "bomberman-session"

func Start(cfg *cfg.Config) {
	log.Info("Trusted origin: %s", cfg.BaseURL)

	authSetup(cfg)

	r := chi.NewRouter()

	// --- Middlewares ---
	r.Use(middleware.Logger)
	csrfMiddleware := csrf.Protect(
		[]byte(cfg.CSRFAuthKey),
		csrf.Secure(cfg.IsProduction),
	)
	r.Use(csrfMiddleware)

	// Serving static files
	FileServer(r, "/static", http.Dir("./static"))

	/// --- Auth Routes ---
	r.Get("/auth/{provider}", githubLogin)
	r.Get("/auth/{provider}/callback", githubLoginCallback)
	r.Post("/logout", logout)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		h := home.Home(csrf.Token(r))
		h.Render(r.Context(), w)
	})

	r.Get("/leaderboard", func(w http.ResponseWriter, r *http.Request) {
		s := leaderboard.Leaderboard(csrf.Token(r))
		s.Render(r.Context(), w)
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
			d := dashboard.Dashboard(nickname, csrf.Token(r))
			d.Render(r.Context(), w)
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
