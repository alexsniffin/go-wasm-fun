package server

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/alexsniffin/website/internal/server/models"
	"github.com/alexsniffin/website/internal/server/processes/http"
	"github.com/alexsniffin/website/internal/server/router"
)

type Server struct {
	cfg    models.Config
	logger zerolog.Logger

	httpServer *http.Server

	fatalErrCh chan error
	shutdown   sync.Once
}

func NewServer(cfg models.Config, logger zerolog.Logger) (*Server, error) {
	newRouter := router.NewRouter(cfg.HTTPRouter, logger)
	newHTTPServer := http.NewServer(cfg.HTTPServer, logger, newRouter)

	return &Server{
		cfg:        cfg,
		logger:     logger,
		httpServer: newHTTPServer,
		fatalErrCh: make(chan error),
	}, nil
}

func (s *Server) Start() {
	go s.httpServer.Start(s.fatalErrCh)

	for err := range s.fatalErrCh {
		if err != nil {
			s.logger.Error().Caller().Err(err).Msg("fatal error received from process")
			s.Shutdown(true)
		}
	}
}

func (s *Server) Shutdown(fromErr bool) {
	s.shutdown.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		s.logger.Info().Msg("attempting graceful shutdown of the server")

		graceful := make(chan bool)

		go func(graceful <-chan bool) {
			for {
				select {
				case <-ctx.Done():
					s.logger.Panic().Msg("shutdown deadline reached, terminating remaining processes ungracefully")
				case <-graceful:
					return
				}
			}
		}(graceful)

		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			s.logger.Error().Caller().Err(err).Msg("failed to shutdown http server gracefully")
		} else {
			s.logger.Info().Msg("shutdown http server gracefully")
		}

		close(s.fatalErrCh)
		close(graceful)

		if fromErr {
			s.logger.Info().Msg("graceful shutdown succeeded, an error was detected, exiting with status code 1")
			os.Exit(1)
		}
	})
}
