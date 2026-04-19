package helpers

import (
	"regexp"
)

var unitRegex = regexp.MustCompile(`[\d.,]*(\D*)$`)

func WidthParser(width string, parseFloatToInt bool) (float64, string) {
	var widthUnit string
	match := unitRegex.FindStringSubmatch(width)
	if len(match) < 2 {
		widthUnit = "px"
	} else {
		widthUnit = match[1]
	}

	parsedWidth, _ := parseFloat(width, widthUnit, parseFloatToInt)

	return parsedWidth, widthUnit
}

func parseFloat(value string, unit string, parseFloatToInt bool) (float64, error) {
	switch unit {
	case "px":
		intValue, err := ParseIntLoose(value)
		return float64(intValue), err
	case "%":
		if parseFloatToInt {
			intValue, err := ParseIntLoose(value)
			return float64(intValue), err
		} else {
			return ParseFloatLoose(value)
		}
	default:
		intValue, err := ParseIntLoose(value)
		return float64(intValue), err
	}
}
