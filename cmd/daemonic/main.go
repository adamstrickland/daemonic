package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/adamstrickland/daemonic/pkg/daemon"
	"github.com/adamstrickland/daemonic/pkg/example"
	"github.com/alecthomas/kong"
	"go.uber.org/zap"
)

var config struct {
	UseZap bool `name:"zap" optional:"" help:"Use zap logger instead of slog."`
}

func main() {
	k, err := kong.New(&config, kong.DefaultEnvars("MOM"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	k.Parse(os.Args[1:])

	var logger daemon.Logger

	if config.UseZap {
		zlogger, err := zap.NewProduction()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating logger: %v\n", err)
			os.Exit(1)
		}
		defer zlogger.Sync()
		logger = daemon.NewZapAdapter(zlogger)
	} else {
		slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		logger = daemon.NewSlogAdapter(slogger)
	}

	archon, err := daemon.NewArchon(daemon.WithLogger(logger))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	app, err := example.NewTicker(example.WithLogger(logger))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := archon.Run(context.Background(), app); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
