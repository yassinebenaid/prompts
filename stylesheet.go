package goclitools

import (
	"errors"
	"fmt"
	"strings"
)

// Prints the message "m" in the specified style "s."
// Perfect for short messages
func Print(m string, s ...Style) {
	fmt.Print(theme(s...), m, Reset)
}

// Prints the message "m" using the "INFO" theme for styling.
// Perfect for short messages
func Info(m string) {
	Alert("INFO", m, BG_Blue)
}

// Prints the message "m" using the "ERROR" theme for styling.
// Perfect for short messages
func Error(m string) {
	Alert("ERROR", m, BG_Red)
}

// Prints m styled by SUCCESS theme
func Success(m string) {
	Alert("SUCCESS", m, BG_Green)
}

// Prints the message "m" using the "WARNING" theme for styling.
// Perfect for short messages
func Warning(m string) {
	Alert("WARNING", m, BG_BrightYellow)
}

// Prints the message "m" labeled by "label" with the specified styles "s."
// The "label" will be formatted according to the provided styles,
// enhancing its visual appearance.
func Alert(label string, m string, s ...Style) {
	t := string(Tab)
	t += theme(s...)

	fmt.Print("\v", t, " ", label, " ", Reset, " ", m, "\v\n")
}

// Prints the message "m" using the "INFO" theme for styling.
// Perfect for long messages
func InfoMessage(m string) {
	PrintMessage(m, BG_Blue)
}

// Prints the message "m" using the "ERROR" theme for styling.
// Perfect for long messages
func ErrorMessage(m string) {
	PrintMessage(m, BG_Red)
}

// Prints the message "m" using the "SUCCESS" theme for styling.
// Perfect for long messages
func SuccessMessage(m string) {
	PrintMessage(m, BG_Green)
}

// Prints the message "m" using the "WARNING" theme for styling.
// Perfect for long messages
func WarningMessage(m string) {
	PrintMessage(m, BG_BrightYellow)
}

// Prints the message "m" in the specified style "s."
// Perfect for long messages
func PrintMessage(m string, s ...Style) {
	t := theme(s...)

	fmt.Print(t, "\x1b[J\v", Tab, m, "\n", t, "\x1b[J\v")
}

// Show a progress bar , this function should be called on a loop ,
// its safe to be used in multiple goroutines
//
// you can choose show/hide the numeric percentage using prc
//
// each time you call it , it will updated the progress bar based on iteration
func ProgressBar(iteration, total float64, prc bool) {
	Progress(&ProgressStyle{
		Total:           total,
		Iteration:       iteration,
		Prc:             prc,
		BarColor:        BG_BrightWhite,
		BackgroundColor: BG_Black,
	})
}

type ProgressStyle struct {
	Iteration       float64
	Total           float64
	Label           string
	Prc             bool
	LabelStyle      []Style
	BarColor        Style
	BackgroundColor Style
}

func (p *ProgressStyle) prc() int {
	return int((p.Iteration / p.Total) * 100)
}

func (p *ProgressStyle) finished() bool {
	return p.Iteration == p.Total
}

// Show a progress bar , this function should be called on a loop ,
// its safe to be used in multiple goroutines
//
// you can choose show/hide the numeric percentage using prc
//
// Unlike ProgressBar which is usefull for simple progress bars
// this function gives more flexibility to customize the bar as you need
func Progress(s *ProgressStyle) error {

	if s.Total == 0 {
		return errors.New("invalid total value , total should be greater than 0")
	}

	var percent string
	var label string

	if s.Prc {
		percent = fmt.Sprintf("%d%%", s.prc())
	}

	if s.Label != "" {
		label = theme(s.LabelStyle...) + s.Label + string(Reset)
	}

	filledLength := int(50 * s.Iteration / s.Total)

	progress := strings.Repeat(string(s.BarColor)+" "+string(Reset), filledLength) + strings.Repeat(string(s.BackgroundColor)+" "+string(Reset), (50-filledLength))

	fmt.Printf("\r%s %s %s", label, progress, percent)

	if s.finished() {
		fmt.Println()
	}

	return nil
}

func theme(s ...Style) string {
	var theme string

	for _, i := range s {
		theme += fmt.Sprint(i)
	}

	return theme

}
