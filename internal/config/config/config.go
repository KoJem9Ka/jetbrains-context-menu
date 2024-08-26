package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"unicode/utf8"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/shared"
)

const configName = "jcm.yml"

var (
	isUsed        = false
	homeDir       = os.Getenv("USERPROFILE")
	currentDir, _ = filepath.Abs(".")
	configAtHome  = filepath.Join(homeDir, configName)
	configAtHere  = filepath.Join(currentDir, configName)
)

var ToolboxStateDefaultPaths = []string{
	filepath.Join(os.Getenv("LOCALAPPDATA"), "JetBrains", "Toolbox", "state.json"),
	filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "JetBrains", "Toolbox", "state.json"),
}

// MustInit tries reads the config file from current or then from home directory.
// If the file is found, it reads the file and assigns the values to the Config variable.
// If the file is not found, it is doing nothing.
func MustInit() {
	var (
		configAt          string
		configFileContent []byte
		err               error
	)
	for _, _configAt := range [...]string{configAtHere, configAtHome} {
		if configFileContent, err = os.ReadFile(_configAt); err == nil {
			configAt = _configAt
			break
		}
	}
	if err != nil || configFileContent == nil {
		return
	}
	err = yaml.Unmarshal(configFileContent, &Config)
	if err != nil && err.Error() != "EOF" {
		color.Red("Error reading config file: %s", err)
		shared.Exit(1)
	}
	shared.LogVerbose("Using config file at: \"%s\"", configAt)
	isUsed = true
	EnsureRegistryPrefixOrExit()
}

// SetDefaultToolsNames sets the default tool names for the tools that are not in the config file.
func SetDefaultToolsNames(defaultToolsNames map[string]string) {
	for id, defaultToolName := range defaultToolsNames {
		if toolName, isExist := Config.ToolsTitles[id]; !isExist || toolName == "" {
			Config.ToolsTitles[id] = defaultToolName
		}
	}
}

// SaveFile saves the config to the file system.
// isAtHomeDir: if true, save the config to the home directory, otherwise save it to the current directory.
// If config already exists, isAtHomeDir will be ignored and the file will be updated.
// Always returns path, and error if occurs with that path.
func SaveFile(isAtHomeDir bool) (path string, isUpdated bool, err error) {
	// TODO: Fix stupidity with isAtHomeDir flag
	if _, err := os.Stat(configAtHome); err == nil {
		isAtHomeDir = true
	}
	if _, err := os.Stat(configAtHere); err == nil {
		isAtHomeDir = false
	}
	if isAtHomeDir {
		path = configAtHome
	} else {
		path = configAtHere
	}
	if _, err := os.Stat(path); err == nil {
		isUpdated = true
	}
	file, err := os.Create(path)
	if err != nil {
		return path, false, err
	}
	defer file.Close()

	node := yaml.Node{}
	err = node.Encode(Config)
	if err != nil {
		return path, false, err
	}

	node.HeadComment = `########################################################
# Config for JetBrains Context Menu (JCM) by @KoJem9Ka #
########################################################`
	shared.SetYamlComment(&node, []string{"registry-uniq-prefix"}, "Unique prefix for made registry keys. Better not touch it, because JCM can't sync it, ...but min length is 10 characters.")
	shared.SetYamlComment(&node, []string{"title"}, "Titles for context menu items")
	shared.SetYamlComment(&node, []string{"title", "open-file-in"}, "Title when right-clicking on a file")
	shared.SetYamlComment(&node, []string{"title", "open-dir-in"}, "Title when right-clicking on a folder")
	shared.SetYamlComment(&node, []string{"title", "open-background-dir-in"}, "Title when right-clicking on folder background")
	shared.SetYamlComment(&node, []string{"toolbox-state-json-path"}, "Path to the state.json file of the JetBrains Toolbox")
	shared.SetYamlComment(&node, []string{"tools-search-paths"}, "Paths for additional search of tools (long recursive file search) (not implemented)")
	shared.SetYamlComment(&node, []string{"tools-titles"}, "Custom titles for tools")
	shared.SetYamlComment(&node, []string{"custom-tools"}, `Custom tools (C:\MyProg\bin.exe: Start &{} here) (not implemented)`)

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	err = encoder.Encode(node)
	if err != nil {
		return path, false, err
	}
	return path, isUpdated, nil
}

// RemoveFile removes the config file from the current directory and home directory, if it exists.
// handler: a function that will be called for each file that is tried to remove.
// returns the number of files tried to remove.
func RemoveFile(handler func(path string, err error)) uint {
	var filesTriedToRemoveCount uint = 0
	for _, path := range [...]string{configAtHome, configAtHere} {
		if _, err := os.Stat(path); err != nil {
			continue
		}
		filesTriedToRemoveCount++
		err := os.Remove(path)
		handler(path, err)
	}
	return filesTriedToRemoveCount
}

func MoveFile() (movedToPath string, err error) {
	configAt, configTo := configAtHere, configAtHome
	if _, err = os.Stat(configAt); err != nil {
		configAt, configTo = configTo, configAtHere
	}
	if _, err = os.Stat(configAt); err != nil {
		return "", errors.New("config file not found")
	}

	input, err := os.ReadFile(configAt)
	if err != nil {
		return "", fmt.Errorf("error reading config file: %w", err)
	}
	if err = os.WriteFile(configTo, input, 0666); err != nil {
		return "", fmt.Errorf("error writing config file: %w", err)
	}
	if err = os.Remove(configAt); err != nil {
		_ = os.Remove(configTo)
		return "", fmt.Errorf("error removing original config file: %w", err)
	}

	return configTo, nil
}

// IsUsed returns true if the config file was read.
func IsUsed() bool {
	return isUsed
}

func EnsureRegistryPrefixOrExit() {
	if utf8.RuneCountInString(Config.RegistryUniqPrefix) < 10 {
		color.Red("Registry prefix \"%s\" is too short. It should be at least 10 characters long.", Config.RegistryUniqPrefix)
		shared.Exit(1)
	}
}
