package registryManager

import (
	"errors"

	rg "golang.org/x/sys/windows/registry"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/toolbox"
)

func IsNeedUpdate(tool toolbox.ToolModel) (bool, error) {
	needUpdate := false

	checkKey := func(key rg.Key, path string, expectedValues map[string]string) error {
		currentKey, err := rg.OpenKey(key, path, rg.QUERY_VALUE)
		if err != nil {
			if errors.Is(err, rg.ErrNotExist) {
				needUpdate = true
				return nil
			}
			return err
		}
		defer func(currentKey rg.Key) {
			_ = currentKey.Close()
		}(currentKey)

		for name, expectedValue := range expectedValues {
			value, _, err := currentKey.GetStringValue(name)
			if err != nil {
				if errors.Is(err, rg.ErrNotExist) {
					needUpdate = true
					return nil
				}
				return err
			}
			if value != expectedValue {
				needUpdate = true
				return nil
			}
		}

		return nil
	}

	for _, keySpec := range keySpecs {
		toolInRegistry := createKeyValuer(tool, keySpec)

		err := checkKey(toolInRegistry.keySpec.key, toolInRegistry.mainPath, map[string]string{
			toolInRegistry.titleKey: toolInRegistry.titleValue,
			toolInRegistry.iconKey:  toolInRegistry.iconValue,
		})
		if err != nil {
			return false, err
		}
		if needUpdate {
			return true, nil
		}

		err = checkKey(toolInRegistry.keySpec.key, toolInRegistry.commandPath, map[string]string{
			toolInRegistry.commandKey: toolInRegistry.commandValue,
		})
		if err != nil {
			return false, err
		}

		if needUpdate {
			return true, nil
		}
	}

	return needUpdate, nil
}
