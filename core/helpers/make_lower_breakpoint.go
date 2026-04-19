package helpers

import (
	"strconv"
)

func MakeLowerBreakpoint(breakpoint string) string {
	pixels, err := ParseIntLoose(breakpoint)
	if err != nil {
		return breakpoint
	}

	return strconv.Itoa(pixels-1) + "px"
}
