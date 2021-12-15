package stringers

import "sort"

func SortString(str string) string {
	runes := []rune(str)
	sort.Sort(runeSorter(runes))
	return string(runes)
}

type runeSorter []rune

func (s runeSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s runeSorter) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s runeSorter) Len() int {
	return len(s)
}
