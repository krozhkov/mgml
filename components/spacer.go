package components

import (
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/parser"
)

var MjSpacerSpec = &core.ComponentSpec{
	TagName: "mj-spacer",
	AllowedAttributes: map[string]string{
		"border":                     "string",
		"border-bottom":              "string",
		"border-left":                "string",
		"border-right":               "string",
		"border-top":                 "string",
		"container-background-color": "color",
		"padding-bottom":             "unit(px,%)",
		"padding-left":               "unit(px,%)",
		"padding-right":              "unit(px,%)",
		"padding-top":                "unit(px,%)",
		"padding":                    "unit(px,%){1,4}",
		"height":                     "unit(px,%)",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "height", Value: "20px"},
	},
	Create: NewMjSpacer,
}

type MjSpacer struct {
	*core.BodyComponent
}

func NewMjSpacer(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	spacer := &MjSpacer{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	spacer.BodyComponent.Component = spacer

	return spacer, nil
}

func (s *MjSpacer) GetTagName() string {
	return s.Spec.TagName
}

func (s *MjSpacer) GetStyles(element string) []*core.Style {
	switch element {
	case "div":
		return []*core.Style{
			{Name: "height", Value: s.GetAttributeOr("height", "")},
			{Name: "line-height", Value: s.GetAttributeOr("height", "")},
		}
	default:
		return nil
	}
}

func (s *MjSpacer) GetChildContext() *core.MJMLContext {
	return s.Context
}

func (s *MjSpacer) Render(w core.MJMLWriter) error {
	if _, err := w.WriteString("<div "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(s).
		Add("style", "div").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">&#8202;</div>\n"); err != nil {
		return err
	}

	return nil
}
