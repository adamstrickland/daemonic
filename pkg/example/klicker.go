package example

import (
	"context"
	"fmt"
	"time"

	"github.com/adamstrickland/daemonic/pkg/daemon"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

const TopicName = "daemonic.klicker"

type Klicker struct {
	ticker        *time.Ticker
	logger        daemon.Logger
	client        *kgo.Client
	closers       []func() error
	bootstrapURIs []string
}

func WithBootstrapURIs(uris []string) AnyOption {
	return func(a any) error {
		if t, ok := a.(*Klicker); ok {
			t.bootstrapURIs = uris
			return nil
		}
		return fmt.Errorf("WithBootstrapURIs can only be used with Klicker type")
	}
}

func NewKlicker(options ...AnyOption) (*Klicker, error) {
	t := &Klicker{
		logger:        nil,
		ticker:        time.NewTicker(1 * time.Second),
		client:        nil,
		closers:       nil,
		bootstrapURIs: nil,
	}

	for _, opt := range options {
		opt(t)
	}

	if t.logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	if len(t.bootstrapURIs) == 0 {
		return nil, fmt.Errorf("at least one bootstrap URI is required")
	}

	return t, nil
}

func (s *Klicker) Setup(ctx context.Context) error {
	adm, err := kadm.NewOptClient(
		kgo.SeedBrokers(s.bootstrapURIs...),
	)
	if err != nil {
		return fmt.Errorf("failed to create kafka admin client: %w", err)
	}
	defer adm.Close()

	tds, err := adm.ListTopics(ctx, TopicName)
	if err != nil {
		return fmt.Errorf("failed to list topics: %w", err)
	}

	if _, exists := tds[TopicName]; exists {
		s.logger.Info("topic already exists, skipping creation", "topic", TopicName)
	} else {
		s.logger.Info("creating topic", "topic", TopicName)
		ctr, err := adm.CreateTopic(ctx, 1, 1, nil, TopicName)
		if err != nil || ctr.Err != nil {
			return fmt.Errorf("failed to create topic %q: %w", TopicName, err)
		}
	}

	klient, err := kgo.NewClient(
		kgo.SeedBrokers(s.bootstrapURIs...),
	)
	if err != nil {
		return fmt.Errorf("failed to create kafka client: %w", err)
	}

	s.client = klient
	s.closers = append(s.closers, func() error {
		s.logger.Info("closing kafka client")
		s.client.Close()
		s.logger.Info("kafka client closed")
		return nil
	})

	return nil
}

func (s *Klicker) Run(ctx context.Context) error {
	defer s.ticker.Stop()

	if s.client == nil {
		return fmt.Errorf("kafka client is not initialized")
	}

	if s.client.Ping(ctx) != nil {
		return fmt.Errorf("kafka client is closed")
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case t := <-s.ticker.C:
			s.onTick(ctx, t)
		}
	}
}

func (s *Klicker) onTick(ctx context.Context, t time.Time) {
	tick := t.Format(time.RFC3339)

	record := &kgo.Record{
		Topic: TopicName,
		Value: []byte(tick),
	}

	err := s.client.ProduceSync(ctx, record).FirstErr()
	if err != nil {
		s.logger.Error("klick", "timestamp", tick, "error", err)
		return
	}

	s.logger.Info("klick", "timestamp", tick)
}

func (s *Klicker) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down klicker")
	for _, closer := range s.closers {
		closer()
	}
	s.logger.Info("klicker is down")

	return nil
}
