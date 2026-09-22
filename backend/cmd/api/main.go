package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cassianobraz/SistemaControleChamados/backend/docs"
	httpapi "github.com/cassianobraz/SistemaControleChamados/backend/internal/delivery/http"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/infrastructure/config"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/infrastructure/database"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/infrastructure/repository"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/infrastructure/seed"
	"github.com/cassianobraz/SistemaControleChamados/backend/internal/usecase"
)

func main() {
	if err := run(); err != nil {
		slog.Error("aplicação encerrada com erro", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client, err := database.Connect(ctx, cfg.MongoURI)
	if err != nil {
		return err
	}
	defer func() {
		disconnectCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Disconnect(disconnectCtx); err != nil {
			slog.Error("falha ao desconectar do mongodb", "error", err)
		}
	}()

	db := client.Database(cfg.MongoDatabase)

	if err := database.MigrateLegacyObjectIDs(ctx, db, "tickets", "responsibles"); err != nil {
		return err
	}

	ticketRepo := repository.NewTicketRepository(db)
	responsibleRepo := repository.NewResponsibleRepository(db)

	if err := seed.Responsibles(ctx, responsibleRepo); err != nil {
		return err
	}

	router := httpapi.NewRouter(httpapi.Dependencies{
		TicketUseCase:      usecase.NewTicketUseCase(ticketRepo, responsibleRepo),
		ResponsibleUseCase: usecase.NewResponsibleUseCase(responsibleRepo, ticketRepo),
		OpenAPISpec:        docs.OpenAPISpec,
	})

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("servidor iniciado", "port", cfg.Port)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}

		serverErr <- nil
	}()

	select {
	case <-ctx.Done():
		slog.Info("sinal de encerramento recebido, desligando o servidor")
	case err := <-serverErr:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return server.Shutdown(shutdownCtx)
}
