package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/adamstrickland/daemonic/pkg/daemon"
	"github.com/adamstrickland/daemonic/pkg/example"
	"go.uber.org/zap"
)

var (
	useZap = true
)

func init() {
	flag.BoolVar(&useZap, "zap", false, "Use zap logger (default: slog)")

	flag.Parse()
}

func main() {
	var logger daemon.Logger

	if useZap {
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
	// app, err := example.NewTicker()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := archon.Run(context.Background(), app); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
