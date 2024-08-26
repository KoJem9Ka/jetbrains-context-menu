package commands

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/config"
)

func createCmdConfigRemove() *cobra.Command {
	removeFileHandler := func(path string, err error) {
		if err == nil {
			color.Green(`Config file removed at "%s"`, path)
		} else {
			color.Red(`Error removing config file at "%s": %s`, path, err)
		}
	}

	run := func(cmd *cobra.Command, args []string) {
		filesTriedToRemoveCount := config.RemoveFile(removeFileHandler)
		if filesTriedToRemoveCount == 0 {
			color.Yellow("No config files found to remove")
		}
	}

	cmd := &cobra.Command{
		Use:     "rm",
		Short:   "Remove config files",
		Long:    "Remove the config file from the current directory and home directory, if it exists.",
		Example: "cfg rm",
		Run:     run,
		Args:    strictDenyCommandArgs,
	}
	return cmd
}
