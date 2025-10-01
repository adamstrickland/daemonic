package main

import (
	"context"

	"github.com/adamstrickland/daemonic/pkg/daemon"
	"github.com/adamstrickland/daemonic/pkg/example"
)

type TickCommand struct{}

func (TickCommand) Run() error {
	logger, closer, err := getLogger(config.UseZap)
	if err != nil {
		return err
	}
	defer closer()

	logger.Info("starting ticker", "config", config)

	archon, err := daemon.NewArchon(daemon.WithLogger(logger))
	if err != nil {
		return err
	}

	app, err := example.NewTicker(example.WithLogger(logger))
	if err != nil {
		return err
	}

	if err := archon.Run(context.Background(), app); err != nil {
		return err
	}

	return nil
}
