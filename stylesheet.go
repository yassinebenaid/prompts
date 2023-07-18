package goclitools

import (
	"errors"
	"fmt"
	"strings"
)

func Text(m string, s ...Style) {
	fmt.Print(theme(s...), m, Reset)
}

func Info(m string) {
	fmt.Print(alert("INFO", m, BackgroundBlue))
}

func Error(m string) {
	fmt.Print(alert("ERROR", m, BackgroundRed))
}

func Success(m string) {
	fmt.Print(alert("SUCCESS", m, BackgroundGreen))
}

func Warning(m string) {
	fmt.Print(alert("WARNING", m, BackgroundBrightYellow))
}

func InfoMessage(m string) {
	Message(m, BackgroundBlue)
}

func ErrorMessage(m string) {
	Message(m, BackgroundRed)
}

func SuccessMessage(m string) {
	Message(m, BackgroundGreen)
}

func WarningMessage(m string) {
	Message(m, BackgroundBrightYellow)
}

func Message(m string, s ...Style) {
	t := theme(s...)

	fmt.Print(t, "\x1b[J\v", Tab, m, "\n", t, "\x1b[J\v")
}

func Ask(m string, styles ...Style) string {
	fmt.Print("\n", Tab)

	if len(styles) > 0 {
		fmt.Print(theme(styles...))
	} else {
		fmt.Print(Bold)
	}

	fmt.Print(m, Reset, "\n")

	return prompt()
}

func ConfirmDef(m string, def bool, styles ...Style) bool {
	fmt.Print("\n", Tab)

	if len(styles) > 0 {
		fmt.Print(theme(styles...))
	} else {
		fmt.Print(Bold)
	}

	d := "no"

	if def {
		d = "yes"
	}

	fmt.Print(m, Reset, Dim, " (yes/no) ", "[", Reset, TextYellow, d, Reset, Dim, "]", Reset, "\n")

	res := prompt()

	if res != "yes" && res != "no" {
		res = d
	}

	return res == "yes"
}

func Confirm(m string, styles ...Style) bool {
	fmt.Print("\n", Tab)

	if len(styles) > 0 {
		for _, i := range styles {
			fmt.Print(i)
		}
	} else {
		fmt.Print(Bold)
	}

	fmt.Print(m, Reset, Dim, " (yes/no) ", Reset, "\n")

	var res string

	for res != "yes" && res != "no" {
		res = prompt()
	}

	return res == "yes"
}

func ProgressBar(iteration, total float64, prc bool) {
	Progress(&ProgressStyle{
		Total:      total,
		Iteration:  iteration,
		Prc:        prc,
		Bar:        BackgroundBrightWhite,
		Background: BackgroundBlack,
	})
}

type ProgressStyle struct {
	Iteration  float64
	Total      float64
	Label      string
	Prc        bool
	LabelStyle []Style
	Bar        Style
	Background Style
}

func (p *ProgressStyle) prc() int {
	return int((p.Iteration / p.Total) * 100)
}

func (p *ProgressStyle) finished() bool {
	return p.Iteration == p.Total
}

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

	progress := strings.Repeat(string(s.Bar)+" "+string(Reset), filledLength) + strings.Repeat(string(s.Background)+" "+string(Reset), (50-filledLength))

	fmt.Printf("\r%s %s %s", label, progress, percent)

	if s.finished() {
		fmt.Println()
	}

	return nil
}
