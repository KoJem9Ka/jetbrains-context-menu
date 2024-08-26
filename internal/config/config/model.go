package config

const RegistryOldPrefix = "OPENWITHKOJEM9KA"

var Config = configModel{
	RegistryUniqPrefix: "JCM_OPEN_WITH_BY_KOJEM9KA",
	Title: titleModel{
		OpenFileIn:          `Открыть файл в &{}`,
		OpenDirIn:           `Открыть директорию в &{}`,
		OpenBackgroundDirIn: `Открыть здесь &{}`,
	},
	ToolboxStateJsonPath: ToolboxStateDefaultPaths[0],
	ToolsSearchPaths:     []string{},
	ToolsTitles:          map[string]string{},
	CustomTools:          map[string]string{},
}

type configModel struct {
	RegistryUniqPrefix   string            `yaml:"registry-uniq-prefix"`
	Title                titleModel        `yaml:"title"`
	ToolboxStateJsonPath string            `yaml:"toolbox-state-json-path"`
	ToolsSearchPaths     []string          `yaml:"tools-search-paths"`
	ToolsTitles          map[string]string `yaml:"tools-titles"`
	CustomTools          map[string]string `yaml:"custom-tools"`
}

type titleModel struct {
	OpenFileIn          string `yaml:"open-file-in"`
	OpenDirIn           string `yaml:"open-dir-in"`
	OpenBackgroundDirIn string `yaml:"open-background-dir-in"`
}
