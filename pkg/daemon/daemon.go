package daemon

import "context"

// Service represents a long-running service (HTTP server, Kafka consumer, etc.)
type Daemon interface {
	Setup(ctx context.Context) error
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
