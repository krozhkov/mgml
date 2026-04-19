package types

import (
	"fmt"
	"regexp"
)

func isColor(name string) bool {
	switch name {
	case "aliceblue",
		"antiquewhite",
		"aqua",
		"aquamarine",
		"azure",
		"beige",
		"bisque",
		"black",
		"blanchedalmond",
		"blue",
		"blueviolet",
		"brown",
		"burlywood",
		"cadetblue",
		"chartreuse",
		"chocolate",
		"coral",
		"cornflowerblue",
		"cornsilk",
		"crimson",
		"cyan",
		"darkblue",
		"darkcyan",
		"darkgoldenrod",
		"darkgray",
		"darkgreen",
		"darkgrey",
		"darkkhaki",
		"darkmagenta",
		"darkolivegreen",
		"darkorange",
		"darkorchid",
		"darkred",
		"darksalmon",
		"darkseagreen",
		"darkslateblue",
		"darkslategray",
		"darkslategrey",
		"darkturquoise",
		"darkviolet",
		"deeppink",
		"deepskyblue",
		"dimgray",
		"dimgrey",
		"dodgerblue",
		"firebrick",
		"floralwhite",
		"forestgreen",
		"fuchsia",
		"gainsboro",
		"ghostwhite",
		"gold",
		"goldenrod",
		"gray",
		"green",
		"greenyellow",
		"grey",
		"honeydew",
		"hotpink",
		"indianred",
		"indigo",
		"inherit",
		"ivory",
		"khaki",
		"lavender",
		"lavenderblush",
		"lawngreen",
		"lemonchiffon",
		"lightblue",
		"lightcoral",
		"lightcyan",
		"lightgoldenrodyellow",
		"lightgray",
		"lightgreen",
		"lightgrey",
		"lightpink",
		"lightsalmon",
		"lightseagreen",
		"lightskyblue",
		"lightslategray",
		"lightslategrey",
		"lightsteelblue",
		"lightyellow",
		"lime",
		"limegreen",
		"linen",
		"magenta",
		"maroon",
		"mediumaquamarine",
		"mediumblue",
		"mediumorchid",
		"mediumpurple",
		"mediumseagreen",
		"mediumslateblue",
		"mediumspringgreen",
		"mediumturquoise",
		"mediumvioletred",
		"midnightblue",
		"mintcream",
		"mistyrose",
		"moccasin",
		"navajowhite",
		"navy",
		"oldlace",
		"olive",
		"olivedrab",
		"orange",
		"orangered",
		"orchid",
		"palegoldenrod",
		"palegreen",
		"paleturquoise",
		"palevioletred",
		"papayawhip",
		"peachpuff",
		"peru",
		"pink",
		"plum",
		"powderblue",
		"purple",
		"rebeccapurple",
		"red",
		"rosybrown",
		"royalblue",
		"saddlebrown",
		"salmon",
		"sandybrown",
		"seagreen",
		"seashell",
		"sienna",
		"silver",
		"skyblue",
		"slateblue",
		"slategray",
		"slategrey",
		"snow",
		"springgreen",
		"steelblue",
		"tan",
		"teal",
		"thistle",
		"tomato",
		"transparent",
		"turquoise",
		"violet",
		"wheat",
		"white",
		"whitesmoke",
		"yellow",
		"yellowgreen":
		return true
	default:
		return false
	}
}

var colorMatcher = regexp.MustCompile(`(?im)^color`)
var colorRgbaValueMatcher = regexp.MustCompile(`(?i)rgba\(\d{1,3},\s?\d{1,3},\s?\d{1,3},\s?\d(\.\d{1,3})?\)`)
var colorRgbValueMatcher = regexp.MustCompile(`(?i)rgb\(\d{1,3},\s?\d{1,3},\s?\d{1,3}\)`)
var colorHexValueMatcher = regexp.MustCompile(`(?i)^#([0-9a-f]{3}){1,2}$`)
var shorthandRegex = regexp.MustCompile(`^#\w{3}$`)
var replaceInputRegex = regexp.MustCompile(`^#(\w)(\w)(\w)$`)
var replaceOutput = "#$1$1$2$2$3$3"

type Color struct {
	value string
}

func NewColorFactory(_ string) (TypeConstructor, error) {
	return NewColor, nil
}

func NewColor(value string) Type {
	return &Color{value}
}

func (t *Color) Name() string {
	return "Color"
}

func (t *Color) Value() string {
	if shorthandRegex.MatchString(t.value) {
		return replaceInputRegex.ReplaceAllString(t.value, replaceOutput)
	}

	return t.value
}

func (t *Color) IsValid() bool {
	if isColor(t.value) {
		return true
	}

	if colorRgbaValueMatcher.MatchString(t.value) || colorRgbValueMatcher.MatchString(t.value) || colorHexValueMatcher.MatchString(t.value) {
		return true
	}

	return false
}

func (t *Color) GetError() error {
	if t.IsValid() {
		return nil
	}

	return fmt.Errorf("has invalid value: %s for type %s", t.value, t.Name())
}
