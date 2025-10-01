package tocker

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adamstrickland/daemonic/pkg/daemon"
)

type TockServer struct {
	logger daemon.Logger
	server *http.Server
	port   int
}

func NewTockServer(options ...AnyOption) (*TockServer, error) {
	t := &TockServer{
		logger: nil,
	}

	for _, opt := range options {
		opt(t)
	}

	if t.logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return t, nil
}

func (s *TockServer) Setup(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/tick", s.handleTick)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	return nil
}

func (s *TockServer) handleTick(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	timestamp := time.Now().Format(time.RFC3339)
	s.logger.Info("tick request received", "timestamp", timestamp, "remote_addr", r.RemoteAddr)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(timestamp))
}

func (s *TockServer) Run(ctx context.Context) error {
	if s.server == nil {
		return fmt.Errorf("HTTP server not initialized")
	}

	// Start server in a goroutine
	errorCh := make(chan error, 1)
	go func() {
		s.logger.Info("starting HTTP server", "addr", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errorCh <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		return nil
	case err := <-errorCh:
		return err
	}
}

func (s *TockServer) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	s.logger.Info("shutting down HTTP server")
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP server shutdown error: %w", err)
	}
	s.logger.Info("HTTP server shutdown complete")

	return nil
}
