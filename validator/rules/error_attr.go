package rules

import (
	"github.com/krozhkov/mgml/parser"
)

func ErrorAttribute(element *parser.MJMLNode, options *RuleOptions) []*RuleError {
	if element.Err == nil {
		return nil
	}

	var errors = []*RuleError{NewRuleError(element.Err.Error(), element)}

	return errors
}
