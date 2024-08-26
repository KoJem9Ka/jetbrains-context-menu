package registryManager

import (
	"fmt"
	"strings"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/toolbox"
)

type keyValuerModel struct {
	tool    toolbox.ToolModel
	keySpec keySpecModel

	mainPath    string
	commandPath string

	titleKey   string
	titleValue string

	iconKey   string
	iconValue string

	commandKey   string
	commandValue string
}

func createKeyValuer(tool toolbox.ToolModel, keySpec keySpecModel) keyValuerModel {
	toolIdHash := tool.IdHash()
	mainPath := keySpec.MainPath(toolIdHash)

	commandArg := "1"
	if strings.Contains(keySpec.path, "background") {
		commandArg = "V"
	}

	return keyValuerModel{
		tool:    tool,
		keySpec: keySpec,

		mainPath:    mainPath,
		commandPath: keySpec.CommandPath(toolIdHash),

		titleKey:   "",
		titleValue: strings.ReplaceAll(*keySpec.title, "{}", tool.DisplayNameUsingConfig()),

		iconKey:   "Icon",
		iconValue: fmt.Sprintf(`"%s",0`, tool.ExecutablePath()),

		commandKey:   "",
		commandValue: fmt.Sprintf(`"%s" "%%%s"`, tool.ExecutablePath(), commandArg),
	}
}

func (this keyValuerModel) BackPaths() []string {
	return []string{this.commandPath, this.mainPath}
}
