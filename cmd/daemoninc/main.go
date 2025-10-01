package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Create root context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Initialize your service(s)
	svc := NewTickerService()

	// Start service in goroutine
	errCh := make(chan error, 1)
	go func() {
		slog.Info("starting service")
		if err := svc.Run(ctx); err != nil {
			errCh <- err
		}
	}()

	// Wait for shutdown signal or error
	select {
	case err := <-errCh:
		return fmt.Errorf("service error: %w", err)
	case sig := <-sigCh:
		slog.Info("received signal", "signal", sig)
	}

	// Graceful shutdown
	cancel() // Signal context cancellation to service

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	slog.Info("shutting down gracefully")
	if err := svc.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	slog.Info("shutdown complete")
	return nil
}

// Service represents a long-running service (HTTP server, Kafka consumer, etc.)
type Service interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

// TickerService writes timestamp every second
type TickerService struct {
	ticker *time.Ticker
	done   chan struct{}
}

func NewTickerService() *TickerService {
	return &TickerService{
		done: make(chan struct{}),
	}
}

func (s *TickerService) Run(ctx context.Context) error {
	s.ticker = time.NewTicker(1 * time.Second)
	defer s.ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case t := <-s.ticker.C:
			slog.Info("tick", "timestamp", t.Format(time.RFC3339))
		case <-s.done:
			return nil
		}
	}
}

func (s *TickerService) Shutdown(ctx context.Context) error {
	close(s.done)
	return nil
}
