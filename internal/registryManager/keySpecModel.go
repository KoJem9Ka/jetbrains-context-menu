package registryManager

import (
	"fmt"
	"path/filepath"

	"golang.org/x/sys/windows/registry"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/config"
)

var keySpecs = [...]keySpecModel{
	{registry.CLASSES_ROOT, "HKEY_CLASSES_ROOT", `*\shell`, &config.Config.Title.OpenFileIn},
	{registry.CLASSES_ROOT, "HKEY_CLASSES_ROOT", `Directory\shell`, &config.Config.Title.OpenDirIn},
	{registry.CLASSES_ROOT, "HKEY_CLASSES_ROOT", `Directory\background\shell`, &config.Config.Title.OpenBackgroundDirIn},
}

type keySpecModel struct {
	key       registry.Key
	keyString string
	path      string
	title     *string
}

func (this keySpecModel) MainPath(toolId string) string {
	return filepath.Join(this.path, fmt.Sprintf("%s_%s", config.Config.RegistryUniqPrefix, toolId))
}

func (this keySpecModel) CommandPath(toolId string) string {
	return filepath.Join(this.MainPath(toolId), "command")
}
