package registryManager

import (
	rg "golang.org/x/sys/windows/registry"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/toolbox"
	"github.com/KoJem9Ka/jetbrains-context-menu/internal/shared"
)

func Upsert(tool toolbox.ToolModel) error {

	performIteration := func(toolInRegistry keyValuerModel) (err error) {
		currentRegKey, _, err := rg.CreateKey(toolInRegistry.keySpec.key, toolInRegistry.mainPath, rg.SET_VALUE)
		if err != nil {
			return err
		}
		defer func(regKey rg.Key) {
			_ = regKey.Close()
			if err != nil {
				_ = rg.DeleteKey(toolInRegistry.keySpec.key, toolInRegistry.mainPath)
				shared.LogVerbose("Rollback due error: deleted rg key: \"%s\"", toolInRegistry.mainPath)
			}
		}(currentRegKey)
		shared.LogVerbose("Created rg key: \"%s\"\n", toolInRegistry.mainPath)
		if err = currentRegKey.SetStringValue(toolInRegistry.titleKey, toolInRegistry.titleValue); err != nil {
			return err
		}
		shared.LogVerbose("Set rg value: \"%s\" \"%s\"=\"%s\"\n", toolInRegistry.mainPath, toolInRegistry.titleKey, toolInRegistry.titleValue)
		if err = currentRegKey.SetExpandStringValue(toolInRegistry.iconKey, toolInRegistry.iconValue); err != nil {
			return err
		}
		shared.LogVerbose("Set rg value: \"%s\" \"%s\"=\"%s\"\n", toolInRegistry.mainPath, toolInRegistry.iconKey, toolInRegistry.iconValue)
		currentRegKey, _, err = rg.CreateKey(toolInRegistry.keySpec.key, toolInRegistry.commandPath, rg.SET_VALUE)
		if err != nil {
			return err
		}
		defer func(regKey rg.Key) {
			_ = regKey.Close()
			if err != nil {
				_ = rg.DeleteKey(toolInRegistry.keySpec.key, toolInRegistry.commandPath)
				shared.LogVerbose("Rollback due error: deleted rg key: \"%s\"", toolInRegistry.commandPath)
			}
		}(currentRegKey)
		shared.LogVerbose("Created rg key: \"%s\"\n", toolInRegistry.commandPath)
		if err = currentRegKey.SetStringValue(toolInRegistry.commandKey, toolInRegistry.commandValue); err != nil {
			return err
		}
		shared.LogVerbose("Set rg value: \"%s\" \"%s\"=\"%s\"\n", toolInRegistry.commandPath, toolInRegistry.commandKey, toolInRegistry.commandValue)
		return nil
	}

	for _, keySpec := range keySpecs {
		toolWrapper := createKeyValuer(tool, keySpec)
		err := performIteration(toolWrapper)
		if err != nil {
			return err
		}
	}

	return nil
}
