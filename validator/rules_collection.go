package validator

import (
	"github.com/krozhkov/mgml/parser"
	"github.com/krozhkov/mgml/validator/rules"
)

type MJMLValidationRule func(element *parser.MJMLNode, options *rules.RuleOptions) []*rules.RuleError

var MJMLRulesCollection = map[string]MJMLValidationRule{
	"validAttributes": rules.ValidateAttribute,
	"validChildren":   rules.ValidateChildren,
	"validTag":        rules.ValidateTag,
	"validTypes":      rules.ValidateType,
	"errorAttr":       rules.ErrorAttribute,
}

func RegisterRule(rule MJMLValidationRule, name string) {
	MJMLRulesCollection[name] = rule
}
