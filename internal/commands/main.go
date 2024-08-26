package commands

import (
	"github.com/fatih/color"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/config"
)

func Execute() {
	config.EnsureRegistryPrefixOrExit()
	if err := createCmdRoot().Execute(); err != nil {
		color.Red(err.Error())
	}
}
