package types

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/krozhkov/mgml/internal/utils"
)

var unitMatcher = regexp.MustCompile(`(?im)^(unit|unitWithNegative)\(.*\)`)
var unitConfigMatcher = regexp.MustCompile(`\(([^)]+)\)`)
var unitArgsMather = regexp.MustCompile(`\{([^}]+)\}`)
var allowNegRegexp = regexp.MustCompile(`^unitWithNegative`)

type Unit struct {
	value string
	units []string
	args  []string
	regex *regexp.Regexp
}

func NewUnitFactory(config string) (TypeConstructor, error) {
	var allowNeg = ""
	if allowNegRegexp.MatchString(config) {
		allowNeg = "-|"
	}

	match := unitConfigMatcher.FindStringSubmatch(config)
	if len(match) < 2 {
		return nil, fmt.Errorf("unit type is invalid %s", config)
	}

	units := utils.MapFunc(strings.Split(match[1], ","), func(s string) string { return strings.TrimSpace(s) })

	match = unitArgsMather.FindStringSubmatch(config)
	var args []string
	if len(match) > 1 {
		args = utils.MapFunc(strings.Split(match[1], ","), func(s string) string { return strings.TrimSpace(s) })
	} else {
		args = []string{"1"}
	}

	var allowAuto = ""
	if slices.Index(units, "auto") != -1 {
		allowAuto = "|auto"
	}

	filteredUnits := utils.FilterFunc(units, func(v string) bool { return v != "auto" })

	regex, err := regexp.Compile(`^(((` + allowNeg + `\d|,|\.){1,}(` + strings.Join(filteredUnits, "|") + `)|0` + allowAuto + `)( )?){` + strings.Join(args, ",") + `}$`)
	if err != nil {
		return nil, err
	}

	return func(value string) Type {
		return NewUnit(value, units, args, regex)
	}, nil
}

func NewUnit(value string, units []string, args []string, regex *regexp.Regexp) Type {
	return &Unit{value, units, args, regex}
}

func (t *Unit) Name() string {
	return "Unit"
}

func (t *Unit) Value() string {
	return t.value
}

func (t *Unit) IsValid() bool {
	return t.regex.MatchString(t.value)
}

func (t *Unit) GetError() error {
	if t.IsValid() {
		return nil
	}

	return fmt.Errorf("has invalid value: %s for type Unit, only accepts (%s) units and %s value(s)", t.value, strings.Join(t.units, ", "), strings.Join(t.args, " to "))
}
