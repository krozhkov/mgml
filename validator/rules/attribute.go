package rules

import (
	"strings"

	"github.com/krozhkov/mgml/parser"
)

func ValidateAttribute(element *parser.MJMLNode, options *RuleOptions) []*RuleError {
	tagName := element.TagName
	attributes := element.Attributes

	component := options.Components[tagName]

	if component == nil {
		return nil
	}

	allowedAttribs := component.AllowedAttributes

	var unknownAttributes = make([]string, 0, 10)

	if attributes != nil {
		for attribute := range attributes.AllFromFront() {
			_, ok := allowedAttribs[attribute]
			if !ok && attribute != "mj-class" && attribute != "css-class" {
				unknownAttributes = append(unknownAttributes, attribute)
			}
		}
	}

	if len(unknownAttributes) == 0 {
		return nil
	}

	var attribute, illegal string

	if len(unknownAttributes) > 1 {
		attribute = "Attributes"
		illegal = "are illegal"
	} else {
		attribute = "Attribute"
		illegal = "is illegal"
	}

	var errors = []*RuleError{NewRuleError(attribute+" "+strings.Join(unknownAttributes, ", ")+" "+illegal, element)}

	return errors
}
