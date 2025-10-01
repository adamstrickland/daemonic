package main

import (
	"context"

	"github.com/adamstrickland/daemonic/pkg/daemon"
	"github.com/adamstrickland/daemonic/pkg/example"
)

type Klick struct {
	BrokerURIs  []string `name:"broker-uris" help:"List of Kafka broker URIs."`
	RegistryURI string   `name:"registry-uri" help:"URI for the schema registry."`
}

func (c Klick) Run() error {
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
