package router

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	// "github.com/urfave/negroni"

	lHandler "github.com/alexsniffin/website/internal/server/handlers/logging"
	"github.com/alexsniffin/website/internal/server/models"
)

// Creates Chi based multiplexer router with middleware
func NewRouter(cfg models.HTTPRouterConfig, logger zerolog.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(lHandler.NewHandlerFunc(logger))
	r.Use(middleware.Timeout(time.Duration(cfg.TimeoutSec) * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   cfg.AllowedMethods,
		AllowedHeaders:   cfg.AllowedHeaders,
		AllowCredentials: false,
	}))

	FileServer(r)

	return r
}

func FileServer(router *chi.Mux) {
	root := "./assets"
	fs := http.FileServer(http.Dir(root))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.RedirectHandler("/404.html", 301).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
