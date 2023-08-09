package wind

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type selectModel struct {
	choices  []string
	cursor   int
	label    string
	selected map[int]struct{}
}

type SelectOptions struct {
	Choices []string
	Label   string
}

// prompt user to select between choices , and return the selected indexes
//
// example :
//
//	wind.SelectBox("you are intersted at ", []string{"gaming", "coding"})
func SelectBox(label string, choices []string) ([]int, error) {
	res := tea.NewProgram(selectModel{
		choices:  choices,
		label:    label,
		selected: map[int]struct{}{},
	})

	selected, err := res.Run()

	s := selected.(selectModel)

	return s.getSelected(), err

}

func (s selectModel) Init() tea.Cmd {
	return nil
}

func (s selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			_, ok := s.selected[s.cursor]
			if ok {
				delete(s.selected, s.cursor)
			} else {
				s.selected[s.cursor] = struct{}{}
			}
		}

	}

	return s, nil
}

func (s selectModel) View() string {
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

		checked := style().Foreground(color("#495867")).Render("☐ ")
		if _, ok := s.selected[i]; ok {
			checked = style().Foreground(color("#07beb8")).Render("☑ ")
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
	m += style().Foreground(color("#444")).Render(" q quit   ⬩   ↵  select.")
	m += "\n\n"
	return m
}

func (s selectModel) getSelected() []int {
	selected := make([]int, 0, len(s.selected))

	for i := range s.selected {
		selected = append(selected, i)
	}

	return selected
}
