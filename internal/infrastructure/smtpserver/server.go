package smtpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
)

type Server struct {
	addr     string
	capturer Capturer
	log      *slog.Logger
	listener net.Listener
}

type Deps struct {
	Log      *slog.Logger
	Capturer Capturer
}

func NewServer(addr string, deps Deps) *Server {
	return &Server{
		addr:     addr,
		capturer: deps.Capturer,
		log:      deps.Log,
	}
}

func (s *Server) Start(ctx context.Context) error {
	var lc net.ListenConfig
	listener, err := lc.Listen(ctx, "tcp", s.addr)
	if err != nil {
		return fmt.Errorf("smtp listen on %s: %w", s.addr, err)
	}
	s.listener = listener
	s.log.Info("smtp server listening", "addr", s.addr)

	go s.accept(ctx)
	return nil
}

func (s *Server) accept(ctx context.Context) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			s.log.Error("smtp accept error", err)
			continue
		}
		go newSession(conn, s.capturer).serve(ctx)
	}
}

func (s *Server) Shutdown() error {
	if s.listener == nil {
		return nil
	}
	return s.listener.Close()
}
