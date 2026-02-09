package router

import (
	"context"
	"net/http"
	"strings"

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

	FileServer(r, "/static", http.Dir("./static"))

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
