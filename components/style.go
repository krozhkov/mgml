package components

import (
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/parser"
)

var MjStyleSpec = &core.ComponentSpec{
	TagName: "mj-style",
	AllowedAttributes: map[string]string{
		"inline": "string",
	},
	Create: NewMjStyle,
}

type MjStyle struct {
	*core.HeadComponent
}

func NewMjStyle(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	style := &MjStyle{
		HeadComponent: core.NewHeadComponent(node, spec, props, context),
	}

	style.HeadComponent.Component = style

	return style, nil
}

func (s *MjStyle) GetTagName() string {
	return s.Spec.TagName
}

func (s *MjStyle) GetStyles(element string) []*core.Style {
	return nil
}

func (s *MjStyle) GetChildContext() *core.MJMLContext {
	return s.Context
}

func (s *MjStyle) Render(w core.MJMLWriter) error {
	context := s.Context

	inline := s.GetAttributeOr("inline", "") == "inline"
	content := s.GetContent()

	context.AddStyle(content, inline)

	return nil
}
