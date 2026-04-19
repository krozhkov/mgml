package components

import (
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/parser"
)

var MjRawSpec = &core.ComponentSpec{
	TagName:      "mj-raw",
	EndingTag:    true,
	IsRawElement: true,
	AllowedAttributes: map[string]string{
		"position": "enum(file-start)",
	},
	Create: NewMjRaw,
}

type MjRaw struct {
	*core.BodyComponent
}

func NewMjRaw(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	raw := &MjRaw{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	raw.BodyComponent.Component = raw

	return raw, nil
}

func (r *MjRaw) GetTagName() string {
	return r.Spec.TagName
}

func (r *MjRaw) GetStyles(element string) []*core.Style {
	return nil
}

func (r *MjRaw) GetChildContext() *core.MJMLContext {
	return r.Context
}

func (r *MjRaw) Render(w core.MJMLWriter) error {
	if _, err := w.WriteString(r.GetContent()); err != nil {
		return err
	}

	return nil
}
