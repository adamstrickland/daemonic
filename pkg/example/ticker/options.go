package ticker

import (
	"github.com/adamstrickland/daemonic/pkg/daemon"
)

type Option func(*Ticker) error

func WithLogger(logger daemon.Logger) Option {
	return func(t *Ticker) error {
		t.logger = logger
		return nil
	}
}
