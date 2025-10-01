package example

import (
	"fmt"

	"github.com/adamstrickland/daemonic/pkg/daemon"
)

type AnyOption func(any) error

func WithLogger(logger daemon.Logger) AnyOption {
	return func(a any) error {
		switch t := a.(type) {
		case *Ticker:
			t.logger = logger
			return nil
		case *Klicker:
			t.logger = logger
			return nil
		case *Tocker:
			t.logger = logger
			return nil
		case *TockServer:
			t.logger = logger
			return nil
		case *TockClient:
			t.logger = logger
			return nil
		}
		return fmt.Errorf("unknown type for WithLogger option")
	}
}
