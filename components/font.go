package components

import (
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/parser"
)

var MjFontSpec = &core.ComponentSpec{
	TagName: "mj-font",
	AllowedAttributes: map[string]string{
		"name": "string",
		"href": "string",
	},
	Create: NewMjFont,
}

type MjFont struct {
	*core.HeadComponent
}

func NewMjFont(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	font := &MjFont{
		HeadComponent: core.NewHeadComponent(node, spec, props, context),
	}

	font.HeadComponent.Component = font

	return font, nil
}

func (f *MjFont) GetTagName() string {
	return f.Spec.TagName
}

func (f *MjFont) GetStyles(element string) []*core.Style {
	return nil
}

func (f *MjFont) GetChildContext() *core.MJMLContext {
	return f.Context
}

func (f *MjFont) Render(w core.MJMLWriter) error {
	context := f.Context

	name := f.GetAttribute("name")
	href := f.GetAttribute("href")

	if name != nil && href != nil {
		context.AddFont(*name, *href)
	}

	return nil
}
