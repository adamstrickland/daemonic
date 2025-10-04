package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/adamstrickland/daemonic/pkg/daemon"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Logger interface {
	daemon.Logger
}

type Gateway struct {
	brokerURIs []string
	topic      string
	client     *kgo.Client
	closers    []func() error
	logger     Logger
	name       string
	handler    Handler
}

func NewGateway(options ...Option) (*Gateway, error) {
	gw := &Gateway{
		brokerURIs: nil,
		topic:      "",
		client:     nil,
		closers:    nil,
		name:       "",
	}

	for _, opt := range options {
		if err := opt(gw); err != nil {
			return nil, err
		}
	}

	if gw.name == "" {
		WithName(fmt.Sprintf("daemonic.gateway.%p", gw))(gw)
	}

	if gw.topic == "" {
		return nil, fmt.Errorf("topic is not configured")
	}

	return gw, nil
}

func (s *Gateway) Setup(ctx context.Context) error {
	if s.client == nil {
		klient, err := kgo.NewClient(
			kgo.SeedBrokers(s.brokerURIs...),
			kgo.ConsumerGroup(s.name),
			kgo.ConsumeTopics(s.topic),
			kgo.FetchIsolationLevel(kgo.ReadCommitted),
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
	}

	return nil
}

func (s *Gateway) Run(ctx context.Context) error {
	errorCount := 0
	maxErrorCount := 5

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			// NOTE: Add backoff here.
			err := s.handle(ctx)
			if err == nil {
				errorCount = 0
			} else {
				s.logger.Error("handling errors", "errors", err)
				errorCount++
				if errorCount >= maxErrorCount {
					return fmt.Errorf("exceeded maximum error count of %d", maxErrorCount)
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (s *Gateway) handle(ctx context.Context) error {
	fetches := s.client.PollFetches(ctx)
	if errs := fetches.Errors(); len(errs) > 0 {
		return fmt.Errorf("fetch errors: %v", errs)
	}

	var err error
	fetches.EachPartition(func(p kgo.FetchTopicPartition) {
		p.EachRecord(func(record *kgo.Record) {
			err := s.client.BeginTransaction(ctx)
			if err != nil {
				s.logger.Warn("beginning transaction", "topic", record.Topic, "partition", record.Partition, "offset", record.Offset, "error", err)
				_ = s.client.AbortTransaction(ctx)
			}

			records, err := s.handler.Handle(ctx, record)
			if err != nil {
				s.logger.Warn("handling record", "topic", record.Topic, "partition", record.Partition, "offset", record.Offset, "error", err)
				_ = s.client.AbortTransaction(ctx)
			}

			if len(records) > 0 {
				err = s.client.ProduceSync(ctx, records...).FirstErr()
				if err != nil {
					s.logger.Warn("producing records", "topic", record.Topic, "partition", record.Partition, "offset", record.Offset, "error", err)
					_ = s.client.AbortTransaction(ctx)
				}
			}

			err := s.client.CommitTransaction(ctx)
			if err != nil {
				s.logger.Warn("committing transaction", "topic", record.Topic, "partition", record.Partition, "offset", record.Offset, "error", err)
				_ = s.client.AbortTransaction(ctx)
			}
		})
	})

	return nil
}

func (s *Gateway) Shutdown(ctx context.Context) error {
	return nil
}
