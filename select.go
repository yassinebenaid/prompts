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
		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			}
		case "down", "j":
			if s.cursor < len(s.choices)-1 {
				s.cursor++
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
	m := Sprint(" ┌── ", T_BrightBlack)
	m += Sprint(s.label, T_BrightBlue)
	m += Sprint(strings.Repeat("─", getTrmW()-charWidth(s.label)-7), T_BrightBlack)
	m += Sprint("┐", T_BrightBlack)
	m += Sprint("\n", Reset)

	for i, choice := range s.choices {
		cursor := Sprint(" |", T_BrightBlack)
		if s.cursor == i {
			cursor += Sprint(" ➜", T_BrightBlue)
		} else {
			cursor += Sprint("  ", T_BrightCyan)
		}

		checked := Sprint(" □ ", Dim)
		if _, ok := s.selected[i]; ok {
			checked = Sprint(" ▣ ", T_BrightCyan)
		}

		m += fmt.Sprintf("%s %s %s", cursor, checked, choice)
		m += strings.Repeat(" ", getTrmW()-11-len(choice))
		m += Sprint("|\n", T_BrightBlack)
	}
	m += Sprint(" └", T_BrightBlack)
	m += Sprint(strings.Repeat("─", getTrmW()-4), T_BrightBlack)
	m += Sprint("┘", T_BrightBlack)
	m += "\nPress q to quit.\n"

	return m
}

func (s selectModel) getSelected() []int {
	selected := make([]int, 0, len(s.selected))

	for i := range s.selected {
		selected = append(selected, i)
	}

	return selected
}
