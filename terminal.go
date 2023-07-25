// wind helps you build cli applications by abstracting all terminal intractions in an easy to use
// apis and let you focus on the core of your application where it takes care of accepting the inputs and returning the output
//
// in other words , wind is the connection between your application and the terminal
//
// it provides several components to format the output and give it custom styles
// , a router to structure the commands tree, components  for accepting  user input and a logger
// that you can use to log errors in a nice styles
package wind

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// clear the terminal . lines represents how many line to earase ,
// if lines = -1 clear all visible lines , like when you click CTRL+L
func ClearTrm(lines int) {
	if lines < 0 {
		lines = getTrmH()
	}

	fmt.Print(Earase)
	for i := 0; i < lines; i++ {
		fmt.Print(UP, Earase)
	}
	fmt.Print(Reset, "\r")
}

// prompt the user to select between a list of options and return the selected key
func Select(label string, options map[int]string) int {
	if len(options) < 1 {
		return 0
	}
	selected := -1
	width := getTrmW()
	keys := quickSort[int](mapKeys[int](options))

	Println(label+"\n", Tab, Bold)

	for _, v := range keys {
		value := Sprint(options[v] + " ")
		key := Sprint(fmt.Sprintf(" %d", v))
		dmt := strings.Repeat(".", width-charWidth(value)-charWidth(key)-4)
		dmt = Sprint(dmt, T_BrightBlack)
		fmt.Println("  " + value + dmt + key)
	}

	for selected == -1 {
		ansr := prompt()
		intansr, err := strconv.ParseInt(ansr, 10, 64)

		if _, ok := options[int(intansr)]; err == nil && ok {
			selected = int(intansr)
		}
	}

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
