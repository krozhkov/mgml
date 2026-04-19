package types

import (
	"fmt"
	"regexp"
)

var stringMatcher = regexp.MustCompile(`(?im)^string`)
var stringValueMatcher = regexp.MustCompile(`.*`)

type NString struct {
	value string
}

func NewNStringFactory(_ string) (TypeConstructor, error) {
	return NewNString, nil
}

func NewNString(value string) Type {
	return &NString{value}
}

func (t *NString) Name() string {
	return "NString"
}

func (t *NString) Value() string {
	return t.value
}

func (t *NString) IsValid() bool {
	return stringValueMatcher.MatchString(t.value)
}

func (t *NString) GetError() error {
	if t.IsValid() {
		return nil
	}

	return fmt.Errorf("has invalid value: %s for type %s", t.value, t.Name())
}
