package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/adamstrickland/daemonic/pkg/daemon"
	"go.uber.org/zap"
)

func getLogger(useZap bool) (daemon.Logger, func(), error) {
	if useZap {
		zlogger, err := zap.NewProduction()
		if err != nil {
			return nil, nil, fmt.Errorf("error creating zap logger: %w", err)
		}

		cleanup := func() {
			zlogger.Sync()
		}

		return daemon.NewZapAdapter(zlogger), cleanup, nil
	} else {
		slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		return daemon.NewSlogAdapter(slogger), func() {}, nil
	}
}
