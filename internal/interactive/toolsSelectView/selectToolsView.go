package toolsSelectView

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/config/toolbox"
	"github.com/KoJem9Ka/jetbrains-context-menu/internal/shared/set"
)

type selectToolsModel struct {
	choices  []toolbox.ToolModel
	cursor   uint
	selected set.Set[toolbox.ToolId]
}

func (this selectToolsModel) Update(msg tea.Msg) (selectToolsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if this.cursor > 0 {
				this.cursor--
			} else {
				this.cursor = uint(len(this.choices) - 1)
			}
		case "down":
			this.cursor = (this.cursor + 1) % uint(len(this.choices))
		case "a", /* RU */ "ф":
			if this.selected.Len() == len(this.choices) {
				this.selected = set.NewSet[toolbox.ToolId]()
			} else {
				for _, choice := range this.choices {
					this.selected.Add(choice.Id())
				}
			}
		case " ":
			choiceId := this.choices[this.cursor].Id()
			if this.selected.Has(choiceId) {
				this.selected.Remove(choiceId)
			} else {
				this.selected.Add(choiceId)
			}
		}
	}

	return this, nil
}

func (this selectToolsModel) View() string {
	var s strings.Builder
	s.WriteString("» Choose which tools to show in the context menu:\n")
	//s.WriteString(color.HiBlackString("\t↑/↓: Select a tool\n"))
	//s.WriteString(color.HiBlackString("\tSpace: Toggle selection\n"))
	//s.WriteString(color.HiBlackString("\ta: Toggle all\n"))
	//s.WriteString(color.HiBlackString("\tEnter: Confirm\n"))
	s.WriteString(color.HiBlackString("    ↑/↓: Select a tool, Space: Toggle selection\n"))
	s.WriteString(color.HiBlackString("    a: Toggle all, Enter: Confirm\n"))
	s.WriteString(color.HiBlackString("    e: Edit tool title\n"))
	for i, choice := range this.choices {
		isChecked := this.selected.Has(choice.Id())
		isCursorAtCurrentChoice := this.cursor == uint(i)

		cursor := " "
		if isCursorAtCurrentChoice {
			cursor = ">"
		}

		checked := "( )"
		title := choice.DisplayNameUsingConfig()
		installLocation := color.HiBlackString(fmt.Sprintf("(%s)", choice.Id()))

		if isChecked {
			checked = color.HiGreenString("(*)")
			title = color.GreenString(title)
		}

		s.WriteString(fmt.Sprintf("%s %s %s %s %s\n", cursor, checked, title, choice.DisplayVersion, installLocation))
	}

	return s.String()
}
