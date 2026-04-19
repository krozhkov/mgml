package helpers

import (
	"strings"

	"github.com/krozhkov/mgml/internal/utils"
)

func SuffixCssClasses(classes, suffix string) string {
	if classes == "" {
		return ""
	}

	return strings.Join(utils.MapFunc(strings.Split(classes, " "), func(c string) string {
		return c + "-" + suffix
	}), " ")
}
