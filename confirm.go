package wind

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type confirmModel struct {
	label     string
	confirmed bool
}

// prompt user to select between choices , and return the selected indexes
//
// example :
//
//	wind.SelectBox("you are intersted at ", []string{"gaming", "coding"})
func ConfirmBox(label string, def bool) (bool, error) {
	res := tea.NewProgram(confirmModel{
		label:     strings.TrimSpace(label),
		confirmed: def,
	})

	selected, err := res.Run()

	s := selected.(confirmModel)

	return s.confirmed, err

}

func (s confirmModel) Init() tea.Cmd {
	return nil
}

func (s confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k", "left":
			s.confirmed = !s.confirmed
		case "down", "j", "right":
			s.confirmed = !s.confirmed
		case "enter", " ", "ctrl+c", "q":
			return s, tea.Quit
		}

	}

	return s, nil
}

func (s confirmModel) View() string {
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

	if s.confirmed {
		m += style().Padding(0, 5, 0, 5).Bold(true).Foreground(color("#4f5d75")).Render("No")
		m += strings.Repeat(" ", 10)
		m += style().Padding(0, 5, 0, 5).Bold(true).Background(color("2")).Render("Yes")
	} else {
		m += style().Padding(0, 5, 0, 5).Bold(true).Background(color("#4f5d75")).Render("No")
		m += strings.Repeat(" ", 10)
		m += style().Padding(0, 5, 0, 5).Bold(true).Foreground(color("2")).Render("Yes")
	}

	m += strings.Repeat(" ", ((getTrmW() - 36) / 2))
	m += style().Foreground(color("#495867")).Render(" │")
	m += "\n"
	m += style().Foreground(color("#495867")).Render(" └")
	m += style().Foreground(color("#495867")).Render(strings.Repeat("─", getTrmW()-4))
	m += style().Foreground(color("#495867")).Render("┘")
	m += "\n"
	m += style().Foreground(color("#444")).Render(" ⇽ /⇾  move cursor   ⬩  ↵  submit.")
	m += "\n\n"
	return m
}
