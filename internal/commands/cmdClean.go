package commands

import (
	"github.com/spf13/cobra"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/registryManager"
)

func createCmdClean() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		registryManager.Cleanup()
	}

	cmd := &cobra.Command{
		Use:     "clean",
		Short:   "Clean the registry changes made by JCM. You will be prompted to confirm the changes.",
		Aliases: []string{"cleanup", "clear"},
		PreRun:  preRunConfigIgnoringToolbox,
		Run:     run,
		Args:    strictDenyCommandArgs,
	}
	return cmd
}
