package wind

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type inputModel struct {
	label string
	value string
}

// prompt user to select between choices , and return the selected indexes
//
// example :
//
//	wind.SelectBox("you are intersted at ", []string{"gaming", "coding"})
func InputBox(label string) (string, error) {
	res := tea.NewProgram(inputModel{
		label: strings.TrimSpace(label),
	})

	selected, err := res.Run()

	s := selected.(inputModel)

	return s.value, err

}

func (s inputModel) Init() tea.Cmd {
	return nil
}

func (s inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ", "ctrl+c", "q":
			return s, tea.Quit
		}
	}

	return s, nil
}

func (s inputModel) View() string {
	m := "\n"
	m += style().Margin(0, 0, 0, 1).Foreground(color("#495867")).Render("┌")
	m += style().Foreground(color("#495867")).Render(strings.Repeat("─", ((getTrmW() - 5 - len(s.label)) / 2)))
	m += style().Margin(0, 1, 0, 1).Foreground(color("#07beb8")).Render(s.label)
	m += style().Foreground(color("#495867")).Render(strings.Repeat("─", ((getTrmW() - 6 - len(s.label)) / 2)))
	m += style().Foreground(color("#495867")).Render("┐")
	m += "\n"

	m += style().Foreground(color("#495867")).Render(" │")
	m += strings.Repeat(" ", getTrmW()-5)
	m += style().Foreground(color("#495867")).Render(" │")
	m += "\n"
	m += style().Foreground(color("#495867")).Render(" │")
	m += strings.Repeat(" ", (getTrmW()-43)/2)

	m += strings.Repeat(" ", ((getTrmW() - 36) / 2))
	m += style().Foreground(color("#495867")).Render(" │")
	m += "\n"
	m += style().Foreground(color("#495867")).Render(" └")
	m += style().Foreground(color("#495867")).Render(strings.Repeat("─", getTrmW()-4))
	m += style().Foreground(color("#495867")).Render("┘")
	m += "\n"
	m += style().Foreground(color("#495867")).Render(" Press enter to quit.")
	m += "\n\n"
	return m
}
