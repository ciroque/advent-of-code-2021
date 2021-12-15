package stringers

import "sort"

func SortStringReverse(str string) string {
	runes := []rune(str)
	sort.Sort(runeReverseSorter(runes))
	return string(runes)
}

type runeReverseSorter []rune

func (s runeReverseSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s runeReverseSorter) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s runeReverseSorter) Len() int {
	return len(s)
}
