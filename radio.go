package wind

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type radioModel struct {
	choices  []string
	cursor   int
	label    string
	selected int
}

// prompt user to select between choices , and return the selected indexes
//
// example :
//
//	wind.SelectBox("you are intersted at ", []string{"gaming", "coding"})
func RadioBox(label string, choices []string) (int, error) {
	res := tea.NewProgram(radioModel{
		choices: choices,
		label:   label,
	})

	selected, err := res.Run()

	s := selected.(radioModel)

	return s.getSelected(), err

}

func (s radioModel) Init() tea.Cmd {
	return nil
}

func (s radioModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return s, tea.Quit
		case "up", "k", "left":
			if s.cursor > 0 {
				s.cursor--
			} else {
				s.cursor = len(s.choices) - 1
			}
		case "down", "j", "right":
			if s.cursor < len(s.choices)-1 {
				s.cursor++
			} else {
				s.cursor = 0
			}
		case "enter", " ":
			s.selected = s.cursor
		}

	}

	return s, nil
}

func (s radioModel) View() string {
	m := "\n"
	m += style().Margin(0, 0, 0, 1).Foreground(color("#495867")).Render("┌─")
	m += style().Margin(0, 1, 0, 1).Foreground(color("#07beb8")).Render(s.label)
	m += style().Foreground(color("#495867")).Render(strings.Repeat("─", getTrmW()-charWidth(s.label)-7))
	m += style().Foreground(color("#495867")).Render("┐")
	m += "\n"

	for i, choice := range s.choices {

		cursor := style().Foreground(color("#495867")).Render(" │")
		if s.cursor == i {
			cursor += style().Foreground(color("#07beb8")).Render("➧")
		} else {
			cursor += " "
		}

		checked := style().Foreground(color("#495867")).Render("○ ")
		if s.selected == i {
			checked = style().Foreground(color("#07beb8")).Render("◉ ")
		}

		m += fmt.Sprintf("%s %s %s", cursor, checked, choice)
		m += strings.Repeat(" ", getTrmW()-9-len(choice))
		m += style().Foreground(color("#495867")).Render("│")
		m += "\n"
	}
	m += style().Foreground(color("#495867")).Render(" └")
	m += style().Foreground(color("#495867")).Render(strings.Repeat("─", getTrmW()-4))
	m += style().Foreground(color("#495867")).Render("┘")
	m += "\n"
	m += style().Foreground(color("#495867")).Render(" Press enter to select and q to quit.")
	m += "\n\n"
	return m
}

func (s radioModel) getSelected() int {
	return s.selected
}