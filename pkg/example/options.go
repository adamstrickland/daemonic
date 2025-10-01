package example

import "github.com/adamstrickland/daemonic/pkg/daemon"

type AnyOption func(any)

func WithLogger(logger daemon.Logger) AnyOption {
	return func(a any) {
		switch t := a.(type) {
		case *Ticker:
			t.logger = logger
		case *Klicker:
			t.logger = logger
		}
	}
}
