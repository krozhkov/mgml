package rules

import "github.com/krozhkov/mgml/core"

type RuleOptions struct {
	Components   map[string]*core.ComponentSpec
	Dependencies map[string]map[string]struct{}
	SkipElements map[string]struct{}
}
