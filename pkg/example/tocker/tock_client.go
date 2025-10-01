package tocker

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/adamstrickland/daemonic/pkg/daemon"
)

type TockClient struct {
	ticker     *time.Ticker
	logger     daemon.Logger
	port       int
	httpClient *http.Client
}

func NewTockClient(options ...AnyOption) (*TockClient, error) {
	t := &TockClient{
		logger:     nil,
		ticker:     time.NewTicker(1 * time.Second),
		httpClient: &http.Client{Timeout: 5 * time.Second},
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

func (s *TockClient) onTick(ctx context.Context, t time.Time) {
	if s.port == 0 {
		s.logger.Info("skipping tick request, no port configured")
		return
	}

	url := fmt.Sprintf("http://localhost:%d/tick", s.port)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		s.logger.Error("failed to create request", "error", err, "url", url)
		return
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error("failed to make tick request", "error", err, "url", url)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("failed to read response body", "error", err, "url", url)
		return
	}

	s.logger.Info("tick request completed", "url", url, "status", resp.Status, "response", string(body))
}

func (s *TockClient) Shutdown(ctx context.Context) error {
	return nil
}
