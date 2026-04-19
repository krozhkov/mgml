package components

import (
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/parser"
)

var MjHeadSpec = &core.ComponentSpec{
	TagName: "mj-head",
	Create:  NewMjHead,
}

type MjHead struct {
	*core.HeadComponent
}

func NewMjHead(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	head := &MjHead{
		HeadComponent: core.NewHeadComponent(node, spec, props, context),
	}

	head.HeadComponent.Component = head

	children, err := head.CreateChildren()
	if err != nil {
		return head, err
	}

	head.Children = children

	return head, nil
}

func (h *MjHead) GetTagName() string {
	return h.Spec.TagName
}

func (h *MjHead) GetStyles(element string) []*core.Style {
	return nil
}

func (h *MjHead) GetChildContext() *core.MJMLContext {
	return h.Context
}

func (h *MjHead) Render(w core.MJMLWriter) error {
	for _, child := range h.Children {
		child.Render(w)
	}

	return nil
}
