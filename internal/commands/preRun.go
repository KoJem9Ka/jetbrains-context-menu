package commands

import (
	"errors"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/toolbox"
	"github.com/KoJem9Ka/jetbrains-context-menu/internal/elevate"
	"github.com/KoJem9Ka/jetbrains-context-menu/internal/shared"
)

func preRunElevateAndConfig(cmd *cobra.Command, args []string) {
	if err := elevate.ElevatedOrError(); err != nil {
		color.Red(err.Error())
		shared.Exit(1)
	}
	if err := initConfig(); err != nil {
		color.Red(err.Error())
		shared.Exit(1)
	}
}

func preRunConfigIgnoringToolbox(cmd *cobra.Command, args []string) {
	if err := initConfig(); err != nil && !errors.As(err, &toolbox.StateReadError{}) {
		color.Red("Error on initializing: %v\n", err)
		shared.Exit(1)
	}
}
