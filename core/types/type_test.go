package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func formatType(t Type) string {
	return fmt.Sprintf("Type: %s — Value: %s — isValid: %t %v", t.Name(), t.Value(), t.IsValid(), t.GetError())
}

func TestColorType(t *testing.T) {
	colortype, err := InitializeType("color")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Type: Color — Value: grey — isValid: true <nil>", formatType(colortype("grey")))
	assert.Equal(t, "Type: Color — Value: rgba(0,255,3,0.3) — isValid: true <nil>", formatType(colortype("rgba(0,255,3,0.3)")))
	assert.Equal(t, "Type: Color — Value: #DDDDFF — isValid: true <nil>", formatType(colortype("#DDF")))
	assert.Equal(t, "Type: Color — Value: #DF — isValid: false has invalid value: #DF for type Color", formatType(colortype("#DF")))
}

func TestBooleanType(t *testing.T) {
	booleantype, err := InitializeType("boolean")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Type: Boolean — Value: true — isValid: true <nil>", formatType(booleantype("true")))
	assert.Equal(t, "Type: Boolean — Value: false — isValid: true <nil>", formatType(booleantype("false")))
	assert.Equal(t, "Type: Boolean — Value: banana — isValid: false has invalid value: banana for type Boolean", formatType(booleantype("banana")))
}

func TestIntegerType(t *testing.T) {
	stringtype, err := InitializeType("integer")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Type: NInteger — Value: 1280 — isValid: true <nil>", formatType(stringtype("1280")))
}

func TestStringType(t *testing.T) {
	stringtype, err := InitializeType("string")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Type: NString — Value: hello world — isValid: true <nil>", formatType(stringtype("hello world")))
}

func TestEnumType(t *testing.T) {
	enumtype, err := InitializeType("enum(top,left,center)")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Type: Enum — Value: left — isValid: true <nil>", formatType(enumtype("left")))
	assert.Equal(t, "Type: Enum — Value: bottom — isValid: false has invalid value: bottom for type Enum, only accepts top, left, center", formatType(enumtype("bottom")))
}

func TestUnitType(t *testing.T) {
	unittype, err := InitializeType("unit(px,%){1,3}")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Type: Unit — Value: 10 20px 20 — isValid: false has invalid value: 10 20px 20 for type Unit, only accepts (px, %) units and 1 to 3 value(s)", formatType(unittype("10 20px 20")))
	assert.Equal(t, "Type: Unit — Value: 10px 20px 20px — isValid: true <nil>", formatType(unittype("10px 20px 20px")))
	assert.Equal(t, "Type: Unit — Value: 10px — isValid: true <nil>", formatType(unittype("10px")))
	assert.Equal(t, "Type: Unit — Value: 10% — isValid: true <nil>", formatType(unittype("10%")))
	assert.Equal(t, "Type: Unit — Value: 10px 10px — isValid: true <nil>", formatType(unittype("10px 10px")))
	assert.Equal(t, "Type: Unit — Value: 0 — isValid: true <nil>", formatType(unittype("0")))
}
