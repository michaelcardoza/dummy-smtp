package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/michaelcardoza/dummy-smtp/internal/config"
	"github.com/michaelcardoza/dummy-smtp/internal/core/mail"
	"github.com/michaelcardoza/dummy-smtp/internal/infrastructure/api"
	"github.com/michaelcardoza/dummy-smtp/internal/infrastructure/httpserver"
	"github.com/michaelcardoza/dummy-smtp/internal/infrastructure/smtpserver"
	"github.com/michaelcardoza/dummy-smtp/internal/infrastructure/web"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	broker := web.NewBroker()
	mailService := mail.NewService(cfg.Storage, broker)

	smtpServer := smtpserver.NewServer(smtpserver.Options{
		Addr:     cfg.STMPAddr,
		Log:      log,
		Capturer: mailService,
	})

	apiHandler := api.NewHandler(mailService)
	webHandler := web.NewHandler(broker)
	httpServer := httpserver.New(httpserver.Options{
		Addr:       cfg.HTTPAddr,
		Log:        log,
		ApiHandler: apiHandler,
		WebHandler: webHandler,
	})

	if err = smtpServer.Start(ctx); err != nil {
		log.Error("failed to start smtp server", "error", err)
		os.Exit(1)
	}

	if err = httpServer.Start(); err != nil {
		log.Error("failed to start http server", "error", err)
		os.Exit(1)
	}

	<-ctx.Done()
	log.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error("http shutdown error", "error", err)
	}

	if err = smtpServer.Shutdown(); err != nil {
		log.Error("smtp shutdown error", "error", err)
	}

	if err = cfg.Storage.Close(shutdownCtx); err != nil {
		log.Error("storage close error", "error", err)
	}
}
