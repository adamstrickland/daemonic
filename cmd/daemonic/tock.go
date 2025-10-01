package main

import (
	"context"

	"github.com/adamstrickland/daemonic/pkg/daemon"
	"github.com/adamstrickland/daemonic/pkg/example"
)

type Tock struct {
	Port      int  `default:"8081" help:"Port to listen on"`
	RunServer bool `name:"server" negatable:"" default:"true" help:"Run the server"`
	RunClient bool `name:"client" default:"false" help:"Run the client"`
}

func (Tock) Run() error {
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

	opts := []example.AnyOption{
		example.WithLogger(logger),
	}
	if config.Tock.RunServer {
		opts = append(opts, example.WithServer(config.Tock.Port))
	}
	if config.Tock.RunClient {
		opts = append(opts, example.WithClient(config.Tock.Port))
	}

	app, err := example.NewTocker(opts...)
	if err != nil {
		return err
	}

	if err := archon.Run(context.Background(), app); err != nil {
		return err
	}

	return nil
}
