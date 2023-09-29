// prompts helps you build cli applications by abstracting all terminal intractions in an easy to use
// apis and let you focus on the core of your application where it takes care of accepting the inputs and returning the output
//
// in other words , prompts is the connection between your application and the terminal
//
// it provides several components to format the output and give it custom styles
// , a router to structure the commands tree, components  for accepting  user input and a logger
// that you can use to log errors in a nice styles
package prompts

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func style() lipgloss.Style {
	return lipgloss.NewStyle()
}

func color(c string) lipgloss.Color {
	return lipgloss.Color(c)
}

func getTrmW() int {
	w, _, _ := term.GetSize(int(os.Stdin.Fd()))
	return w
}
