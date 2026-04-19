package types

import "regexp"

type typeConstructors struct {
	matcher *regexp.Regexp
	factory TypeFactory
}

var constructors = []typeConstructors{
	{matcher: booleanMatcher, factory: NewBooleanFactory},
	{matcher: colorMatcher, factory: NewColorFactory},
	{matcher: stringMatcher, factory: NewNStringFactory},
	{matcher: integerMatcher, factory: NewNIntegerFactory},
	{matcher: enumMatcher, factory: NewEnumFactory},
	{matcher: unitMatcher, factory: NewUnitFactory},
}
