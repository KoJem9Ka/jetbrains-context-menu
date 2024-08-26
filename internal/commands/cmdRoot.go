package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/configShared"
	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/toolbox"
	"github.com/KoJem9Ka/jetbrains-context-menu/internal/interactive/toolsSelectView"
	"github.com/KoJem9Ka/jetbrains-context-menu/internal/registryManager"
	"github.com/KoJem9Ka/jetbrains-context-menu/internal/shared"
)

func createCmdRoot() *cobra.Command {
	run := func(cmd *cobra.Command, args []string) {
		toolsSelected, isCancelled, isConfigNeedBeSaved := toolsSelectView.MustSelectTools(toolbox.State.Tools)
		if isCancelled {
			color.Yellow("Cancelled")
			return
		}
		if isConfigNeedBeSaved {
			saveConfig(false)
		}

		var appliedCount uint = 0
		for _, tool := range toolbox.State.Tools {
			isEnabledPrev, _ := registryManager.Has(tool)
			isEnabledNext := toolsSelected.Has(tool.Id())

			if !isEnabledPrev && !isEnabledNext {
				continue
			}

			if isEnabledPrev && !isEnabledNext {
				err := registryManager.Remove(tool)
				if err != nil {
					color.Red("Error disabling %v: %v", tool.DisplayNameUsingConfig(), err)
				} else {
					path := color.HiBlackString(fmt.Sprintf("(%s)", tool.InstallLocation))
					color.Yellow("Disable %s %s", tool.DisplayNameUsingConfig(), path)
				}
			} else if !isEnabledPrev && isEnabledNext {
				err := registryManager.Upsert(tool)
				if err != nil {
					color.Red("Error enabling %v: %v", tool.DisplayNameUsingConfig(), err)
				} else {
					path := color.HiBlackString(fmt.Sprintf("(%s)", tool.InstallLocation))
					color.Green("Enabled %s %s", tool.DisplayNameUsingConfig(), path)
				}
			} else if isNeedToolUpdate, err := registryManager.IsNeedUpdate(tool); isNeedToolUpdate {
				if err != nil {
					color.Red("Error checking tool: %v", err)
					continue
				}
				err := registryManager.Upsert(tool)
				if err != nil {
					color.Red("Error updating %v: %v", tool.DisplayNameUsingConfig(), err)
				} else {
					path := color.HiBlackString(fmt.Sprintf("(%s)", tool.InstallLocation))
					color.Green("Updated %s %s", tool.DisplayNameUsingConfig(), path)
				}
			} else {
				continue
			}
			appliedCount++
		}

		if appliedCount == 0 {
			color.Yellow("No changes applied")
		}
	}

	cmd := &cobra.Command{
		Use:    "jcm",
		Short:  "JCM is a tool to manage JetBrains Tools in Windows context menu",
		Long:   "JetBrains Context Menu is a tool to manage JetBrains Tools in Windows context menu.\nRun it directly to start!",
		PreRun: preRunElevateAndConfig,
		Run:    run,
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			shared.Exit(0)
		},
		Args: strictDenyCommandArgs,
	}

	cmd.AddCommand(createCmdConfigCreate())
	cmd.AddCommand(createCmdClean())
	cmd.AddCommand(createCmdSync())

	cmd.PersistentFlags().BoolVarP(&configShared.Verbose, "verbose", "v", false, "verbose output")

	cmd.CompletionOptions.DisableDefaultCmd = true // Disable default completion command
	cobra.MousetrapHelpText = ""                   // Disable mousetrap in Windows explorer

	return cmd
}
