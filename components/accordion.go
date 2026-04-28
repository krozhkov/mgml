package components

import (
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/parser"
)

var accordionChildrenAttributes = []string{
	"border",
	"icon-align",
	"icon-width",
	"icon-height",
	"icon-position",
	"icon-wrapped-url",
	"icon-wrapped-alt",
	"icon-unwrapped-url",
	"icon-unwrapped-alt",
}

var MjAccordionSpec = &core.ComponentSpec{
	TagName: "mj-accordion",
	AllowedAttributes: map[string]string{
		"container-background-color": "color",
		"border":                     "string",
		"font-family":                "string",
		"icon-align":                 "enum(top,middle,bottom)",
		"icon-width":                 "unit(px,%)",
		"icon-height":                "unit(px,%)",
		"icon-wrapped-url":           "string",
		"icon-wrapped-alt":           "string",
		"icon-unwrapped-url":         "string",
		"icon-unwrapped-alt":         "string",
		"icon-position":              "enum(left,right)",
		"padding-bottom":             "unit(px,%)",
		"padding-left":               "unit(px,%)",
		"padding-right":              "unit(px,%)",
		"padding-top":                "unit(px,%)",
		"padding":                    "unit(px,%){1,4}",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "border", Value: "2px solid black"},
		{Key: "font-family", Value: "Ubuntu, Helvetica, Arial, sans-serif"},
		{Key: "icon-align", Value: "middle"},
		{Key: "icon-wrapped-url", Value: "https://i.imgur.com/bIXv1bk.png"},
		{Key: "icon-wrapped-alt", Value: "+"},
		{Key: "icon-unwrapped-url", Value: "https://i.imgur.com/w4uTygT.png"},
		{Key: "icon-unwrapped-alt", Value: "-"},
		{Key: "icon-position", Value: "right"},
		{Key: "icon-height", Value: "32px"},
		{Key: "icon-width", Value: "32px"},
		{Key: "padding", Value: "10px 25px"},
	},
	Create: NewMjAccordion,
}

type MjAccordion struct {
	*core.BodyComponent
}

func NewMjAccordion(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	accordion := &MjAccordion{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	accordion.BodyComponent.Component = accordion

	childrenAttr := accordion.getChildrenAttr()
	children, err := accordion.CreateChildren(nil, &core.CreateChildrenOptions{
		Attributes: childrenAttr,
	})
	if err != nil {
		return accordion, err
	}

	accordion.Children = children

	return accordion, nil
}

func (a *MjAccordion) GetTagName() string {
	return a.Spec.TagName
}

func (a *MjAccordion) HeadStyle(breakpoint string) string {
	return `
      noinput.mj-accordion-checkbox { display:block!important; }

      @media yahoo, only screen and (min-width:0) {
        .mj-accordion-element { display:block; }
        input.mj-accordion-checkbox, .mj-accordion-less { display:none!important; }
        input.mj-accordion-checkbox + * .mj-accordion-title { cursor:pointer; touch-action:manipulation; -webkit-user-select:none; -moz-user-select:none; user-select:none; }
        input.mj-accordion-checkbox + * .mj-accordion-content { overflow:hidden; display:none; }
        input.mj-accordion-checkbox + * .mj-accordion-more { display:block!important; }
        input.mj-accordion-checkbox:checked + * .mj-accordion-content { display:block; }
        input.mj-accordion-checkbox:checked + * .mj-accordion-more { display:none!important; }
        input.mj-accordion-checkbox:checked + * .mj-accordion-less { display:block!important; }
      }

      .moz-text-html input.mj-accordion-checkbox + * .mj-accordion-title { cursor: auto; touch-action: auto; -webkit-user-select: auto; -moz-user-select: auto; user-select: auto; }
      .moz-text-html input.mj-accordion-checkbox + * .mj-accordion-content { overflow: hidden; display: block; }
      .moz-text-html input.mj-accordion-checkbox + * .mj-accordion-ico { display: none; }

      @goodbye { @gmail }
    `
}

func (a *MjAccordion) GetStyles(element string) []*core.Style {
	switch element {
	case "table":
		return []*core.Style{
			{Name: "width", Value: "100%"},
			{Name: "border-collapse", Value: "collapse"},
			{Name: "border", Value: a.GetAttributeOr("border", "")},
			{Name: "border-bottom", Value: "none"},
			{Name: "font-family", Value: a.GetAttributeOr("font-family", "")},
		}
	default:
		return nil
	}
}

func (a *MjAccordion) GetChildContext() *core.MJMLContext {
	return a.Context
}

func (a *MjAccordion) getChildrenAttr() []*core.Attribute {
	attr := make([]*core.Attribute, 0, len(accordionChildrenAttributes))

	for _, name := range accordionChildrenAttributes {
		value := a.GetAttribute(name)
		if value != nil {
			attr = append(attr, &core.Attribute{Key: name, Value: *value})
		}
	}

	accordionFontFamily := a.GetAttribute("font-family")
	if accordionFontFamily != nil {
		attr = append(attr, &core.Attribute{Key: "accordionFontFamily", Value: *accordionFontFamily})
	}

	return attr
}

func (a *MjAccordion) Render(w core.MJMLWriter) error {
	if _, err := w.WriteString("<table "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(a).
		Add("cellspacing", "0").
		Add("cellpadding", "0").
		Add("class", "mj-accordion").
		Add("style", "table").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><tbody>"); err != nil {
		return err
	}

	for _, child := range a.Children {
		if err := child.Render(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("</tbody></table>\n"); err != nil {
		return err
	}

	return nil
}
