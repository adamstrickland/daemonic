package klicker

import (
	"github.com/adamstrickland/daemonic/pkg/daemon"
)

type Option func(*Klicker) error

func WithBootstrapURIs(uris []string) Option {
	return func(t *Klicker) error {
		t.bootstrapURIs = uris
		return nil
	}
}

func WithLogger(logger daemon.Logger) Option {
	return func(t *Klicker) error {
		t.logger = logger
		return nil
	}
}
