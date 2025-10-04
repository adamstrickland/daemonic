package gateway

import "github.com/twmb/franz-go/pkg/kgo"

type Handler interface {
	Handle(*kgo.Record) ([]*kgo.Record, error)
}
