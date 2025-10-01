package tocker

import (
	"fmt"

	"github.com/adamstrickland/daemonic/pkg/daemon"
)

type AnyOption func(any) error

// type TockerOption func(*Tocker) error

func WithLogger(logger daemon.Logger) AnyOption {
	return func(a any) error {
		switch t := a.(type) {
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

func WithServer(port int) AnyOption {
	return func(a any) error {
		if t, ok := a.(*Tocker); ok {
			ts, _ := NewTockServer(WithLogger(t.logger),
				WithPort(t.port))
			t.tockServer = ts
			return nil
		}

		return fmt.Errorf("WithServer can only be used with Tocker type")
	}
}

func WithClient(port int) AnyOption {
	return func(a any) error {
		if t, ok := a.(*Tocker); ok {
			tc, _ := NewTockClient(WithLogger(t.logger),
				WithPort(t.port))
			t.tockClient = tc
			return nil
		}

		return fmt.Errorf("WithClient can only be used with Tocker type")
	}
}

func WithPort(port int) AnyOption {
	return func(a any) error {
		switch t := a.(type) {
		case *TockServer:
			t.port = port
			return nil
		case *TockClient:
			t.port = port
			return nil
		}
		return fmt.Errorf("unknown type for WithPort option")
	}
}
