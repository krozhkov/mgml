package components

import (
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/parser"
)

var MjBreakpointSpec = &core.ComponentSpec{
	TagName:   "mj-breakpoint",
	EndingTag: true,
	AllowedAttributes: map[string]string{
		"width": "unit(px)",
	},
	Create: NewMjBreakpoint,
}

type MjBreakpoint struct {
	*core.HeadComponent
}

func NewMjBreakpoint(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	breakpoint := &MjBreakpoint{
		HeadComponent: core.NewHeadComponent(node, spec, props, context),
	}

	breakpoint.HeadComponent.Component = breakpoint

	return breakpoint, nil
}

func (b *MjBreakpoint) GetTagName() string {
	return b.Spec.TagName
}

func (b *MjBreakpoint) GetStyles(element string) []*core.Style {
	return nil
}

func (b *MjBreakpoint) GetChildContext() *core.MJMLContext {
	return b.Context
}

func (b *MjBreakpoint) Render(w core.MJMLWriter) error {
	context := b.Context

	width := b.GetAttributeOr("width", "")

	context.GlobalData.Breakpoint = width

	return nil
}
