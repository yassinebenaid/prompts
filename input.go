package wind

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error

	inputModel struct {
		required  bool
		validator func(value string) error
		label     string
		err       error
		input     textinput.Model
	}

	InputOptions struct {
		Label       string
		Placeholder string
		Required    bool
		Validator   func(value string) error
	}
)

// prompt user to select between choices , and return the selected indexes
//
// example :
//
//	wind.SelectBox("you are intersted at ", []string{"gaming", "coding"})
func InputBox(options InputOptions) (string, error) {
	input := getInput()
	input.Placeholder = strings.TrimSpace(options.Placeholder)

	result := tea.NewProgram(inputModel{
		label:     strings.TrimSpace(options.Label),
		required:  options.Required,
		validator: options.Validator,
		input:     input,
	})

	model, err := result.Run()

	m := model.(inputModel)

	return m.input.Value(), err

}

func (model inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (model inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "ctrl+c", "esc":
			if model.validate() {
				return model, tea.Quit
			}
			return model, nil
		}
	case errMsg:
		model.err = msg
		return model, nil
	}

	var cmd tea.Cmd
	model.input, cmd = model.input.Update(msg)

	return model, cmd
}

func (model inputModel) View() string {

	v := model.input.View()
	m := "\n"
	m += style().Margin(0, 0, 0, 1).Foreground(color("#495867")).Render("┌─")

	if len(model.label) > 0 {
		model.label = " " + model.label + " "
	}

	m += style().Foreground(color("#07beb8")).Render(model.label)
	m += style().Foreground(color("#495867")).Render(strings.Repeat("─", getTrmW()-charWidth(model.label)-5))
	m += style().Foreground(color("#495867")).Render("┐")
	m += "\n"

	m += style().Foreground(color("#495867")).Render(" │")
	m += v

	ph := model.input.Placeholder
	if len(ph) > 0 {
		if len(model.input.Value()) > 0 {
			m += style().Foreground(color("#495867")).Render(" │")
		} else {
			m += strings.Repeat(" ", getTrmW()-6-len(ph))
			m += style().Foreground(color("#495867")).Render(" │")
		}
	} else {
		m += style().Foreground(color("#495867")).Render(" │")
	}
	m += "\n"
	m += style().Foreground(color("#495867")).Render(" └")
	m += style().Foreground(color("#495867")).Render(strings.Repeat("─", getTrmW()-4))
	m += style().Foreground(color("#495867")).Render("┘")
	m += "\n"

	if model.err != nil {
		m += style().Foreground(color("#fb8500")).Italic(true).Render(" ⚠  " + model.err.Error())
	} else {
		m += style().Foreground(color("#495867")).Render(" Press enter to quit.")
	}
	m += "\n\n"
	return m
}

func getInput() textinput.Model {
	textInput := textinput.New()
	textInput.Focus()
	textInput.Prompt = " "
	textInput.CharLimit = 150
	textInput.Width = getTrmW() - 7

	return textInput
}

func (model *inputModel) validate() bool {
	if model.required && model.input.Value() == "" {
		model.err = fmt.Errorf("this field is required")
		return false
	}

	if model.validator != nil {
		model.err = model.validator(model.input.Value())
		return model.err == nil
	}

	return true
}