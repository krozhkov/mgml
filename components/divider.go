package components

import (
	"io"
	"strconv"

	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/core/helpers"
	"github.com/krozhkov/mgml/parser"
)

var MjDividerSpec = &core.ComponentSpec{
	TagName: "mj-divider",
	AllowedAttributes: map[string]string{
		"border-color":               "color",
		"border-style":               "string",
		"border-width":               "unit(px)",
		"container-background-color": "color",
		"padding":                    "unit(px,%){1,4}",
		"padding-bottom":             "unit(px,%)",
		"padding-left":               "unit(px,%)",
		"padding-right":              "unit(px,%)",
		"padding-top":                "unit(px,%)",
		"width":                      "unit(px,%)",
		"align":                      "enum(left,center,right)",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "border-color", Value: "#000000"},
		{Key: "border-style", Value: "solid"},
		{Key: "border-width", Value: "4px"},
		{Key: "padding", Value: "10px 25px"},
		{Key: "width", Value: "100%"},
		{Key: "align", Value: "center"},
	},
	Create: NewMjDivider,
}

type MjDivider struct {
	*core.BodyComponent
}

func NewMjDivider(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	divider := &MjDivider{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	divider.BodyComponent.Component = divider

	return divider, nil
}

func (d *MjDivider) GetTagName() string {
	return d.Spec.TagName
}

func (d *MjDivider) GetStyles(element string) []*core.Style {
	computeAlign := "0px auto"
	align := d.GetAttributeOr("align", "")
	if align == "left" {
		computeAlign = "0px"
	} else if align == "right" {
		computeAlign = "0px 0px 0px auto"
	}

	borderTopStyle := d.GetAttributeOr("border-style", "")
	borderTopWidth := d.GetAttributeOr("border-width", "")
	borderTopColor := d.GetAttributeOr("border-color", "")
	borderTop := borderTopStyle + " " + borderTopWidth + " " + borderTopColor

	attrs := []*core.Style{
		{Name: "border-top", Value: borderTop},
		{Name: "font-size", Value: "1px"},
		{Name: "margin", Value: computeAlign},
	}

	switch element {
	case "p":
		attrs = append(attrs, &core.Style{Name: "width", Value: d.GetAttributeOr("width", "")})
		return attrs
	case "outlook":
		attrs = append(attrs, &core.Style{Name: "width", Value: d.getOutlookWidth()})
		return attrs
	default:
		return nil
	}
}

func (d *MjDivider) GetChildContext() *core.MJMLContext {
	return d.Context
}

func (d *MjDivider) getOutlookWidth() string {
	containerWidth := d.Context.ContainerWidth
	paddingSize := d.GetShorthandAttrValue("padding", "left") +
		d.GetShorthandAttrValue("padding", "right")

	width := d.GetAttributeOr("width", "")

	parsedWidth, unit := helpers.WidthParser(width, true)

	switch unit {
	case "%":
		effectiveWidth := containerWidth - paddingSize
		percentMultiplier := parsedWidth / 100.0
		value := float64(effectiveWidth) * percentMultiplier
		return strconv.FormatFloat(value, 'f', -1, 64) + "px"
	case "px":
		return width
	default:
		return strconv.Itoa(containerWidth-paddingSize) + "px"
	}
}

func (d *MjDivider) renderAfter(w io.StringWriter) error {
	if _, err := w.WriteString("<!--[if mso | IE]><table "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(d).
		AddNullable("align", d.GetAttribute("align")).
		Add("border", "0").
		Add("cellpadding", "0").
		Add("cellspacing", "0").
		Add("style", "outlook").
		Add("role", "presentation").
		Add("width", d.getOutlookWidth()).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><tr><td style=\"height:0;line-height:0;\">&nbsp;</td></tr></table><![endif]-->"); err != nil {
		return err
	}

	return nil
}

func (d *MjDivider) Render(w core.MJMLWriter) error {
	if _, err := w.WriteString("<p "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(d).
		Add("style", "p").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("></p>"); err != nil {
		return err
	}

	if err := d.renderAfter(w); err != nil {
		return err
	}

	return nil
}
