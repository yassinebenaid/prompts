package prompts

import (
	"fmt"
	"regexp"

	"github.com/charmbracelet/lipgloss"
)

// Prints the message "m" using the "INFO" theme for styling.
// Perfect for short messages
func Info(m string) {
	Alert("INFO", "4", m)
}

// Prints the message "m" using the "ERROR" theme for styling.
// Perfect for short messages
func Error(m string) {
	Alert("ERROR", "1", m)
}

// Prints m styled by SUCCESS theme
func Success(m string) {
	Alert("SUCCESS", "2", m)
}

// Prints the message "m" using the "WARNING" theme for styling.
// Perfect for short messages
func Warning(m string) {
	Alert("WARNING", "#fca311", m)
}

// Prints the alert "m" labeled by "label" with the specified color
//
// example  :
//
//	prompts.Alert("WARNING", "#fca311", "hello world")
func Alert(label string, color string, m string) {
	label = style().
		Padding(0, 1, 0, 1).
		Margin(0, 2, 0, 0).
		Background(lipgloss.Color(color)).
		Render(label)

	fmt.Println(style().
		Margin(1, 0, 1, 2).
		Render(label + m))
}

// Prints the message "m" using the "INFO" theme for styling.
// Perfect for long messages
func InfoMessage(m string) {
	Message(m, "4")
}

// Prints the message "m" using the "ERROR" theme for styling.
// Perfect for long messages
func ErrorMessage(m string) {
	Message(m, "2")
}

// Prints the message "m" using the "SUCCESS" theme for styling.
// Perfect for long messages
func SuccessMessage(m string) {
	Message(m, "1")
}

// Prints the message "m" using the "WARNING" theme for styling.
// Perfect for long messages
func WarningMessage(m string) {
	Message(m, "#fca311")
}

// Prints the message "m" labeled by "label" with the specified color
//
// example  :
//
//	prompts.Message("hello world", "#fca311")
func Message(m string, color string) {
	fmt.Println(style().
		Margin(1, 0, 1, 2).
		Padding(1, 0, 1, 2).
		Width(getTrmW() - 2).
		Background(lipgloss.Color(color)).
		Render(m))
}

func charWidth(m string) int {
	r := regexp.MustCompile(`\x1b\[([^m]+)m`)
	return len(r.ReplaceAll([]byte(m), []byte("")))
}
