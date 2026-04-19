package components

import (
	"slices"

	"github.com/elliotchance/orderedmap/v3"
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/parser"
)

var accordionElementChildrenAttributes = []string{
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

var MjAccordionElementSpec = &core.ComponentSpec{
	TagName: "mj-accordion-element",
	AllowedAttributes: map[string]string{
		"background-color":   "color",
		"border":             "string",
		"font-family":        "string",
		"icon-align":         "enum(top,middle,bottom)",
		"icon-width":         "unit(px,%)",
		"icon-height":        "unit(px,%)",
		"icon-wrapped-url":   "string",
		"icon-wrapped-alt":   "string",
		"icon-unwrapped-url": "string",
		"icon-unwrapped-alt": "string",
		"icon-position":      "enum(left,right)",
	},
	Create: NewMjAccordionElement,
}

type MjAccordionElement struct {
	*core.BodyComponent
}

func NewMjAccordionElement(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	element := &MjAccordionElement{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	element.BodyComponent.Component = element

	nodes := node.Children
	childrenAttr := element.getChildrenAttr()

	if slices.IndexFunc(nodes, func(n *parser.MJMLNode) bool { return n.TagName == "mj-accordion-title" }) == -1 {
		nodes = slices.Insert(nodes, 0, &parser.MJMLNode{
			TagName:    "mj-accordion-title",
			Attributes: orderedmap.NewOrderedMapWithElements(childrenAttr...),
		})
	}
	if slices.IndexFunc(nodes, func(n *parser.MJMLNode) bool { return n.TagName == "mj-accordion-text" }) == -1 {
		nodes = append(nodes, &parser.MJMLNode{
			TagName:    "mj-accordion-text",
			Attributes: orderedmap.NewOrderedMapWithElements(childrenAttr...),
		})
	}

	children, err := element.CreateChildren(nodes, &core.CreateChildrenOptions{
		Attributes: childrenAttr,
	})
	if err != nil {
		return element, err
	}

	element.Children = children

	return element, nil
}

func (e *MjAccordionElement) GetTagName() string {
	return e.Spec.TagName
}

func (e *MjAccordionElement) GetStyles(element string) []*core.Style {
	switch element {
	case "td":
		return []*core.Style{
			{Name: "padding", Value: "0px"},
			{Name: "background-color", Value: e.GetAttributeOr("background-color", "")},
		}
	case "label":
		return []*core.Style{
			{Name: "font-size", Value: "13px"},
			{Name: "font-family", Value: e.GetAttributeOr("font-family", "")},
		}
	case "input":
		return []*core.Style{
			{Name: "display", Value: "none"},
		}
	default:
		return nil
	}
}

func (e *MjAccordionElement) GetChildContext() *core.MJMLContext {
	return e.Context
}

func (e *MjAccordionElement) getChildrenAttr() []*core.Attribute {
	attr := make([]*core.Attribute, 0, len(accordionElementChildrenAttributes))

	for _, name := range accordionElementChildrenAttributes {
		value := e.GetAttribute(name)
		if value != nil {
			attr = append(attr, &core.Attribute{Key: name, Value: *value})
		}
	}

	return attr
}

func (e *MjAccordionElement) Render(w core.MJMLWriter) error {
	if _, err := w.WriteString("<tr "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(e).
		AddNullable("class", e.GetAttribute("css-class")).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><td "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(e).
		Add("style", "td").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><label "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(e).
		Add("class", "mj-accordion-element").
		Add("style", "label").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if _, err := w.WriteString("<!--[if !mso | IE]><!--><input "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(e).
		Add("class", "mj-accordion-checkbox").
		Add("type", "checkbox").
		Add("style", "input").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("/><!--<![endif]-->"); err != nil {
		return err
	}

	if _, err := w.WriteString("<div>"); err != nil {
		return err
	}

	for _, child := range e.Children {
		if err := child.Render(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("</div></label></td></tr>\n"); err != nil {
		return err
	}

	return nil
}
