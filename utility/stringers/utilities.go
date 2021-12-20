package stringers

import "fmt"

func ContainsAllCharacters(str, substr string) bool {
	foundCount := 0
	for _, toFind := range substr {
		for _, in := range str {
			if toFind == in {
				foundCount++
			}
		}
	}
	return foundCount == len(substr)
}

func PrintStringSlice(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}
