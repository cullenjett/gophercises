package main

import (
	"fmt"
	"strings"
)

func main() {
	var length, delta int
	var input string
	fmt.Scanf("%d\n", &length)
	fmt.Scanf("%s\n", &input)
	fmt.Scanf("%d\n", &delta)

	alphabetLower := "abcdefghijklmnopqrstuvwxyz"
	alphabetUpper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	result := ""
	for _, ch := range input {
		switch {
		case strings.IndexRune(alphabetLower, ch) >= 0:
			result = result + string(rotate(ch, delta, []rune(alphabetLower)))
		case strings.IndexRune(alphabetUpper, ch) >= 0:
			result = result + string(rotate(ch, delta, []rune(alphabetUpper)))
		default:
			result = result + string(ch)
		}
	}

	fmt.Println(result)
}

func rotate(s rune, delta int, key []rune) rune {
	idx := -1
	for i, r := range key {
		if r == s {
			idx = i
			break
		}
	}
	if idx < 0 {
		panic("PANIC!")
	}
	idx = (idx + delta) % len(key)

	return key[idx]
}
