package registryManager

import (
	"errors"

	rg "golang.org/x/sys/windows/registry"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/toolbox"
)

// Has checks if at least one rg key related to the tool exists.
// The function returns true if at least one key is successfully opened, regardless of any errors encountered.
// If any errors occur during the process, the one of the errors is returned along with the boolean result.
func Has(tool toolbox.ToolModel) (hasKey bool, firstError error) {
	for _, keySpec := range keySpecs {
		toolInRegistry := createKeyValuer(tool, keySpec)
		for _, path := range toolInRegistry.BackPaths() {
			key, err := rg.OpenKey(keySpec.key, path, rg.QUERY_VALUE)
			if err == nil {
				hasKey = true
				key.Close()
			} else if !errors.Is(err, rg.ErrNotExist) && firstError == nil {
				firstError = err
			}
		}
	}
	return hasKey, firstError
}
