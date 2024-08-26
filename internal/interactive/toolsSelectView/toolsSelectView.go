package toolsSelectView

import (
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/toolbox"
	"github.com/KoJem9Ka/jetbrains-context-menu/internal/registryManager"
	"github.com/KoJem9Ka/jetbrains-context-menu/internal/shared"
	"github.com/KoJem9Ka/jetbrains-context-menu/internal/shared/set"
)

func MustSelectTools(tools []toolbox.ToolModel) (selected set.Set[toolbox.ToolId], isCancelled bool, isConfigNeedSave bool) {
	result, err := tea.NewProgram(createMainModel(tools)).Run()
	if err != nil {
		color.Red("Internal error on selecting tools: %v\n", err)
		shared.Exit(1)
	}
	mainModel, ok := result.(mainModel)
	if !ok {
		color.Red("Internal error on casting result: %v\n", err)
		shared.Exit(1)
	}
	return mainModel.selectToolsModel.selected, mainModel.isCancelled, mainModel.needSaveConfig
}

type activeModel uint

const (
	selectToolsModelId activeModel = iota
	inputNameModelId
)

type mainModel struct {
	activeModel      activeModel
	selectToolsModel selectToolsModel
	inputNameModel   inputNameModel
	isCancelled      bool
	isQuitting       bool
	needSaveConfig   bool
}

func createMainModel(tools []toolbox.ToolModel) mainModel {
	selected := set.NewSet[toolbox.ToolId]()

	for _, tool := range tools {
		isHas, _ := registryManager.Has(tool)
		if isHas {
			selected.Add(tool.Id())
		}
	}

	return mainModel{
		selectToolsModel: selectToolsModel{
			choices:  tools,
			selected: selected,
		},
	}
}

func (this mainModel) Init() tea.Cmd {
	return nil
}

func (this mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			switch this.activeModel {
			case selectToolsModelId:
				return cancel(this)
			case inputNameModelId:
				this.activeModel = selectToolsModelId
			}
		case "ctrl+c":
			return cancel(this)
		case "enter":
			switch this.activeModel {
			case inputNameModelId:
				value := strings.TrimSpace(this.inputNameModel.textInput.Value())
				if this.inputNameModel.textInput.Err == nil && value != "" {
					this.activeModel = selectToolsModelId
					if value != this.inputNameModel.tool.DisplayNameUsingConfig() {
						this.inputNameModel.tool.SetDisplayName(value)
						this.needSaveConfig = true
					}
				}
			case selectToolsModelId:
				return quit(this)
			}
		case "e", /* RU */ "у":
			if this.activeModel == selectToolsModelId {
				this.inputNameModel = createInputNameModel(this.selectToolsModel.choices[this.selectToolsModel.cursor])
				this.activeModel = inputNameModelId
				return this, tea.Batch(cmds...)
			}
		}
	}

	switch this.activeModel {
	case selectToolsModelId:
		this.selectToolsModel, cmd = this.selectToolsModel.Update(msg)
		cmds = append(cmds, cmd)
	case inputNameModelId:
		this.inputNameModel, cmd = this.inputNameModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return this, tea.Batch(cmds...)
}

func (this mainModel) View() string {
	if this.isQuitting {
		return ""
	}

	switch this.activeModel {
	case selectToolsModelId:
		return this.selectToolsModel.View()
	case inputNameModelId:
		return this.inputNameModel.View()
	}

	return color.RedString("Internal error: unknown active interactive model")
}

// Utils

func cancel(model mainModel) (mainModel, tea.Cmd) {
	model.isCancelled = true
	return quit(model)
}

func quit(model mainModel) (mainModel, tea.Cmd) {
	model.isQuitting = true
	return model, tea.Quit
}
