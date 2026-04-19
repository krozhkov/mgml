package helpers

import (
	"strconv"
	"unicode"
)

func ParseIntLoose(s string) (int, error) {
	for i, c := range s {
		if !unicode.IsDigit(c) {
			return strconv.Atoi(s[:i])
		}
	}

	return strconv.Atoi(s)
}

func ParseFloatLoose(s string) (float64, error) {
	for i, c := range s {
		if !unicode.IsDigit(c) && c != '.' {
			return strconv.ParseFloat(s[:i], 64)
		}
	}

	return strconv.ParseFloat(s, 64)
}

func TrimNumber(s string) string {
	for i, c := range s {
		if !unicode.IsDigit(c) && c != '.' {
			return s[:i]
		}
	}

	return s
}
