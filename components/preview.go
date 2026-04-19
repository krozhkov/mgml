package components

import (
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/parser"
)

var MjPreviewSpec = &core.ComponentSpec{
	TagName:   "mj-preview",
	EndingTag: true,
	Create:    NewMjPreview,
}

type MjPreview struct {
	*core.HeadComponent
}

func NewMjPreview(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	preview := &MjPreview{
		HeadComponent: core.NewHeadComponent(node, spec, props, context),
	}

	preview.HeadComponent.Component = preview

	return preview, nil
}

func (p *MjPreview) GetTagName() string {
	return p.Spec.TagName
}

func (p *MjPreview) GetStyles(element string) []*core.Style {
	return nil
}

func (p *MjPreview) GetChildContext() *core.MJMLContext {
	return p.Context
}

func (p *MjPreview) Render(w core.MJMLWriter) error {
	context := p.Context

	context.AddPreview(p.GetContent())

	return nil
}
