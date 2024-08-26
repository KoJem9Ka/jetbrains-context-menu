package commands

import (
	"strings"

	"github.com/spf13/cobra"
)

func createCmdConfigCreate() *cobra.Command {
	var isAtHomeDir bool

	example := strings.Trim(`
  cfg           write config file in the current directory    
  cfg -H        write config file in home directory
  cfg rm        remove the config file from both directories if exists`, "\n")

	run := func(cmd *cobra.Command, args []string) {
		saveConfig(isAtHomeDir)
	}

	cmd := &cobra.Command{
		Use:     "cfg",
		Short:   "Create or update config file",
		Example: example,
		PreRun:  preRunConfigIgnoringToolbox,
		Run:     run,
		Args:    strictDenyCommandArgs,
	}

	cmd.Flags().BoolVarP(&isAtHomeDir, "home", "H", false, "Write config to home directory")

	cmd.AddCommand(createCmdConfigRemove())
	cmd.AddCommand(createCmdConfigMove())

	return cmd
}
