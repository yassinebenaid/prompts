// wind helps you build cli applications by abstracting all terminal intractions in an easy to use
// apis and let you focus on the core of your application where it takes care of accepting the inputs and returning the output
//
// in other words , wind is the connection between your application and the terminal
//
// it provides several components to format the output and give it custom styles
// , a router to structure the commands tree, components  for accepting  user input and a logger
// that you can use to log errors in a nice styles
package wind

import "github.com/charmbracelet/lipgloss"

func mapKeys[TKey int | float64, TValue any](m map[TKey]TValue) []TKey {
	keys := []TKey{}

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func quickSort[TItem int | float64](s []TItem) []TItem {
	if len(s) <= 1 {
		return s
	}

	pivot := s[0]
	lower := []TItem{}
	greater := []TItem{}

	for _, item := range s {
		if item < pivot {
			lower = append(lower, item)
		} else if item > pivot {
			greater = append(greater, item)
		}
	}

	sorted := []TItem{}
	sorted = append(sorted, quickSort[TItem](lower)...)
	sorted = append(sorted, pivot)
	sorted = append(sorted, quickSort[TItem](greater)...)

	return sorted
}

func style() lipgloss.Style {
	return lipgloss.NewStyle()
}
func color(c string) lipgloss.Color {
	return lipgloss.Color(c)
}
