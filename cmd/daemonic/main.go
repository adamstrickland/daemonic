package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

var config struct {
	// commands
	Klick Klick `cmd:"" help:"Run the Klicker application."`
	Tick  Tick  `cmd:"" help:"Run the Ticker application."`
	Tock  Tock  `cmd:"" help:"Run the Tocker application."`

	// top-level options
	UseZap bool `name:"zap" optional:"" help:"Use zap logger instead of slog."`
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
