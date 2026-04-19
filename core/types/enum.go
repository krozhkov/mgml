package types

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/krozhkov/mgml/internal/utils"
)

var enumMatcher = regexp.MustCompile(`(?im)^enum`)
var enumConfigMatcher = regexp.MustCompile(`\(([^)]+)\)`)

type Enum struct {
	value    string
	matchers []string
	regex    *regexp.Regexp
}

func NewEnumFactory(config string) (TypeConstructor, error) {
	match := enumConfigMatcher.FindStringSubmatch(config)
	if len(match) < 2 {
		return nil, fmt.Errorf("enum type is invalid %s", config)
	}

	matchers := utils.MapFunc(strings.Split(match[1], ","), func(s string) string { return strings.TrimSpace(s) })
	regex, err := regexp.Compile(`^(` + strings.Join(matchers, "|") + `)$`)
	if err != nil {
		return nil, err
	}

	return func(value string) Type {
		return NewEnum(value, matchers, regex)
	}, nil
}

func NewEnum(value string, matchers []string, regex *regexp.Regexp) Type {
	return &Enum{value, matchers, regex}
}

func (t *Enum) Name() string {
	return "Enum"
}

func (t *Enum) Value() string {
	return t.value
}

func (t *Enum) IsValid() bool {
	return t.regex.MatchString(t.value)
}

func (t *Enum) GetError() error {
	if t.IsValid() {
		return nil
	}

	return fmt.Errorf("has invalid value: %s for type Enum, only accepts %s", t.value, strings.Join(t.matchers, ", "))
}
