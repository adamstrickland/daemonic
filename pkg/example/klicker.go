package example

import (
	"context"
	"fmt"
	"time"

	"github.com/adamstrickland/daemonic/pkg/daemon"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Klicker struct {
	ticker        *time.Ticker
	logger        daemon.Logger
	client        *kgo.Client
	closers       []func()
	bootstrapURIs []string
}

func WithBootstrapURIs(uris []string) AnyOption {
	return func(a any) {
		if t, ok := a.(*Klicker); ok {
			t.bootstrapURIs = uris
		}
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
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(s.bootstrapURIs...),
		kgo.ConsumerGroup("klicker"),
		// kgo.ConsumeTopics("foo"),
	)
	if err != nil {
		return fmt.Errorf("failed to create kafka client: %w", err)
	}

	s.client = cl
	s.closers = append(s.closers, cl.Close)

	return nil
}

func (s *Klicker) Run(ctx context.Context) error {
	defer s.ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case t := <-s.ticker.C:
			s.logger.Info("klick", "timestamp", t.Format(time.RFC3339))
		}
	}
}

func (s *Klicker) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down klicker")
	for _, closer := range s.closers {
		closer()
	}
	s.logger.Info("klicker is down")

	return nil
}

// cl, err := kgo.NewClient(
// 	kgo.SeedBrokers(seeds...),
// 	kgo.ConsumerGroup("my-group-identifier"),
// 	kgo.ConsumeTopics("foo"),
// )
// if err != nil {
// 	panic(err)
// }
// defer cl.Close()
//
// ctx := context.Background()
//
// // 1.) Producing a message
// // All record production goes through Produce, and the callback can be used
// // to allow for synchronous or asynchronous production.
// var wg sync.WaitGroup
// wg.Add(1)
// record := &kgo.Record{Topic: "foo", Value: []byte("bar")}
// cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
// 	defer wg.Done()
// 	if err != nil {
// 		fmt.Printf("record had a produce error: %v\n", err)
// 	}
//
// })
// wg.Wait()
