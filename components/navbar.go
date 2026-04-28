package components

import (
	"fmt"

	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/core/helpers"
	"github.com/krozhkov/mgml/parser"
)

var MjNavbarSpec = &core.ComponentSpec{
	TagName: "mj-navbar",
	AllowedAttributes: map[string]string{
		"align":               "enum(left,center,right)",
		"base-url":            "string",
		"hamburger":           "string",
		"ico-align":           "enum(left,center,right)",
		"ico-open":            "string",
		"ico-close":           "string",
		"ico-color":           "color",
		"ico-font-size":       "unit(px,%)",
		"ico-font-family":     "string",
		"ico-text-transform":  "string",
		"ico-padding":         "unit(px,%){1,4}",
		"ico-padding-left":    "unit(px,%)",
		"ico-padding-top":     "unit(px,%)",
		"ico-padding-right":   "unit(px,%)",
		"ico-padding-bottom":  "unit(px,%)",
		"padding":             "unit(px,%){1,4}",
		"padding-left":        "unit(px,%)",
		"padding-top":         "unit(px,%)",
		"padding-right":       "unit(px,%)",
		"padding-bottom":      "unit(px,%)",
		"ico-text-decoration": "string",
		"ico-line-height":     "unit(px,%,)",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "align", Value: "center"},
		{Key: "ico-align", Value: "center"},
		{Key: "ico-open", Value: "&#9776;"},
		{Key: "ico-close", Value: "&#8855;"},
		{Key: "ico-color", Value: "#000000"},
		{Key: "ico-font-size", Value: "30px"},
		{Key: "ico-font-family", Value: "Ubuntu, Helvetica, Arial, sans-serif"},
		{Key: "ico-text-transform", Value: "uppercase"},
		{Key: "ico-padding", Value: "10px"},
		{Key: "ico-text-decoration", Value: "none"},
		{Key: "ico-line-height", Value: "30px"},
	},
	Create: NewMjNavbar,
}

type MjNavbar struct {
	*core.BodyComponent
}

