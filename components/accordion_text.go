package components

import (
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/parser"
)

var MjAccordionTextSpec = &core.ComponentSpec{
	TagName:   "mj-accordion-text",
	EndingTag: true,
	AllowedAttributes: map[string]string{
		"background-color": "color",
		"font-size":        "unit(px)",
		"font-family":      "string",
		"font-weight":      "string",
		"letter-spacing":   "unitWithNegative(px,em)",
		"line-height":      "unit(px,%,)",
		"color":            "color",
		"padding-bottom":   "unit(px,%)",
		"padding-left":     "unit(px,%)",
		"padding-right":    "unit(px,%)",
		"padding-top":      "unit(px,%)",
		"padding":          "unit(px,%){1,4}",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "font-size", Value: "13px"},
		{Key: "line-height", Value: "1"},
		{Key: "padding", Value: "16px"},
	},
	Create: NewMjAccordionText,
}

type MjAccordionText struct {
	*core.BodyComponent
}

func NewMjAccordionText(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	text := &MjAccordionText{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	text.BodyComponent.Component = text

	return text, nil
}

func (t *MjAccordionText) GetTagName() string {
	return t.Spec.TagName
}

func (t *MjAccordionText) GetStyles(element string) []*core.Style {
	switch element {
	case "td":
		return []*core.Style{
			{Name: "background", Value: t.GetAttributeOr("background-color", "")},
			{Name: "font-size", Value: t.GetAttributeOr("font-size", "")},
			{Name: "font-family", Value: t.resolveFontFamily()},
			{Name: "font-weight", Value: t.GetAttributeOr("font-weight", "")},
			{Name: "letter-spacing", Value: t.GetAttributeOr("letter-spacing", "")},
			{Name: "line-height", Value: t.GetAttributeOr("line-height", "")},
			{Name: "color", Value: t.GetAttributeOr("color", "")},
			{Name: "padding", Value: t.GetAttributeOr("padding", "")},
			{Name: "padding-bottom", Value: t.GetAttributeOr("padding-bottom", "")},
			{Name: "padding-left", Value: t.GetAttributeOr("padding-left", "")},
			{Name: "padding-right", Value: t.GetAttributeOr("padding-right", "")},
			{Name: "padding-top", Value: t.GetAttributeOr("padding-top", "")},
		}
	case "table":
		return []*core.Style{
			{Name: "width", Value: "100%"},
			{Name: "border-bottom", Value: t.GetAttributeOr("border", "")},
		}
	default:
		return nil
	}
}

func (t *MjAccordionText) GetChildContext() *core.MJMLContext {
	return t.Context
}

func (t *MjAccordionText) renderContent(w core.MJMLWriter) error {
	if _, err := w.WriteString("<td "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(t).
		AddNullable("class", t.GetAttribute("css-class")).
		Add("style", "td").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if _, err := w.WriteString(t.GetContent()); err != nil {
		return err
	}

	if _, err := w.WriteString("</td>"); err != nil {
		return err
	}

	return nil
}

func (t *MjAccordionText) resolveFontFamily() string {
	if t.Node != nil && t.Node.Attributes != nil {
		fontFamily, ok := t.Node.Attributes.Get("font-family")
		if ok {
			return fontFamily
		}
	}

	elementFontFamily := t.GetAttribute("elementFontFamily")
	if elementFontFamily != nil {
		return *elementFontFamily
	}

	accordionFontFamily := t.GetAttribute("accordionFontFamily")
	if accordionFontFamily != nil {
		return *accordionFontFamily
	}

	fontFamily := t.GetAttribute("font-family")
	if fontFamily != nil {
		return *fontFamily
	}

	return ""
}

func (t *MjAccordionText) Render(w core.MJMLWriter) error {
	if _, err := w.WriteString("<div "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(t).
		Add("class", "mj-accordion-content").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><table "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(t).
		Add("cellspacing", "0").
		Add("cellpadding", "0").
		Add("style", "table").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><tbody><tr>"); err != nil {
		return err
	}

	if err := t.renderContent(w); err != nil {
		return err
	}

	if _, err := w.WriteString("</tr></tbody></table></div>"); err != nil {
		return err
	}

	return nil
}
