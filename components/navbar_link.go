package components

import (
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/core/helpers"
	"github.com/krozhkov/mgml/parser"
)

var MjNavbarLinkSpec = &core.ComponentSpec{
	TagName:   "mj-navbar-link",
	EndingTag: true,
	AllowedAttributes: map[string]string{
		"color":           "color",
		"font-family":     "string",
		"font-size":       "unit(px)",
		"font-style":      "string",
		"font-weight":     "string",
		"href":            "string",
		"name":            "string",
		"target":          "string",
		"rel":             "string",
		"letter-spacing":  "unitWithNegative(px,em)",
		"line-height":     "unit(px,%,)",
		"padding-bottom":  "unit(px,%)",
		"padding-left":    "unit(px,%)",
		"padding-right":   "unit(px,%)",
		"padding-top":     "unit(px,%)",
		"padding":         "unit(px,%){1,4}",
		"text-decoration": "string",
		"text-transform":  "string",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "color", Value: "#000000"},
		{Key: "font-family", Value: "Ubuntu, Helvetica, Arial, sans-serif"},
		{Key: "font-size", Value: "13px"},
		{Key: "font-weight", Value: "normal"},
		{Key: "line-height", Value: "22px"},
		{Key: "padding", Value: "15px 10px"},
		{Key: "target", Value: "_blank"},
		{Key: "text-decoration", Value: "none"},
		{Key: "text-transform", Value: "uppercase"},
	},
	Create: NewMjNavbarLink,
}

type MjNavbarLink struct {
	*core.BodyComponent
}

func NewMjNavbarLink(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	link := &MjNavbarLink{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	link.BodyComponent.Component = link

	return link, nil
}

func (l *MjNavbarLink) GetTagName() string {
	return l.Spec.TagName
}

func (l *MjNavbarLink) GetStyles(element string) []*core.Style {
	switch element {
	case "a":
		return []*core.Style{
			{Name: "display", Value: "inline-block"},
			{Name: "color", Value: l.GetAttributeOr("color", "")},
			{Name: "font-family", Value: l.GetAttributeOr("font-family", "")},
			{Name: "font-size", Value: l.GetAttributeOr("font-size", "")},
			{Name: "font-style", Value: l.GetAttributeOr("font-style", "")},
			{Name: "font-weight", Value: l.GetAttributeOr("font-weight", "")},
			{Name: "letter-spacing", Value: l.GetAttributeOr("letter-spacing", "")},
			{Name: "line-height", Value: l.GetAttributeOr("line-height", "")},
			{Name: "text-decoration", Value: l.GetAttributeOr("text-decoration", "")},
			{Name: "text-transform", Value: l.GetAttributeOr("text-transform", "")},
			{Name: "padding", Value: l.GetAttributeOr("padding", "")},
			{Name: "padding-top", Value: l.GetAttributeOr("padding-top", "")},
			{Name: "padding-left", Value: l.GetAttributeOr("padding-left", "")},
			{Name: "padding-right", Value: l.GetAttributeOr("padding-right", "")},
			{Name: "padding-bottom", Value: l.GetAttributeOr("padding-bottom", "")},
		}
	case "td":
		return []*core.Style{
			{Name: "padding", Value: l.GetAttributeOr("padding", "")},
			{Name: "padding-top", Value: l.GetAttributeOr("padding-top", "")},
			{Name: "padding-left", Value: l.GetAttributeOr("padding-left", "")},
			{Name: "padding-right", Value: l.GetAttributeOr("padding-right", "")},
			{Name: "padding-bottom", Value: l.GetAttributeOr("padding-bottom", "")},
		}
	default:
		return nil
	}
}

func (l *MjNavbarLink) GetChildContext() *core.MJMLContext {
	return l.Context
}

func (l *MjNavbarLink) renderContent(w core.MJMLWriter) error {
	href := l.GetAttributeOr("href", "")
	navbarBaseUrl := l.GetAttribute("navbarBaseUrl")
	link := href
	if navbarBaseUrl != nil && *navbarBaseUrl != "" {
		link = *navbarBaseUrl + href
	}

	className := "mj-link"
	cssClass := l.GetAttribute("css-class")
	if cssClass != nil && *cssClass != "" {
		className += " " + *cssClass
	}

	if _, err := w.WriteString("<a "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(l).
		Add("class", className).
		Add("href", link).
		AddNullable("rel", l.GetAttribute("rel")).
		AddNullable("target", l.GetAttribute("target")).
		AddNullable("name", l.GetAttribute("name")).
		Add("style", "a").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if _, err := w.WriteString(l.GetContent()); err != nil {
		return err
	}

	if _, err := w.WriteString("\n</a>"); err != nil {
		return err
	}

	return nil
}

func (l *MjNavbarLink) Render(w core.MJMLWriter) error {
	if _, err := w.WriteString("<!--[if mso | IE]>"); err != nil {
		return err
	}

	if _, err := w.WriteString("<td "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(l).
		Add("style", "td").
		AddIf("class", helpers.SuffixCssClasses(l.GetAttributeOr("css-class", ""), "outlook"), l.GetAttribute("css-class") != nil).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if _, err := w.WriteString("<![endif]-->"); err != nil {
		return err
	}

	if err := l.renderContent(w); err != nil {
		return err
	}

	if _, err := w.WriteString("<!--[if mso | IE]></td><![endif]-->\n"); err != nil {
		return err
	}

	return nil
}
