package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/config"
	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/toolbox"
	"github.com/KoJem9Ka/jetbrains-context-menu/internal/shared"
)

func strictDenyCommandArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		msg := "unexpected argument"
		if len(args) > 1 {
			msg += "s"
		}
		return errors.New(fmt.Sprint(msg, ": ", strings.Join(args, " ")))
	}
	return nil
}

func initConfig() (err error) {
	config.MustInit()
	if config.IsUsed() {
		err = toolbox.Init([]string{config.Config.ToolboxStateJsonPath})
	} else {
		err = toolbox.Init(config.ToolboxStateDefaultPaths)
	}
	toolsNames := map[string]string{}
	for _, tool := range toolbox.State.Tools {
		toolsNames[tool.IdStr()] = tool.DisplayName
	}
	config.SetDefaultToolsNames(toolsNames)
	return err
}

func saveConfig(isAtHomeDir bool) {
	path, isUpdated, err := config.SaveFile(isAtHomeDir)
	if err != nil {
		color.Red("Can't create config at \"%v\": %v\n", path, err)
		shared.Exit(1)
	}
	if isUpdated {
		color.Yellow("Updated config at \"%v\"\n", path)
	} else {
		color.Green("Written config at \"%v\"\n", path)
	}
}
