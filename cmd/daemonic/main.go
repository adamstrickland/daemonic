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

type KlickCommand struct {
	BrokerURIs  []string `name:"broker-uris" help:"List of Kafka broker URIs."`
	RegistryURI string   `name:"registry-uri" help:"URI for the schema registry."`
}

type TickCommand struct{}

var config struct {
	UseZap bool         `name:"zap" optional:"" help:"Use zap logger instead of slog."`
	Klick  KlickCommand `cmd:"" help:"Run the Klicker application."`
	Tick   TickCommand  `cmd:"" help:"Run the Ticker application."`
}

func main() {
	k, err := kong.New(&config, kong.DefaultEnvars("MOM"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cmd, err := k.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	err = cmd.Run()
	cmd.FatalIfErrorf(err)
}

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

func (c KlickCommand) Run() error {
	logger, closer, err := getLogger(config.UseZap)
	if err != nil {
		return err
	}
	defer closer()

	logger.Info("starting klicker", "config", config)

	archon, err := daemon.NewArchon(daemon.WithLogger(logger))
	if err != nil {
		return err
	}

	app, err := example.NewKlicker(example.WithLogger(logger),
		example.WithBootstrapURIs(c.BrokerURIs))
	if err != nil {
		return err
	}

	if err := archon.Run(context.Background(), app); err != nil {
		return err
	}

	return nil
}

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
