package components

import (
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/parser"
)

var MjAccordionTitleSpec = &core.ComponentSpec{
	TagName:   "mj-accordion-title",
	EndingTag: true,
	AllowedAttributes: map[string]string{
		"background-color": "color",
		"color":            "color",
		"font-size":        "unit(px)",
		"font-family":      "string",
		"font-weight":      "string",
		"padding-bottom":   "unit(px,%)",
		"padding-left":     "unit(px,%)",
		"padding-right":    "unit(px,%)",
		"padding-top":      "unit(px,%)",
		"padding":          "unit(px,%){1,4}",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "font-size", Value: "13px"},
		{Key: "padding", Value: "16px"},
	},
	Create: NewMjAccordionTitle,
}

type MjAccordionTitle struct {
	*core.BodyComponent
}

func NewMjAccordionTitle(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	title := &MjAccordionTitle{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	title.BodyComponent.Component = title

	return title, nil
}

func (t *MjAccordionTitle) GetTagName() string {
	return t.Spec.TagName
}

func (t *MjAccordionTitle) GetStyles(element string) []*core.Style {
	switch element {
	case "td":
		return []*core.Style{
			{Name: "width", Value: "100%"},
			{Name: "background-color", Value: t.GetAttributeOr("background-color", "")},
			{Name: "color", Value: t.GetAttributeOr("color", "")},
			{Name: "font-size", Value: t.GetAttributeOr("font-size", "")},
			{Name: "font-family", Value: t.resolveFontFamily()},
			{Name: "font-weight", Value: t.GetAttributeOr("font-weight", "")},
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
	case "td2":
		return []*core.Style{
			{Name: "padding", Value: "16px"},
			{Name: "background", Value: t.GetAttributeOr("background-color", "")},
			{Name: "vertical-align", Value: t.GetAttributeOr("icon-align", "")},
		}
	case "img":
		return []*core.Style{
			{Name: "display", Value: "none"},
			{Name: "width", Value: t.GetAttributeOr("icon-width", "")},
			{Name: "height", Value: t.GetAttributeOr("icon-height", "")},
		}
	default:
		return nil
	}
}

func (t *MjAccordionTitle) GetChildContext() *core.MJMLContext {
	return t.Context
}

func (t *MjAccordionTitle) renderTitle(w core.MJMLWriter) error {
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

func (t *MjAccordionTitle) renderIcons(w core.MJMLWriter) error {
	if _, err := w.WriteString("<!--[if !mso | IE]><!--><td "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(t).
		Add("class", "mj-accordion-ico").
		Add("style", "td2").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><img "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(t).
		AddNullable("src", t.GetAttribute("icon-wrapped-url")).
		AddNullable("alt", t.GetAttribute("icon-wrapped-alt")).
		Add("class", "mj-accordion-more").
		Add("style", "img").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("/><img "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(t).
		AddNullable("src", t.GetAttribute("icon-unwrapped-url")).
		AddNullable("alt", t.GetAttribute("icon-unwrapped-alt")).
		Add("class", "mj-accordion-less").
		Add("style", "img").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("/></td><!--<![endif]-->"); err != nil {
		return err
	}

	return nil
}

func (t *MjAccordionTitle) resolveFontFamily() string {
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

func (t *MjAccordionTitle) Render(w core.MJMLWriter) error {
	isRight := t.GetAttributeOr("icon-position", "") == "right"

	if _, err := w.WriteString("<div "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(t).
		Add("class", "mj-accordion-title").
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

	if isRight {
		if err := t.renderTitle(w); err != nil {
			return err
		}
		if err := t.renderIcons(w); err != nil {
			return err
		}
	} else {
		if err := t.renderIcons(w); err != nil {
			return err
		}
		if err := t.renderTitle(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("</tr></tbody></table></div>"); err != nil {
		return err
	}

	return nil
}
