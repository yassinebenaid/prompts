package prompts

import (
	"strings"
)

func routerHelp(router *Router) string {
	var termWidth = getTrmW()
	var help string
	width := 0

	for prefix := range router.routes {
		width = max(len(prefix), width, 10)
	}

	for prefix := range router.groups {
		width = max(len(prefix), width, 10)
	}

	help += style().Padding(1, 0, 0, 1).Render(router.config.Name)
	help += style().Padding(0, 1, 0, 1).Bold(true).Render(router.config.Version)
	help += style().Render(router.config.Description)
	help += "\n"

	help += style().Padding(1, 0, 0, 1).Bold(true).Render("COMMANDS :")
	help += "\n"

	for prefix, route := range router.routes {
		help += style().PaddingLeft(1).Width(width).Foreground(color("#07beb8")).Render(prefix)
		help += style().PaddingLeft(1).Render(formatDescription(termWidth, width, route.description))
		help += "\n"
	}

	for prefix, route := range router.groups {
		help += style().PaddingLeft(1).Width(width).Foreground(color("#07beb8")).Render(prefix)
		help += style().PaddingLeft(1).Render(formatDescription(termWidth, width, route.description))
		help += "\n"
	}

	return help
}

func formatDescription(tw int, w int, d string) string {

	if len(d) < tw-w {
		return d
	}

	var chunked string
	var cols int

	for _, letter := range d {
		if cols >= tw-w-2 {
			chunked += "\n" + strings.Repeat(" ", w+1)
			cols = 0
		}

		chunked += string(letter)
		cols++
	}

	return chunked
}
