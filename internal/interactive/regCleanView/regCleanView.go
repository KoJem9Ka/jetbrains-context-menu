package regCleanView

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"golang.org/x/sys/windows/registry"

	"github.com/KoJem9Ka/jetbrains-context-menu/internal/shared"
	"github.com/KoJem9Ka/jetbrains-context-menu/internal/shared/set"
)

type RegRow struct {
	Key         registry.Key
	MainPath    string
	CommandPath string
	Title       string
	Command     string
}

func (this RegRow) id() string {
	return this.MainPath
}

func MustSelectRegKeysToDelete(regEntities []RegRow) []RegRow {
	modelResult, err := tea.NewProgram(createModel(regEntities)).Run()
	if err != nil {
		color.Red("Internal error on selecting registry rows: %v\n", err)
		shared.Exit(1)
	}
	modelCasted, ok := modelResult.(model)
	if !ok {
		color.Red("Internal error on casting result: %v\n", err)
		shared.Exit(1)
	}
	result := make([]RegRow, 0, modelCasted.selected.Len())
	for _, choice := range modelCasted.choices {
		if modelCasted.selected.Has(choice.id()) {
			result = append(result, choice)
		}
	}
	return result
}

type model struct {
	choices    []RegRow
	cursor     uint
	selected   set.Set[string]
	isQuitting bool
}

func createModel(regEntities []RegRow) model {
	return model{
		choices:  regEntities,
		selected: set.NewSet[string](),
	}
}

func (this model) Init() tea.Cmd {
	return nil
}

func (this model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			this.selected.Clear()
			fallthrough
		case "enter":
			this.isQuitting = true
			return this, tea.Quit
		case "up", "k":
			if this.cursor > 0 {
				this.cursor--
			} else {
				this.cursor = uint(len(this.choices) - 1)
			}
		case "down", "j":
			this.cursor = (this.cursor + 1) % uint(len(this.choices))
		case "a":
			if this.selected.Len() == len(this.choices) {
				this.selected = set.NewSet[string]()
			} else {
				for _, choice := range this.choices {
					this.selected.Add(choice.id())
				}
			}
		case " ":
			choiceId := this.choices[this.cursor].id()
			if this.selected.Has(choiceId) {
				this.selected.Remove(choiceId)
			} else {
				this.selected.Add(choiceId)
			}
		}
	}

	return this, nil
}

func (this model) View() string {
	if this.isQuitting {
		return ""
	}

	var s strings.Builder
	s.WriteString("» Choose which registry rows to delete:\n")
	s.WriteString(color.HiBlackString("\t(↑/↓) Select a tool, (Space) Toggle selection\n\t(a) Toggle all, (Enter) Confirm\n"))
	for i, choice := range this.choices {
		isChecked := this.selected.Has(choice.id())
		isCursorAtCurrentChoice := this.cursor == uint(i)

		cursor := " "
		if isCursorAtCurrentChoice {
			cursor = ">"
		}

		checked := "( )"
		title := choice.Title

		if isChecked {
			checked = color.RedString("(*)")
			title = color.RedString(title)
		}

		description := color.HiBlackString(fmt.Sprintf("\tRegPath: %s\n\tCommand: %s", choice.MainPath, choice.Command))
		s.WriteString(fmt.Sprintf("%s %s %s\n%s\n", cursor, checked, title, description))
	}

	return s.String()
}
