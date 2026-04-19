package helpers

import (
	"regexp"
	"strconv"
	"strings"
)

func ShorthandParser(cssValue string, direction string) int {
	splittedCssValue := strings.Fields(cssValue)
	var directions map[string]int

	switch len(splittedCssValue) {
	case 2:
		directions = map[string]int{"top": 0, "bottom": 0, "left": 1, "right": 1}
	case 3:
		directions = map[string]int{"top": 0, "left": 1, "right": 1, "bottom": 2}
	case 4:
		directions = map[string]int{"top": 0, "right": 1, "bottom": 2, "left": 3}
	default:
		val, err := ParseIntLoose(cssValue)
		if err != nil {
			return 0
		}
		return val
	}

	index, ok := directions[direction]
	if !ok {
		return 0
	}

	if index >= len(splittedCssValue) {
		return 0
	}

	str := splittedCssValue[index]

	val, err := ParseIntLoose(str)
	if err != nil {
		return 0
	}

	return val
}

var borderMatcher = regexp.MustCompile(`(?:(?:^| )(\d+))`)

func BorderParser(border string) int {
	match := borderMatcher.FindStringSubmatch(border)
	if len(match) < 2 {
		return 0
	}

	value, err := strconv.Atoi(match[1])
	if err != nil {
		return 0
	}

	return value
}
