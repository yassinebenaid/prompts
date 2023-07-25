package goclitools

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
