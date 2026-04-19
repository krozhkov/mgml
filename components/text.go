package components

import (
	"fmt"
	"io"

	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/parser"
)

var MjTextSpec = &core.ComponentSpec{
	TagName:   "mj-text",
	EndingTag: true,
	AllowedAttributes: map[string]string{
		"align":                      "enum(left,right,center,justify)",
		"background-color":           "color",
		"color":                      "color",
		"container-background-color": "color",
		"font-family":                "string",
		"font-size":                  "unit(px)",
		"font-style":                 "string",
		"font-weight":                "string",
		"height":                     "unit(px,%)",
		"letter-spacing":             "unitWithNegative(px,em)",
		"line-height":                "unit(px,%,)",
		"padding-bottom":             "unit(px,%)",
		"padding-left":               "unit(px,%)",
		"padding-right":              "unit(px,%)",
		"padding-top":                "unit(px,%)",
		"padding":                    "unit(px,%){1,4}",
		"text-decoration":            "string",
		"text-transform":             "string",
		"vertical-align":             "enum(top,bottom,middle)",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "align", Value: "left"},
		{Key: "color", Value: "#000000"},
		{Key: "font-family", Value: "Ubuntu, Helvetica, Arial, sans-serif"},
		{Key: "font-size", Value: "13px"},
		{Key: "line-height", Value: "1"},
		{Key: "padding", Value: "10px 25px"},
	},
	Create: NewMjText,
}

type MjText struct {
	*core.BodyComponent
}

func NewMjText(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	text := &MjText{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	text.BodyComponent.Component = text

	return text, nil
}

func (t *MjText) GetTagName() string {
	return t.Spec.TagName
}

func (t *MjText) GetStyles(element string) []*core.Style {
	switch element {
	case "text":
		return []*core.Style{
			{Name: "font-family", Value: t.GetAttributeOr("font-family", "")},
			{Name: "font-size", Value: t.GetAttributeOr("font-size", "")},
			{Name: "font-style", Value: t.GetAttributeOr("font-style", "")},
			{Name: "font-weight", Value: t.GetAttributeOr("font-weight", "")},
			{Name: "letter-spacing", Value: t.GetAttributeOr("letter-spacing", "")},
			{Name: "line-height", Value: t.GetAttributeOr("line-height", "")},
			{Name: "text-align", Value: t.GetAttributeOr("align", "")},
			{Name: "text-decoration", Value: t.GetAttributeOr("text-decoration", "")},
			{Name: "text-transform", Value: t.GetAttributeOr("text-transform", "")},
			{Name: "color", Value: t.GetAttributeOr("color", "")},
			{Name: "height", Value: t.GetAttributeOr("height", "")},
		}
	default:
		return nil
	}
}

func (t *MjText) GetChildContext() *core.MJMLContext {
	return t.Context
}

func (t *MjText) renderContent(w io.StringWriter) error {
	if _, err := w.WriteString("<div "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(t).
		Add("style", "text").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if _, err := w.WriteString(t.GetContent()); err != nil {
		return err
	}

	if _, err := w.WriteString("</div>"); err != nil {
		return err
	}

	return nil
}

func (t *MjText) Render(w core.MJMLWriter) error {
	height := t.GetAttribute("height")

	if height != nil && *height != "" {
		if _, err := w.WriteConditionalTag(fmt.Sprintf("<table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\"><tr><td height=\"%s\" style=\"vertical-align:top;height:%s;\">", *height, *height), false); err != nil {
			return err
		}

		if err := t.renderContent(w); err != nil {
			return err
		}

		if _, err := w.WriteConditionalTag("</td></tr></table>", false); err != nil {
			return err
		}

		return nil
	} else {
		return t.renderContent(w)
	}
}
