// wind helps you build cli applications as soon as possible;
// it provides you with bunch of pre built components you can use to intract with the terminal
//
// it provides several components to format the output, give it custom styles
// custom logger, router and deal with user input , all in one package
package wind

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
