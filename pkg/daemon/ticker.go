package daemon

import (
	"context"
	"log/slog"
	"time"
)

// TickerDaemon writes timestamp every second
type Ticker struct {
	ticker *time.Ticker
	done   chan struct{}
}

func NewTicker() *Ticker {
	return &Ticker{
		done: make(chan struct{}),
	}
}

func (s *Ticker) Setup(ctx context.Context) error {
	return nil
}

func (s *Ticker) Run(ctx context.Context) error {
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

func (s *Ticker) Shutdown(ctx context.Context) error {
	close(s.done)
	return nil
}
