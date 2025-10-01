package example

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adamstrickland/daemonic/pkg/daemon"
)

func WithServer(port int) AnyOption {
	return func(a any) error {
		if t, ok := a.(*Tocker); ok {
			ts, _ := NewTockServer(WithLogger(t.logger),
				WithPort(t.port))
			t.tockServer = ts
			return nil
		}

		return fmt.Errorf("WithServer can only be used with Tocker type")
	}
}

func WithClient(port int) AnyOption {
	return func(a any) error {
		if t, ok := a.(*Tocker); ok {
			tc, _ := NewTockClient(WithLogger(t.logger),
				WithPort(t.port))
			t.tockClient = tc
			return nil
		}

		return fmt.Errorf("WithClient can only be used with Tocker type")
	}
}

func WithPort(port int) AnyOption {
	return func(a any) error {
		switch t := a.(type) {
		case *TockServer:
			t.port = port
			return nil
		case *TockClient:
			t.port = port
			return nil
		}
		return fmt.Errorf("unknown type for WithPort option")
	}
}

type Tocker struct {
	tockServer *TockServer
	tockClient *TockClient
	logger     daemon.Logger
	port       int
}

func NewTocker(options ...AnyOption) (*Tocker, error) {
	t := &Tocker{
		tockServer: nil,
		tockClient: nil,
		logger:     nil,
		port:       8080,
	}

	for _, opt := range options {
		opt(t)
	}

	if t.logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	if t.tockClient == nil {
		tc, err := NewTockClient(WithLogger(t.logger),
			WithPort(t.port))
		if err != nil {
			return nil, fmt.Errorf("failed to create tock client: %w", err)
		}
		t.tockClient = tc
	}

	return t, nil
}

func (s *Tocker) Setup(ctx context.Context) error {
	err := s.tockServer.Setup(ctx)
	if err != nil {
		return fmt.Errorf("tock server setup failed: %w", err)
	}

	err = s.tockClient.Setup(ctx)
	if err != nil {
		return fmt.Errorf("tock client setup failed: %w", err)
	}

	return nil
}

func (s *Tocker) Run(ctx context.Context) error {
	// Run both the server and client concurrently
	errorCh := make(chan error, 2)

	// Start TockServer
	go func() {
		if err := s.tockServer.Run(ctx); err != nil {
			errorCh <- fmt.Errorf("tock server error: %w", err)
		}
	}()

	// Start TockClient
	go func() {
		if err := s.tockClient.Run(ctx); err != nil {
			errorCh <- fmt.Errorf("tock client error: %w", err)
		}
	}()

	// Wait for context cancellation or any service error
	select {
	case <-ctx.Done():
		return nil
	case err := <-errorCh:
		return err
	}
}

func (s *Tocker) Shutdown(ctx context.Context) error {
	err := s.tockServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("tock server shutdown failed: %w", err)
	}

	err = s.tockClient.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("tock client shutdown failed: %w", err)
	}

	return nil
}

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

type TockClient struct {
	ticker *time.Ticker
	logger daemon.Logger
	port   int
}

func NewTockClient(options ...AnyOption) (*TockClient, error) {
	t := &TockClient{
		logger: nil,
		ticker: time.NewTicker(1 * time.Second),
	}

	for _, opt := range options {
		opt(t)
	}

	if t.logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return t, nil
}

func (s *TockClient) Setup(ctx context.Context) error {
	return nil
}

func (s *TockClient) Run(ctx context.Context) error {
	defer s.ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case t := <-s.ticker.C:
			s.onTick(ctx, t)
		}
	}
}

func (s *TockClient) onTick(_ context.Context, t time.Time) {
	// s.logger.Info("tick", "timestamp", t.Format(time.RFC3339))
}

func (s *TockClient) Shutdown(ctx context.Context) error {
	return nil
}
