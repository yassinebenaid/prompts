package prompts

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error

	inputModel struct {
		required   bool
		validator  func(value string) error
		label      string
		err        error
		input      textinput.Model
		value      string
		secure     bool
		terminated bool
	}

	InputOptions struct {
		Secure      bool
		Label       string
		Placeholder string
		Required    bool
		Validator   func(value string) error
	}
)

// prompt user to fill an input, it returns the user input string
//
// example :
//
//	prompts.InputBox(prompts.InputOptions{
//		Secure:      false, // hides the user input, very common for passwords
//		Label:       "what is your name?",
//		Placeholder: "what is your name",
//		Required:    true,
//		Validator: func(value string) error { // will be called when user submit, and returned error will be displayed to the user below the input
//			if len(value) < 3 {
//				return fmt.Errorf("minimum len is 3")
//			}
//			return nil
//		},
//	})
func InputBox(options InputOptions) (string, error) {
	input := getInput()
	input.Placeholder = strings.TrimSpace(options.Placeholder)

	result := tea.NewProgram(inputModel{
		label:     strings.TrimSpace(options.Label),
		required:  options.Required,
		validator: options.Validator,
		input:     input,
		secure:    options.Secure,
	})

	model, err := result.Run()

	m := model.(inputModel)

	return m.value, err

}

func (model inputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (model inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model.value = model.input.Value()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "ctrl+c", "esc":
			if model.validate() {
				model.terminated = true
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

	if model.secure && len(model.input.Value()) > 0 {
		model.input.SetValue(strings.Repeat("⬩", len(model.input.Value())))
	}

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
		if len(model.value) > 0 {
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
	} else if !model.terminated {
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
	if model.required && model.value == "" {
		model.err = fmt.Errorf("this field is required")
		return false
	}

	if model.validator != nil {
		model.err = model.validator(model.value)
		return model.err == nil
	}

	return true
}
