package rules

import (
	"fmt"
	"strings"

	"github.com/krozhkov/mgml/parser"
)

type RuleError struct {
	Line             int
	TagName          string
	Message          string
	FormattedMessage string
}

func NewRuleError(message string, element *parser.MJMLNode) *RuleError {
	line := element.Line
	tagName := element.TagName
	absoluteFilePath := element.AbsoluteFilePath
	formattedMessage := fmt.Sprintf("Line %d of %s%s (%s) — %s", line, absoluteFilePath, formatInclude(element), tagName, message)

	return &RuleError{line, tagName, message, formattedMessage}
}

func (r *RuleError) Error() string {
	return r.FormattedMessage
}

func formatInclude(element *parser.MJMLNode) string {
	includedIn := element.IncludedIn

	if len(includedIn) == 0 {
		return ""
	}

	formattedIncluded := make([]string, 0, len(includedIn))
	for i := len(includedIn) - 1; i >= 0; i-- {
		line := includedIn[i].Line
		file := includedIn[i].File

		formattedIncluded = append(formattedIncluded, fmt.Sprintf("line %d of file %s", line, file))
	}

	return fmt.Sprintf(", included at %s", strings.Join(formattedIncluded, ", itself included at "))
}
