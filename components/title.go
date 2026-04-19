package components

import (
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/parser"
)

var MjTitleSpec = &core.ComponentSpec{
	TagName:   "mj-title",
	EndingTag: true,
	Create:    NewMjTitle,
}

type MjTitle struct {
	*core.HeadComponent
}

func NewMjTitle(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	title := &MjTitle{
		HeadComponent: core.NewHeadComponent(node, spec, props, context),
	}

	title.HeadComponent.Component = title

	return title, nil
}

func (b *MjTitle) GetTagName() string {
	return b.Spec.TagName
}

func (b *MjTitle) GetStyles(element string) []*core.Style {
	return nil
}

func (b *MjTitle) GetChildContext() *core.MJMLContext {
	return b.Context
}

func (b *MjTitle) Render(w core.MJMLWriter) error {
	context := b.Context

	content := b.GetContent()

	context.GlobalData.Title = &content

	return nil
}
