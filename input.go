package wind

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type inputModel struct {
	label string
	err   error
	input textinput.Model
}

type (
	errMsg error
)

// prompt user to select between choices , and return the selected indexes
//
// example :
//
//	wind.SelectBox("you are intersted at ", []string{"gaming", "coding"})
func InputBox(label string) (string, error) {
	ti := textinput.New()
	ti.Focus()
	ti.Prompt = ""
	ti.CharLimit = 150
	ti.Width = getTrmW() - 6

	res := tea.NewProgram(inputModel{
		label: strings.TrimSpace(label),
		input: ti,
	})

	selected, err := res.Run()

	s := selected.(inputModel)

	return s.input.Value(), err

}

func (s inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (s inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "ctrl+c", "q", "esc":
			return s, tea.Quit
		}
	case errMsg:
		s.err = msg
		return s, nil
	}

	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	return s, cmd
}

func (s inputModel) View() string {

	v := s.input.View()
	m := "\n"
	m += style().Margin(0, 0, 0, 1).Foreground(color("#495867")).Render("┌─")
	m += style().Margin(0, 1, 0, 1).Foreground(color("#07beb8")).Render(s.label)
	m += style().Foreground(color("#495867")).Render(strings.Repeat("─", getTrmW()-charWidth(s.label)-7))
	m += style().Foreground(color("#495867")).Render("┐")
	m += "\n"

	m += style().Foreground(color("#495867")).Render(" │")
	m += v
	// m += strings.Repeat(" ")
	m += style().Foreground(color("#495867")).Render(" │")

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
