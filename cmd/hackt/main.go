package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/tsatke/hackt"
)

var (
	root = &cobra.Command{
		Use:     "hackt",
		Example: "hackt myproject",
		Args:    cobra.ExactArgs(1),
		Version: "0.0.1",
		RunE:    runApplication,
	}
)

func main() {
	if err := root.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func runApplication(cmd *cobra.Command, args []string) error {
	logFile, err := os.OpenFile("hackt.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("unable to open log file: %w", err)
	}

	cw := zerolog.ConsoleWriter{
		Out:     logFile,
		NoColor: false,
	}
	log := zerolog.New(cw).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Logger()

	projectName := args[0]

	app := hackt.NewApplication(log)
	app.Run()

	if err := app.LoadProjectFromPath(projectName); err != nil {
		return fmt.Errorf("load project: %w", err)
	}

	err = app.Wait()
	log.Info().Err(err).Msg("shutting down...")
	return err
}
