package components

import (
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/core/helpers"
	"github.com/krozhkov/mgml/parser"
)

var MjBodySpec = &core.ComponentSpec{
	TagName: "mj-body",
	AllowedAttributes: map[string]string{
		"background-color": "color",
		"width":            "unit(px)",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "width", Value: "600px"},
	},
	Create: NewMjBody,
}

type MjBody struct {
	*core.BodyComponent
}

func NewMjBody(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	body := &MjBody{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	body.BodyComponent.Component = body

	width := body.GetAttribute("width")
	if width != nil && *width != "" {
		value, _ := helpers.WidthParser(*width, true)
		context.ContainerWidth = int(value)
	}

	children, err := body.CreateChildren(nil, &core.CreateChildrenOptions{})
	if err != nil {
		return body, err
	}

	body.Children = children

	return body, nil
}

func (b *MjBody) GetTagName() string {
	return b.Spec.TagName
}

func (b *MjBody) GetStyles(element string) []*core.Style {
	switch element {
	case "div":
		return []*core.Style{
			{Name: "background-color", Value: b.GetAttributeOr("background-color", "")},
		}
	default:
		return nil
	}
}

func (b *MjBody) GetChildContext() *core.MJMLContext {
	return b.Context
}

func (b *MjBody) Render(w core.MJMLWriter) error {
	context := b.Context

	bgColor := b.GetAttribute("background-color")
	if bgColor != nil && *bgColor != "" {
		context.GlobalData.BackgroundColor = bgColor
	}

	if _, err := w.WriteString("<div "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(b).
		AddNullable("aria-label", context.GlobalData.Title).
		Add("aria-roledescription", "email").
		AddNullable("class", b.GetAttribute("css-class")).
		Add("style", "div").
		Add("role", "article").
		Add("lang", context.GlobalData.Lang).
		Add("dir", context.GlobalData.Dir).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	for _, child := range b.Children {
		child.Render(w)
	}

	if _, err := w.WriteString("</div>"); err != nil {
		return err
	}

	return nil
}
