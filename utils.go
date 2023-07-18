package goclitools

import (
	"bufio"
	"fmt"
	"os"
)

func alert(label string, m string, s ...Style) string {
	t := string(Tab)
	t += theme(s...)

	return fmt.Sprint("\v", t, " ", label, " ", Reset, " ", m, "\v\n")

}

func theme(s ...Style) string {
	var theme string

	for _, i := range s {
		theme += fmt.Sprint(i)
	}

	return theme

}

func prompt() string {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("➜  ")

	if scanner.Scan() {
		return scanner.Text()
	}

	return ""
}

func PadString(s string, p int) string {
	len := len(s)

	if len > p {
		return s
	}

	var ps string

	for i := 0; i < p; i++ {
		if i < len {
			ps += string(s[i])
		} else {
			ps += " "
		}
	}

	return ps
}
