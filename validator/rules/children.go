package rules

import (
	"fmt"
	"strings"

	"github.com/krozhkov/mgml/parser"
)

func ValidateChildren(element *parser.MJMLNode, options *RuleOptions) []*RuleError {
	tagName := element.TagName
	children := element.Children

	components := options.Components
	dependencies := options.Dependencies
	skipElements := options.SkipElements

	component := components[tagName]

	if component == nil || len(children) == 0 {
		return nil
	}

	var errors []*RuleError

	for _, child := range children {
		childTagName := child.TagName
		ChildComponent := components[childTagName]
		parentDependencies := dependencies[tagName]

		var shouldSkip bool
		if skipElements != nil {
			_, shouldSkip = skipElements[childTagName]
		}

		var isAllowed bool
		if parentDependencies != nil {
			_, isAllowed = parentDependencies[childTagName]
			if !isAllowed {
				_, isAllowed = parentDependencies["*"]
			}
		}

		childIsValid := ChildComponent == nil || shouldSkip || isAllowed

		if !childIsValid {
			var allowedDependencies = make([]string, 0, 10)

			for key, set := range dependencies {
				if _, ok := set[childTagName]; ok {
					allowedDependencies = append(allowedDependencies, key)
				}
			}

			errors = append(
				errors,
				NewRuleError(fmt.Sprintf("%s cannot be used inside %s, only inside: %s", childTagName, tagName, strings.Join(allowedDependencies, ", ")), element),
			)
		}
	}

	return errors
}
