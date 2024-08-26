package toolsSelectView

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/toolbox"
)

type inputNameModel struct {
	tool      toolbox.ToolModel
	textInput textinput.Model
}

func createInputNameModel(tool toolbox.ToolModel) inputNameModel {
	textInput := textinput.New()
	textInput.SetValue(tool.DisplayNameUsingConfig())
	textInput.Placeholder = tool.DisplayName
	textInput.Focus()
	textInput.CharLimit = 40
	textInput.Validate = func(s string) error {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			return errors.New("cannot be empty")
		}
		if len(s) < 3 {
			return errors.New("must be at least 3 characters long")
		}
		return nil
	}
	textInput.Width = 40
	textInput.SetSuggestions([]string{tool.DisplayName, tool.DisplayNameUsingConfig()})
	textInput.ShowSuggestions = true
	return inputNameModel{tool, textInput}
}

func (this inputNameModel) Update(msg tea.Msg) (inputNameModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+r":
			if this.textInput.Value() != this.tool.DisplayNameUsingConfig() {
				this.textInput.SetValue(this.tool.DisplayNameUsingConfig())
			} else {
				this.textInput.SetValue(this.tool.DisplayName)
			}
		case "alt+backspace":
			this.textInput.Reset()
		case "tab":
			currentSuggestion := this.textInput.CurrentSuggestion()
			if currentSuggestion != "" {
				this.textInput.SetValue(currentSuggestion)
			} else {
				this.textInput.SetValue(this.textInput.Placeholder)
			}
		}
	}

	this.textInput, cmd = this.textInput.Update(msg)

	return this, cmd
}

func (this inputNameModel) View() string {
	var s strings.Builder
	installLocation := color.HiBlackString(fmt.Sprintf("(%s)", this.tool.InstallLocation))
	s.WriteString(fmt.Sprintf("» Input custom title for \"%s\":\n    %s\n", this.tool.DisplayName, installLocation))
	s.WriteString(color.HiBlackString("    Esc: Cancel, Enter: Confirm,\n"))
	s.WriteString(color.HiBlackString("    Ctrl+R: Reset, Alt+Backspace: Clear\n"))
	s.WriteString(this.textInput.View())
	if this.textInput.Err != nil {
		s.WriteString(fmt.Sprintf("\n  Err: %s", color.RedString(this.textInput.Err.Error())))
	}
	return s.String()
}
