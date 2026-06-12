package api

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
	Log         *slog.Logger
	MailService MailService
}

func NewServer(addr string, deps Deps) *Server {
	handler := NewHandler(deps.MailService, deps.Log)
	return &Server{
		log: deps.Log,
		server: &http.Server{
			Addr:         addr,
			Handler:      handler.Routes(),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	s.log.Info("http server listening", "addr", s.server.Addr)
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.log.Error("http server error", "error", err)
		}
	}()
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
