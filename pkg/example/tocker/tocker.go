package tocker

import (
	"context"
	"fmt"

	"github.com/adamstrickland/daemonic/pkg/daemon"
)

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

	return t, nil
}

func (s *Tocker) Setup(ctx context.Context) error {
	if s.tockServer != nil {
		s.logger.Info("setting up tock server")
		err := s.tockServer.Setup(ctx)
		if err != nil {
			return fmt.Errorf("tock server setup failed: %w", err)
		}
	}

	if s.tockClient != nil {
		s.logger.Info("setting up tock client")
		err := s.tockClient.Setup(ctx)
		if err != nil {
			return fmt.Errorf("tock client setup failed: %w", err)
		}
	}

	return nil
}

func (s *Tocker) Run(ctx context.Context) error {
	// Run both the server and client concurrently
	errorCh := make(chan error, 2)

	// Start TockServer
	if s.tockServer != nil {
		go func() {
			if err := s.tockServer.Run(ctx); err != nil {
				errorCh <- fmt.Errorf("tock server error: %w", err)
			}
		}()
	}

	// Start TockClient
	if s.tockClient != nil {
		go func() {
			if err := s.tockClient.Run(ctx); err != nil {
				errorCh <- fmt.Errorf("tock client error: %w", err)
			}
		}()
	}

	// Wait for context cancellation or any service error
	select {
	case <-ctx.Done():
		return nil
	case err := <-errorCh:
		return err
	}
}

func (s *Tocker) Shutdown(ctx context.Context) error {
	if s.tockServer != nil {
		s.logger.Info("shutting down tock server")
		err := s.tockServer.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("tock server shutdown failed: %w", err)
		}
	}

	if s.tockClient != nil {
		s.logger.Info("shutting down tock client")
		err := s.tockClient.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("tock client shutdown failed: %w", err)
		}
	}

	return nil
}
