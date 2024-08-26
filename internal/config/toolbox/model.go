package toolbox

import (
	"hash/crc32"
	"path/filepath"
	"strconv"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/config"
)

var State state

type state struct {
	//Version    int         `json:"version"`
	AppVersion string      `json:"appVersion"`
	Tools      []ToolModel `json:"tools"`
}

type ToolModel struct {
	//ChannelId       string `json:"channelId"`
	//ToolId          string `json:"toolId"`
	//ProductCode     string `json:"productCode"`
	//Tag             string `json:"tag"`
	DisplayName    string `json:"displayName"`
	DisplayVersion string `json:"displayVersion"`
	//BuildNumber     string `json:"buildNumber"`
	InstallLocation string `json:"installLocation"`
	LaunchCommand   string `json:"launchCommand"`
}

type ToolId string

func (this ToolModel) DisplayNameUsingConfig() string {
	displayNameUsingConfig, ok := config.Config.ToolsTitles[this.IdStr()]
	if ok {
		return displayNameUsingConfig
	}
	return this.DisplayName
}

func (this ToolModel) SetDisplayName(name string) {
	config.Config.ToolsTitles[this.IdStr()] = name
}

func (this ToolModel) Id() ToolId {
	return ToolId(this.InstallLocation)
}

func (this ToolModel) IdStr() string {
	return string(this.Id())
}

func (this ToolModel) ExecutablePath() string {
	return filepath.Join(this.InstallLocation, this.LaunchCommand)
}

func (this ToolModel) IdHash() string {
	crc32Sum := crc32.ChecksumIEEE([]byte(this.IdStr()))
	return strconv.FormatUint(uint64(crc32Sum), 10)
}
