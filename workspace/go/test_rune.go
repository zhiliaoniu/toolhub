package main

import (
	"fmt"
	"strings"
)

func sanitize(s string) string {
	return strings.Map(sanitizeRune, s)
}

func sanitizeRune(r rune) rune {
	switch {
	case 'a' <= r && r <= 'z':
		return r
	case '0' <= r && r <= '9':
		return r
	case 'A' <= r && r <= 'Z':
		return r
	default:
		return '_'
	}
}

func main() {
	a := "ab杨胜智bc"
	for i, c := range a {
		fmt.Printf("%d: %c\n", i, c)
	}
	fmt.Println(len(a))

	fmt.Println(sanitize(a))
}
