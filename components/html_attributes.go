package components

import (
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/internal/utils"
	"github.com/krozhkov/mgml/parser"
)

var MjHtmlAttributesSpec = &core.ComponentSpec{
	TagName: "mj-html-attributes",
	Create:  NewMjHtmlAttributes,
}

type MjHtmlAttributes struct {
	*core.HeadComponent
}

func NewMjHtmlAttributes(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	attrs := &MjHtmlAttributes{
		HeadComponent: core.NewHeadComponent(node, spec, props, context),
	}

	attrs.HeadComponent.Component = attrs

	return attrs, nil
}

func (h *MjHtmlAttributes) GetTagName() string {
	return h.Spec.TagName
}

func (h *MjHtmlAttributes) GetStyles(element string) []*core.Style {
	return nil
}

func (h *MjHtmlAttributes) GetChildContext() *core.MJMLContext {
	return h.Context
}

func (h *MjHtmlAttributes) Render(w core.MJMLWriter) error {
	context := h.Context

	children := utils.FilterFunc(h.Node.Children, func(n *parser.MJMLNode) bool {
		return n.TagName == "mj-selector"
	})

	for _, child := range children {
		path, ok := child.Attributes.Get("path")

		if !ok {
			continue
		}

		attributes := utils.MapFunc(
			utils.FilterFunc(child.Children, func(n *parser.MJMLNode) bool {
				return n.TagName == "mj-html-attribute" && n.Attributes.GetOrDefault("name", "") != ""
			}),
			func(n *parser.MJMLNode) core.HtmlAttribute {
				return core.HtmlAttribute{Name: n.Attributes.GetOrDefault("name", ""), Value: n.Content}
			},
		)

		context.AddHtmlAttributes(path, attributes)
	}

	return nil
}
