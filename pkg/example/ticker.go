package example

import (
	"context"
	"fmt"
	"time"

	"github.com/adamstrickland/daemonic/pkg/daemon"
)

type Ticker struct {
	ticker *time.Ticker
	logger daemon.Logger
}

type Option func(*Ticker)

func WithLogger(logger daemon.Logger) Option {
	return func(t *Ticker) {
		t.logger = logger
	}
}

func NewTicker(options ...Option) (*Ticker, error) {
	t := &Ticker{
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

func (s *Ticker) Setup(ctx context.Context) error {
	return nil
}

func (s *Ticker) Run(ctx context.Context) error {
	defer s.ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case t := <-s.ticker.C:
			s.logger.Info("tick", "timestamp", t.Format(time.RFC3339))
		}
	}
}

func (s *Ticker) Shutdown(ctx context.Context) error {
	return nil
}
