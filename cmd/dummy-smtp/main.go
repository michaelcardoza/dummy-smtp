// cmd/dummy-smtp/main.go
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/michaelcardoza/dummy-smtp/internal/core/mail"
	"github.com/michaelcardoza/dummy-smtp/internal/infrastructure/api"
	"github.com/michaelcardoza/dummy-smtp/internal/infrastructure/httpserver"
	"github.com/michaelcardoza/dummy-smtp/internal/infrastructure/smtpserver"
	"github.com/michaelcardoza/dummy-smtp/internal/infrastructure/storage/memory"
	"github.com/michaelcardoza/dummy-smtp/internal/infrastructure/web"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	smtpAddr := getenv("SMTP_ADDR", ":1025")
	httpAddr := getenv("HTTP_ADDR", ":8025")

	repo := memory.New()
	broker := web.NewBroker()
	mailService := mail.NewService(repo, broker)

	smtpServer := smtpserver.NewServer(smtpAddr, smtpserver.Deps{
		Log:      log,
		Capturer: mailService,
	})

	apiHandler := api.NewHandler(mailService, log)
	webHandler, _ := web.NewHandler(broker)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", apiHandler.Routes())
	mux.Handle("/", webHandler.Routes())

	httpServer := httpserver.New(httpAddr, httpserver.Deps{
		Log:     log,
		Handler: mux,
	})

	if err := smtpServer.Start(ctx); err != nil {
		log.Error("failed to start smtp server", "error", err)
		os.Exit(1)
	}

	if err := httpServer.Start(); err != nil {
		log.Error("failed to start http server", "error", err)
		os.Exit(1)
	}

	<-ctx.Done()
	log.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("http shutdown error", "error", err)
	}

	if err := smtpServer.Shutdown(); err != nil {
		log.Error("smtp shutdown error", "error", err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
