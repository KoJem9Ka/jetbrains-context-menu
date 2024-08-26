package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/toolbox"
	"github.com/KoJem9Ka/jetbrains-context-menu/internal/registryManager"
	"github.com/KoJem9Ka/jetbrains-context-menu/internal/shared"
)

func createCmdSync() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		changes := uint(0)
		for _, tool := range toolbox.State.Tools {
			if isHas, err := registryManager.Has(tool); err != nil || !isHas {
				if err != nil {
					color.Red("Error checking tool: %v", err)
					shared.Exit(1)
				}
				continue
			}
			if isNeedUpdate, err := registryManager.IsNeedUpdate(tool); err != nil || !isNeedUpdate {
				if err != nil {
					color.Red("Error checking tool: %v", err)
					shared.Exit(1)
				}
				continue
			}
			err := registryManager.Upsert(tool)
			if err != nil {
				color.Red("Error updating tool to registry: %v", err)
				shared.Exit(1)
			}
			path := color.HiBlackString(fmt.Sprintf("(%s)", tool.InstallLocation))
			color.Green("Updated: %s %s", tool.DisplayNameUsingConfig(), path)
			changes++
		}
		if changes == 0 {
			color.Yellow("No changes were required or applied.")
		}
	}

	cmd := &cobra.Command{
		Use:    "sync",
		Short:  "Sync all currently enabled tools at registry with config",
		PreRun: preRunElevateAndConfig,
		Run:    run,
		Args:   strictDenyCommandArgs,
	}
	return cmd
}
