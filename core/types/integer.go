package types

import (
	"fmt"
	"regexp"
)

var integerMatcher = regexp.MustCompile(`(?im)^integer`)
var integerValueMatcher = regexp.MustCompile(`\d+`)

type NInteger struct {
	value string
}

func NewNIntegerFactory(_ string) (TypeConstructor, error) {
	return NewNInteger, nil
}

func NewNInteger(value string) Type {
	return &NInteger{value}
}

func (t *NInteger) Name() string {
	return "NInteger"
}

func (t *NInteger) Value() string {
	return t.value
}

func (t *NInteger) IsValid() bool {
	return integerValueMatcher.MatchString(t.value)
}

func (t *NInteger) GetError() error {
	if t.IsValid() {
		return nil
	}

	return fmt.Errorf("has invalid value: %s for type %s", t.value, t.Name())
}
