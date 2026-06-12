package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	log    *slog.Logger
	server *http.Server
}

type Deps struct {
	Log     *slog.Logger
	Handler http.Handler
}

func New(addr string, deps Deps) *Server {
	return &Server{
		log: deps.Log,
		server: &http.Server{
			Addr:         addr,
			Handler:      deps.Handler,
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
