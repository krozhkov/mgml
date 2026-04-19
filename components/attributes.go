package components

import (
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/parser"
)

var MjAttributesSpec = &core.ComponentSpec{
	TagName: "mj-attributes",
	Create:  NewMjAttributes,
}

type MjAttributes struct {
	*core.HeadComponent
}

func NewMjAttributes(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	attributes := &MjAttributes{
		HeadComponent: core.NewHeadComponent(node, spec, props, context),
	}

	attributes.HeadComponent.Component = attributes

	return attributes, nil
}

func (a *MjAttributes) GetTagName() string {
	return a.Spec.TagName
}

func (a *MjAttributes) GetStyles(element string) []*core.Style {
	return nil
}

func (a *MjAttributes) GetChildContext() *core.MJMLContext {
	return a.Context
}

func (a *MjAttributes) Render(w core.MJMLWriter) error {
	context := a.Context

	for _, child := range a.Node.Children {
		if child.Attributes == nil {
			continue
		}

		if child.TagName == "mj-class" {
			if child.Attributes.Has("name") {
				name, _ := child.Attributes.Get("name")
				copy := child.Attributes.Copy()
				copy.Delete("name")
				context.GlobalData.Classes[name] = copy
			}
		} else {
			context.GlobalData.DefaultAttributes[child.TagName] = child.Attributes
		}
	}

	return nil
}
