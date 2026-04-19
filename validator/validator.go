package validator

import (
	"maps"

	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/parser"
	"github.com/krozhkov/mgml/validator/rules"
)

func MJMLValidator(element *parser.MJMLNode, options *core.MJMLOptions) []*rules.RuleError {
	dependencies := make(map[string]map[string]struct{}, len(globalDependencies)+len(options.Dependencies))
	for key, val := range globalDependencies {
		dependencies[key] = maps.Clone(val)
	}
	for name, deps := range options.Dependencies {
		upsertDependency(dependencies, name, deps)
	}

	skipElements := make(map[string]struct{}, len(options.SkipElements))
	for _, value := range options.SkipElements {
		skipElements[value] = struct{}{}
	}

	ruleOptions := &rules.RuleOptions{
		Components:   options.Components,
		Dependencies: dependencies,
		SkipElements: skipElements,
	}

	return validate(element, ruleOptions)
}

func validate(element *parser.MJMLNode, options *rules.RuleOptions) []*rules.RuleError {
	tagName := element.TagName
	children := element.Children

	var errors []*rules.RuleError

	var skipValidation bool
	if options.SkipElements != nil {
		_, skipValidation = options.SkipElements[tagName]
	}

	if tagName != "mjml" && !skipValidation {
		for _, rule := range MJMLRulesCollection {
			ruleError := rule(element, options)
			if ruleError != nil {
				errors = append(errors, ruleError...)
			}
		}
	}

	if len(children) > 0 {
		for _, child := range children {
			childErrors := validate(child, options)
			if childErrors != nil {
				errors = append(errors, childErrors...)
			}
		}
	}

	return errors
}
