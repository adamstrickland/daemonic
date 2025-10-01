package daemon

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Archon struct {
	logger  Logger
	timeout time.Duration
	signals []os.Signal
}

type ArchonOption func(*Archon)

func WithSlog(slogger *slog.Logger) ArchonOption {
	return func(a *Archon) {
		a.logger = NewSlogAdapter(slogger)
	}
}

func WithZap(zlogger *zap.Logger) ArchonOption {
	return func(a *Archon) {
		a.logger = NewZapAdapter(zlogger)
	}
}

func WithTimeout(d time.Duration) ArchonOption {
	return func(a *Archon) {
		a.timeout = d
	}
}

func NewArchon(options ...ArchonOption) (*Archon, error) {
	archon := &Archon{
		logger:  nil,
		timeout: 30 * time.Second,
		signals: []os.Signal{os.Interrupt, syscall.SIGTERM},
	}

	WithSlog(slog.New(slog.NewJSONHandler(os.Stdout, nil)))(archon)

	for _, opt := range options {
		opt(archon)
	}

	return archon, nil
}

func (a *Archon) Run(ctx context.Context, daemon Daemon) error {
	// Create root context with cancellation
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, a.signals...)
	defer signal.Stop(sigCh)

	// Setup the service
	if err := daemon.Setup(runCtx); err != nil {
		return fmt.Errorf("service setup failed: %w", err)
	}

	// Start service in goroutine
	errCh := make(chan error, 1)
	go func() {
		a.logger.Info("starting service")
		if err := daemon.Run(runCtx); err != nil {
			errCh <- err
		}
	}()

	// Wait for shutdown signal or error
	select {
	case err := <-errCh:
		return fmt.Errorf("service error: %w", err)
	case sig := <-sigCh:
		a.logger.Info("received signal", "signal", sig)
	}

	// Graceful shutdown
	cancel() // Signal context cancellation to service

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), a.timeout)
	defer shutdownCancel()

	a.logger.Info("shutting down gracefully", "timeout", a.timeout)
	if err := daemon.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	a.logger.Info("shutdown complete")

	return nil
}
