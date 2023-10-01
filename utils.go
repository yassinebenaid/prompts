// prompts offers several configurable, ituitive and intractive CLI components
// including inputs, selectbox, radio box and alerts formatter, just to name a few
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
