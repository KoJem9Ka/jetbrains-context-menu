package commands

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/config"
)

func createCmdConfigMove() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		movedToPath, err := config.MoveFile()
		if err != nil {
			color.Red(err.Error())
			return
		}
		color.Green("Config moved to: %s", movedToPath)
	}

	cmd := &cobra.Command{
		Use:     "move",
		Short:   "Move config from current to home dir or vice versa",
		Example: "cfg move",
		Run:     run,
		Args:    strictDenyCommandArgs,
	}

	return cmd
}
