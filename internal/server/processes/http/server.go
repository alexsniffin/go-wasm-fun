package http

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/urfave/negroni"

	"github.com/alexsniffin/website/internal/server/models"
)

type Server struct {
	*http.Server

	cfg    models.HTTPServerConfig
	logger zerolog.Logger
}

func NewServer(cfg models.HTTPServerConfig, logger zerolog.Logger, routerHandler http.Handler) *Server {
	n := negroni.New()
	n.UseHandler(routerHandler)

	return &Server{
		&http.Server{
			Addr:    fmt.Sprint(":", cfg.Port),
			Handler: n,
		},
		cfg,
		logger,
	}
}

// Start an HTTP server which will block the current goroutine. Will write an error to the `errCh` if a problem occurs.
func (h *Server) Start(errCh chan<- error) {
	h.logger.Info().Msg(fmt.Sprint("running server on 0.0.0.0:", h.cfg.Port))

	err := h.ListenAndServe()
	if err != http.ErrServerClosed {
		h.logger.Error().Caller().Err(err).Msg("http server stopped unexpected")
		errCh <- err
	} else {
		h.logger.Info().Msg("http server process stopped")
	}
}
