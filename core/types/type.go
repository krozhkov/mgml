package types

import (
	"fmt"
)

type Type interface {
	Name() string
	Value() string
	IsValid() bool
	GetError() error
}

type TypeConstructor func(value string) Type

type TypeFactory func(typeConfig string) (TypeConstructor, error)

var existingTypes = make(map[string]TypeConstructor)

func InitializeType(typeConfig string) (TypeConstructor, error) {
	existing, ok := existingTypes[typeConfig]
	if ok {
		return existing, nil
	}

	var typeConstructor TypeConstructor
	var err error
	for _, item := range constructors {
		if item.matcher.MatchString(typeConfig) {
			typeConstructor, err = item.factory(typeConfig)
			if err != nil {
				return nil, err
			}

			break
		}
	}

	if typeConstructor == nil {
		return nil, fmt.Errorf("no type found for %s", typeConfig)
	}

	existingTypes[typeConfig] = typeConstructor

	return typeConstructor, nil
}
