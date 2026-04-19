package helpers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type shorthandTest struct {
	input  string
	output map[string]int
}

func TestShorhand(t *testing.T) {
	var testValues = []shorthandTest{
		{
			input:  "1px",
			output: map[string]int{"top": 1, "right": 1, "bottom": 1, "left": 1},
		},
		{
			input:  "1px 0",
			output: map[string]int{"top": 1, "right": 0, "bottom": 1, "left": 0},
		},
		{
			input:  "1px 2px 3px",
			output: map[string]int{"top": 1, "right": 2, "bottom": 3, "left": 2},
		},
		{
			input:  "1px 2px 3px 4px",
			output: map[string]int{"top": 1, "right": 2, "bottom": 3, "left": 4},
		},
		{
			input:  " 1px 2px  3px 4px ",
			output: map[string]int{"top": 1, "right": 2, "bottom": 3, "left": 4},
		},
	}
	var directions = []string{"top", "right", "bottom", "left"}

	for _, tt := range testValues {
		for _, dir := range directions {
			t.Run(fmt.Sprintf("should parse %s", tt.input), func(t *testing.T) {
				assert.Equal(t, tt.output[dir], ShorthandParser(tt.input, dir))
			})
		}
	}
}
