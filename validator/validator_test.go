package validator

import (
	"testing"

	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/parser"
	"github.com/stretchr/testify/assert"
)

func parse(mjml string) (*parser.MJMLNode, error) {
	return parser.MJMLParser([]byte(mjml), parser.MJMLParserOptions{
		KeepComments: true,
		Components:   []parser.MJMLComponent{},
		FilePath:     ".",
	}, nil)
}

func TestValidateType(t *testing.T) {
	input := `<mj-section background-color="#CCCCCC" full-width="full-width"></mj-section>`
	node, err := parse(input)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, node.TagName, "mj-section")
	assert.Equal(t, node.Attributes.Len(), 2)

	comp := &core.ComponentSpec{
		AllowedAttributes: map[string]string{
			"background-color": "color",
			"full-width":       "enum(auto, full, brief)",
		},
	}

	result := MJMLValidator(node, &core.MJMLOptions{Components: map[string]*core.ComponentSpec{"mj-section": comp}})

	assert.Equal(t, len(result), 1)
	assert.Equal(t, result[0].Message, "Attribute full-width has invalid value: full-width for type Enum, only accepts auto, full, brief")
}

func TestValidateAttribute(t *testing.T) {
	input := `<mj-section mj-class="test" background-color="#CCCCCC" full-width="full" unknown-attr="true"></mj-section>`
	node, err := parse(input)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, node.TagName, "mj-section")
	assert.Equal(t, node.Attributes.Len(), 4)

	comp := &core.ComponentSpec{
		AllowedAttributes: map[string]string{
			"background-color": "color",
			"full-width":       "enum(auto, full, brief)",
		},
	}

	result := MJMLValidator(node, &core.MJMLOptions{Components: map[string]*core.ComponentSpec{"mj-section": comp}})

	assert.Equal(t, len(result), 1)
	assert.Equal(t, result[0].Message, "Attribute unknown-attr is illegal")
}

func TestValidateChildren(t *testing.T) {
	input := `<mj-section><mj-head><mj-text>Hello</mj-text></mj-head></mj-section>`
	node, err := parse(input)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, node.TagName, "mj-section")

	options := &core.MJMLOptions{
		Components: map[string]*core.ComponentSpec{
			"mj-section": {},
			"mj-head":    {},
			"mj-text":    {},
		},
		Dependencies: map[string][]string{
			"mj-section": {"mj-column"},
			"mj-head":    {"mj-text"},
			"mjml":       {"mj-head"},
		},
	}

	result := MJMLValidator(node, options)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, result[0].Message, "mj-head cannot be used inside mj-section, only inside: mjml")
}

func TestValidateAnyChildren(t *testing.T) {
	input := `<mjml><mj-head><mj-attributes><mj-text/><mj-all /></mj-attributes></mj-head><mj-body></mj-body></mjml>`
	node, err := parse(input)
	if err != nil {
		t.Error(err)
	}

	options := &core.MJMLOptions{
		Components: map[string]*core.ComponentSpec{
			"mj-head":       {},
			"mj-attributes": {},
			"mj-body":       {},
			"mj-text":       {},
		},
		Dependencies: map[string][]string{
			"mjml":          {"mj-body", "mj-head"},
			"mj-head":       {"mj-attributes"},
			"mj-body":       {"mj-text"},
			"mj-attributes": {"*"},
		},
	}

	result := MJMLValidator(node, options)

	assert.Equal(t, len(result), 0)
}

func TestValidateTag(t *testing.T) {
	input := `<mj-section><mj-text>Hello</mj-text></mj-section>`
	node, err := parse(input)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, node.TagName, "mj-section")

	options := &core.MJMLOptions{
		Components: map[string]*core.ComponentSpec{
			"mj-section": {},
		},
		Dependencies: map[string][]string{
			"mj-section": {"mj-text"},
		},
	}

	result := MJMLValidator(node, options)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, result[0].Message, "element mj-text doesn't exist or is not registered")
}

func TestValidateParserErrors(t *testing.T) {
	input := `<mjml><mj-body><mj-section><mj-include path="./incl.mjml" /></mj-section></mj-body></mjml>`
	node, err := parse(input)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, node.TagName, "mjml")

	options := &core.MJMLOptions{
		Components: map[string]*core.ComponentSpec{
			"mj-body":    {},
			"mj-section": {},
			"mj-raw":     {},
		},
		Dependencies: map[string][]string{
			"mjml":       {"mj-body"},
			"mj-body":    {"mj-section"},
			"mj-section": {"mj-raw"},
		},
	}

	result := MJMLValidator(node, options)

	assert.Equal(t, len(result), 1)
	assert.Contains(t, result[0].Message, "mj-include fails to read file ")
}
