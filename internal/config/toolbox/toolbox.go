package toolbox

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/shared"
)

func Init(paths []string) error {
	path, err := getStatePath(paths)
	if err != nil {
		return StateReadError{fmt.Errorf("can't identify JetBrains tools: %w", err)}
	}
	State, err = readToolboxState(path)
	if err != nil {
		return StateReadError{fmt.Errorf("can't read JetBrains Toolbox state.json: %w", err)}
	}
	shared.LogVerbose("Using JetBrains Toolbox v%s at: \"%s\"", State.AppVersion, path)
	return nil
}

func getStatePath(paths []string) (path string, err error) {
	for _, path = range paths {
		if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
			return path, nil
		}
	}
	var _paths any = paths
	if len(paths) == 1 {
		_paths = paths[0]
	}
	return "", fmt.Errorf(`JetBrains Toolbox state.json not found in "%s"`, _paths)
}

func readToolboxState(path string) (state state, err error) {
	file, err := os.Open(path)
	if err != nil {
		return state, fmt.Errorf("failed to open: %w", err)
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&state); err != nil {
		return state, fmt.Errorf("failed to decode: %w", err)
	}
	return state, nil
}
