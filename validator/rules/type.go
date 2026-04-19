package rules

import (
	"fmt"

	"github.com/krozhkov/mgml/core/types"
	"github.com/krozhkov/mgml/parser"
)

func ValidateType(element *parser.MJMLNode, options *RuleOptions) []*RuleError {
	tagName := element.TagName
	attributes := element.Attributes

	component := options.Components[tagName]

	if component == nil {
		return nil
	}

	var errors []*RuleError

	if attributes != nil {
		for attr, value := range attributes.AllFromFront() {
			allowedAttribs := component.AllowedAttributes
			attrType := allowedAttribs[attr]
			if attrType != "" {
				typeChecker, err := types.InitializeType(attrType)
				if err != nil {
					errors = append(errors, NewRuleError(err.Error(), element))
				}

				result := typeChecker(value)
				if !result.IsValid() {
					errors = append(errors, NewRuleError(fmt.Sprintf("Attribute %s %v", attr, result.GetError()), element))
				}
			}
		}
	}

	return errors
}
