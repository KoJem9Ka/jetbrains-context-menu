package registryManager

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/sys/windows/registry"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/config"
	"github.com/KoJem9Ka/jetbrains-context-menu/internal/interactive/regCleanView"
)

func Cleanup() {
	config.EnsureRegistryPrefixOrExit()

	getSubKeysOf := func(keySpec keySpecModel) []string {
		key, err := registry.OpenKey(keySpec.key, keySpec.path, registry.ENUMERATE_SUB_KEYS)
		if err != nil {
			color.Red("Error opening key: %v\n", err)
			return nil
		}
		defer func(key registry.Key) {
			_ = key.Close()
		}(key)

		subKeys, err := key.ReadSubKeyNames(-1)
		if err != nil {
			color.Red("Error reading subkeys: %v\n", err)
			return nil
		}

		return subKeys
	}

	getDefaultStringValueOf := func(key registry.Key, path string) string {
		openedKey, err := registry.OpenKey(key, path, registry.QUERY_VALUE)
		if err != nil {
			color.Red("Error opening key: %v\n", err)
			return ""
		}
		defer func(openedKey registry.Key) {
			_ = openedKey.Close()
		}(openedKey)

		value, _, err := openedKey.GetStringValue("")
		if err != nil {
			color.Red("Error reading value: %v\n", err)
			return ""
		}

		return value
	}

	keysToDelete := make([]regCleanView.RegRow, 0, 20)
	for _, keySpec := range keySpecs {
		for _, subKey := range getSubKeysOf(keySpec) {
			for _, registryPrefix := range []string{config.Config.RegistryUniqPrefix, config.RegistryOldPrefix} {
				if strings.HasPrefix(subKey, registryPrefix+"_") {
					mainPath := filepath.Join(keySpec.path, subKey)
					commandPath := filepath.Join(mainPath, "command")
					keysToDelete = append(keysToDelete, regCleanView.RegRow{
						Key:         keySpec.key,
						MainPath:    mainPath,
						CommandPath: commandPath,
						Title:       getDefaultStringValueOf(keySpec.key, mainPath),
						Command:     getDefaultStringValueOf(keySpec.key, commandPath),
					})
				}
			}
		}
	}

	if len(keysToDelete) == 0 {
		color.Green("Nothing to cleanup")
		return
	}

	sort.Slice(keysToDelete, func(i, j int) bool {
		return keysToDelete[i].MainPath < keysToDelete[j].MainPath
	})
	keysToDelete = regCleanView.MustSelectRegKeysToDelete(keysToDelete)
	for _, key := range keysToDelete {
		for _, path := range []string{key.CommandPath, key.MainPath} {
			err := registry.DeleteKey(key.Key, path)
			if err != nil {
				color.Red("Error deleting registry key \"%s\": %v\n", path, err)
				continue
			}
			color.Yellow("Deleted registry key \"%s\"\n", path)
		}
	}
}