func NewMjNavbar(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	navbar := &MjNavbar{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	navbar.BodyComponent.Component = navbar

	children, err := navbar.CreateChildren(nil, &core.CreateChildrenOptions{
		Attributes: []*core.Attribute{
			{Key: "navbarBaseUrl", Value: navbar.GetAttributeOr("base-url", "")},
		},
	})

	if err != nil {
		return navbar, err
	}

	navbar.Children = children

	return navbar, nil
}

func (n *MjNavbar) GetTagName() string {
	return n.Spec.TagName
}

func (n *MjNavbar) HeadStyle(breakpoint string) string {
	return fmt.Sprintf(`
      noinput.mj-menu-checkbox { display:block!important; max-height:none!important; visibility:visible!important; }

      @media only screen and (max-width:%s) {
        .mj-menu-checkbox[type="checkbox"] ~ .mj-inline-links { display:none!important; }
        .mj-menu-checkbox[type="checkbox"]:checked ~ .mj-inline-links,
        .mj-menu-checkbox[type="checkbox"] ~ .mj-menu-trigger { display:block!important; max-width:none!important; max-height:none!important; font-size:inherit!important; }
        .mj-menu-checkbox[type="checkbox"] ~ .mj-inline-links > a { display:block!important; }
        .mj-menu-checkbox[type="checkbox"]:checked ~ .mj-menu-trigger .mj-menu-icon-close { display:block!important; }
        .mj-menu-checkbox[type="checkbox"]:checked ~ .mj-menu-trigger .mj-menu-icon-open { display:none!important; }
      }
	`,
		helpers.MakeLowerBreakpoint(breakpoint),
	)
}

func (n *MjNavbar) GetStyles(element string) []*core.Style {
	switch element {
	case "div":
		return []*core.Style{
			{Name: "align", Value: n.GetAttributeOr("align", "")},
			{Name: "width", Value: "100%"},
		}
	case "label":
		return []*core.Style{
			{Name: "display", Value: "block"},
			{Name: "cursor", Value: "pointer"},
			{Name: "mso-hide", Value: "all"},
			{Name: "-moz-user-select", Value: "none"},
			{Name: "user-select", Value: "none"},
			{Name: "color", Value: n.GetAttributeOr("ico-color", "")},
			{Name: "font-size", Value: n.GetAttributeOr("ico-font-size", "")},
			{Name: "font-family", Value: n.GetAttributeOr("ico-font-family", "")},
			{Name: "text-transform", Value: n.GetAttributeOr("ico-text-transform", "")},
			{Name: "text-decoration", Value: n.GetAttributeOr("ico-text-decoration", "")},
			{Name: "line-height", Value: n.GetAttributeOr("ico-line-height", "")},
			{Name: "padding", Value: n.GetAttributeOr("ico-padding", "")},
			{Name: "padding-top", Value: n.GetAttributeOr("ico-padding-top", "")},
			{Name: "padding-right", Value: n.GetAttributeOr("ico-padding-right", "")},
			{Name: "padding-bottom", Value: n.GetAttributeOr("ico-padding-bottom", "")},
			{Name: "padding-left", Value: n.GetAttributeOr("ico-padding-left", "")},
		}
	case "trigger":
		return []*core.Style{
			{Name: "display", Value: "none"},
			{Name: "max-height", Value: "0px"},
			{Name: "max-width", Value: "0px"},
			{Name: "font-size", Value: "0px"},
			{Name: "overflow", Value: "hidden"},
		}
	case "icoOpen":
		return []*core.Style{
			{Name: "mso-hide", Value: "all"},
		}
	case "icoClose":
		return []*core.Style{
			{Name: "display", Value: "none"},
			{Name: "mso-hide", Value: "all"},
		}
	default:
		return nil
	}
}

func (n *MjNavbar) GetChildContext() *core.MJMLContext {
	return n.Context
}

func (n *MjNavbar) renderHamburger(w core.MJMLWriter) error {
	labelKey := helpers.GenRandomHexString(16)

	if _, err := w.WriteMsoConditionalTag(fmt.Sprintf(`<input type="checkbox" id="%s" class="mj-menu-checkbox" style="display:none !important; max-height:0; visibility:hidden;" />`, labelKey), true); err != nil {
		return err
	}

	if _, err := w.WriteString("<div "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(n).
		Add("class", "mj-menu-trigger").
		Add("style", "trigger").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><label "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(n).
		Add("for", labelKey).
		Add("class", "mj-menu-label").
		Add("style", "label").
		AddNullable("align", n.GetAttribute("ico-align")).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><span "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(n).
		Add("class", "mj-menu-icon-open").
		Add("style", "icoOpen").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if _, err := w.WriteString(n.GetAttributeOr("ico-open", "")); err != nil {
		return err
	}

	if _, err := w.WriteString("\n</span><span "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(n).
		Add("class", "mj-menu-icon-close").
		Add("style", "icoClose").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if _, err := w.WriteString(n.GetAttributeOr("ico-close", "")); err != nil {
		return err
	}

	if _, err := w.WriteString("</span></label></div>"); err != nil {
		return err
	}

	return nil
}

func (n *MjNavbar) Render(w core.MJMLWriter) error {
	if n.GetAttributeOr("hamburger", "") == "hamburger" {
		if err := n.renderHamburger(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("<div "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(n).
		Add("class", "mj-inline-links").
		// Add("style", "div"). // an error - this.htmlAttributes('div')
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if _, err := w.WriteConditionalTag(fmt.Sprintf(`<table role="presentation" border="0" cellpadding="0" cellspacing="0" align="%s"><tr>`, n.GetAttributeOr("align", "")), false); err != nil {
		return err
	}

	for _, child := range n.Children {
		if err := child.Render(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteConditionalTag("</tr></table>", false); err != nil {
		return err
	}

	if _, err := w.WriteString("</div>"); err != nil {
		return err
	}

	return nil
}
