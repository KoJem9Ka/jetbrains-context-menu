package registryManager

import (
	"errors"

	rg "golang.org/x/sys/windows/registry"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/toolbox"
	"github.com/KoJem9Ka/jetbrains-context-menu/internal/shared"
)

func Remove(tool toolbox.ToolModel) error {
	errs := make([]error, 0)
	for _, keySpec := range keySpecs {
		keyValuer := createKeyValuer(tool, keySpec)
		for _, path := range keyValuer.BackPaths() {
			if err := rg.DeleteKey(keyValuer.keySpec.key, path); err != nil {
				errs = append(errs, err)
				continue
			}
			shared.LogVerbose("Deleted rg key: \"%s\"\n", path)
		}
	}
	return errors.Join(errs...)
}
