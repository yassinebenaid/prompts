package goclitools

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

// clear the terminal . lines represents how many line to earase ,
// if lines = -1 clear all visible lines , like when you click CTRL+L
func ClearTrm(lines int) {
	if lines < 0 {
		_, lines, _ = term.GetSize(int(os.Stdin.Fd()))
	}

	fmt.Print(Earase)
	for i := 0; i < lines; i++ {
		fmt.Print(UP, Earase)
	}
	fmt.Print(Reset, "\r")
}

func Select(options map[int]string) int {
	if len(options) < 1 {
		return 0
	}

	selected := 0

	pw, err := term.ReadPassword(int(os.Stdin.Fd()))

	fmt.Println(pw, err)
	return selected
}

// Ask the user for an input and returns the result back
//
// prints m as label for the input with default styles,
// you can customize the label styles using s
func Ask(m string, s ...Style) string {
	fmt.Print(Tab)

	if len(s) > 0 {
		fmt.Print(theme(s...))
	} else {
		fmt.Print(Bold)
	}

	fmt.Print(m, Reset, "\n")

	return prompt()
}

// Ask the user for a secure input and returns the result back,
// this is commonly used for passwords and sensetive data
//
// prints m as label for the input with default styles,
// you can customize the label styles using s
func AskSecurly(m string, s ...Style) (string, error) {
	fmt.Print(Tab)

	if len(s) > 0 {
		fmt.Print(theme(s...))
	} else {
		fmt.Print(Bold)
	}

	fmt.Print(m, Reset, ":", Tab)

	res, err := term.ReadPassword(int(os.Stdin.Fd()))

	return string(res), err
}

// Ask the user for confirmation and returns the result back
//
// the user will have to type either yes or no , and the results will be as boolean
//
// you can choose a default value using def, it will be returned if the user types other words or if pressed enter
//
// prints m as label for the input with default styles,
// you can customize the label styles using s
func ConfirmDef(m string, def bool, s ...Style) bool {
	fmt.Print("\n", Tab)

	if len(s) > 0 {
		fmt.Print(theme(s...))
	} else {
		fmt.Print(Bold)
	}

	d := "no"

	if def {
		d = "yes"
	}

	fmt.Print(m, Reset, Dim, " (yes/no) ", "[", Reset, T_Yellow, d, Reset, Dim, "]", Reset, "\n")

	res := prompt()

	if res != "yes" && res != "no" {
		res = d
	}

	return res == "yes"
}

// Ask the user for confirmation and returns the result back
//
// the user will have to type either yes or no to confirm  , other words or enter will keep showing the input
// consider using ConfirmDef if you want to choose a default value
//
// you can choose a default value using def, it will be returned if the user types other words or if pressed enter
//
// prints m as label for the input with default styles,
// you can customize the label styles using s
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

func prompt() string {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("➜  ")

	if scanner.Scan() {
		return scanner.Text()
	}

	return ""
}

func getTrmW() int {
	w, _, _ := term.GetSize(int(os.Stdin.Fd()))
	return w
}
func getTrmH() int {
	_, h, _ := term.GetSize(int(os.Stdin.Fd()))
	return h
}
