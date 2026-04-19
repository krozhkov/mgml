package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type widthParserTest struct {
	input           string
	parseFloatToInt bool
	outputWidth     float64
	outputUnit      string
}

func TestWidthParser(t *testing.T) {
	var testValues = []widthParserTest{
		{
			input:           "1px",
			parseFloatToInt: true,
			outputWidth:     1,
			outputUnit:      "px",
		},
		{
			input:           "33.3px",
			parseFloatToInt: true,
			outputWidth:     33,
			outputUnit:      "px",
		},
		{
			input:           "33.3%",
			parseFloatToInt: true,
			outputWidth:     33,
			outputUnit:      "%",
		},
		{
			input:           "33.3%",
			parseFloatToInt: false,
			outputWidth:     33.3,
			outputUnit:      "%",
		},
	}

	for _, tt := range testValues {
		t.Run(fmt.Sprintf("should parse %s", tt.input), func(t *testing.T) {
			value, unit := WidthParser(tt.input, tt.parseFloatToInt)
			assert.Equal(t, tt.outputWidth, value)
			assert.Equal(t, tt.outputUnit, unit)
		})
	}
}
