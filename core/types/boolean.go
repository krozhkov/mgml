package types

import (
	"fmt"
	"regexp"
	"strings"
)

var booleanMatcher = regexp.MustCompile(`(?im)^boolean`)

type Boolean struct {
	value string
}

func NewBooleanFactory(_ string) (TypeConstructor, error) {
	return NewBoolean, nil
}

func NewBoolean(value string) Type {
	return &Boolean{value}
}

func (t *Boolean) Name() string {
	return "Boolean"
}

func (t *Boolean) Value() string {
	return t.value
}

func (t *Boolean) IsValid() bool {
	return strings.EqualFold(t.value, "true") || strings.EqualFold(t.value, "false")
}

func (t *Boolean) GetError() error {
	if t.IsValid() {
		return nil
	}

	return fmt.Errorf("has invalid value: %s for type %s", t.value, t.Name())
}
