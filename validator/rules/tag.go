package rules

import (
	"fmt"

	"github.com/krozhkov/mgml/parser"
)

func ValidateTag(element *parser.MJMLNode, options *RuleOptions) []*RuleError {
	tagName := element.TagName

	// Tags that have no associated components but are allowed even so
	if tagName == "mj-all" || tagName == "mj-class" || tagName == "mj-selector" || tagName == "mj-html-attribute" {
		return nil
	}

	component := options.Components[tagName]

	if component == nil {
		var errors = []*RuleError{NewRuleError(fmt.Sprintf("element %s doesn't exist or is not registered", tagName), element)}

		return errors
	}

	return nil
}
