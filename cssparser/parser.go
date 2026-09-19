package cssparser

import (
	"errors"
	"io"
	"iter"
	"strings"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

type Declaration struct {
	Property string
	Value    string
}

type Rule struct {
	Selector     string
	Declarations []*Declaration
}

func ParseCss(input string) ([]*Rule, error) {
	result := make([]*Rule, 0)
	parser := css.NewParser(parse.NewInputString(input), false)

	var selectorStack []string
	var ruleStack []*Rule

	for {
		gt, _, data := parser.Next()
		if gt == css.ErrorGrammar {
			if err := parser.Err(); err != nil && err != io.EOF {
				return nil, err
			}
			if len(ruleStack) > 0 {
				return nil, errors.New("parser is in an invalid state")
			}
			break
		}

		switch gt {
		case css.AtRuleGrammar, css.BeginAtRuleGrammar:
			return nil, errors.New("parser does not support At-rules")

		case css.BeginRulesetGrammar:
			var sb strings.Builder
			for _, val := range parser.Values() {
				sb.Write(val.Data)
			}
			rawSelector := strings.TrimSpace(sb.String())

			var resolvedSelector string
			if len(selectorStack) == 0 {
				resolvedSelector = rawSelector
			} else {
				parentSelector := selectorStack[len(selectorStack)-1]
				resolvedSelector = resolveNesting(parentSelector, rawSelector)
			}

			selectorStack = append(selectorStack, resolvedSelector)
			rule := &Rule{
				Selector:     resolvedSelector,
				Declarations: make([]*Declaration, 0),
			}
			ruleStack = append(ruleStack, rule)
			result = append(result, rule)

		case css.EndRulesetGrammar:
			if len(selectorStack) > 0 {
				selectorStack = selectorStack[:len(selectorStack)-1]
			}
			if len(ruleStack) > 0 {
				ruleStack = ruleStack[:len(ruleStack)-1]
			}

		case css.DeclarationGrammar, css.CustomPropertyGrammar:
			if len(ruleStack) == 0 {
				continue
			}
			property := strings.TrimSpace(string(data))
			var sb strings.Builder
			for _, val := range parser.Values() {
				sb.Write(val.Data)
			}
			value := strings.TrimSpace(sb.String())

			currentRule := ruleStack[len(ruleStack)-1]
			currentRule.Declarations = append(currentRule.Declarations, &Declaration{
				Property: property,
				Value:    value,
			})
		}
	}

	return result, nil
}

func resolveNesting(parent, child string) string {
	parents := splitSelectors(parent)
	children := splitSelectors(child)

	var resolved []string
	for _, c := range children {
		hasAmp := hasAmp(c)
		for _, p := range parents {
			if hasAmp {
				resolved = append(resolved, replaceAmp(c, p))
			} else {
				resolved = append(resolved, p+" "+c)
			}
		}
	}

	return strings.Join(resolved, ", ")
}

func splitSelectors(s string) []string {
	var result []string
	var depth int
	var inQuote rune
	start := 0
	for i, r := range s {
		if inQuote != 0 {
			if r == inQuote && (i == 0 || s[i-1] != '\\') {
				inQuote = 0
			}
			continue
		}
		if r == '"' || r == '\'' {
			inQuote = r
			continue
		}
		if r == '(' || r == '[' {
			depth++
		} else if r == ')' || r == ']' {
			if depth > 0 {
				depth--
			}
		} else if r == ',' && depth == 0 {
			part := strings.TrimSpace(s[start:i])
			if part != "" {
				result = append(result, part)
			}
			start = i + 1
		}
	}
	if start < len(s) {
		part := strings.TrimSpace(s[start:])
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func hasAmp(selector string) bool {
	var has bool
	for range ampIterator(selector) {
		has = true
		break
	}
	return has
}

func replaceAmp(selector, parent string) string {
	var sb strings.Builder
	var lastPos int
	for pos := range ampIterator(selector) {
		sb.WriteString(selector[lastPos:pos])
		sb.WriteString(parent)
		lastPos = pos + 1
	}
	sb.WriteString(selector[lastPos:])
	return sb.String()
}

func ampIterator(s string) iter.Seq[int] {
	return func(yield func(int) bool) {
		var inQuote rune
		for i, r := range s {
			if inQuote != 0 {
				if r == inQuote && (i == 0 || s[i-1] != '\\') {
					inQuote = 0
				}
				continue
			}
			if r == '"' || r == '\'' {
				inQuote = r
				continue
			}
			if r == '&' && (i == 0 || s[i-1] != '\\') {
				if !yield(i) {
					return
				}
			}
		}
	}
}
