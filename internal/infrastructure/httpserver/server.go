package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/michaelcardoza/dummy-smtp/internal/infrastructure/api"
	"github.com/michaelcardoza/dummy-smtp/internal/infrastructure/web"
)

type Server struct {
	log    *slog.Logger
	server *http.Server
}

type Options struct {
	Addr       string
	Log        *slog.Logger
	ApiHandler *api.Handler
	WebHandler *web.Handler
}

func New(opts Options) *Server {
	mux := http.NewServeMux()
	mux.Handle("/api/v1/", opts.ApiHandler.Routes())
	mux.Handle("/", opts.WebHandler.Routes())

	return &Server{
		log: opts.Log,
		server: &http.Server{
			Addr:         opts.Addr,
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 0,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	s.log.Info("http server listening", "addr", s.server.Addr)
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			s.log.Error("http server error", "error", err)
		}
	}()
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
