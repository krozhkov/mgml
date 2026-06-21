package components

import (
	"fmt"
	"strconv"

	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/core/helpers"
	"github.com/krozhkov/mgml/parser"
)

var MjButtonSpec = &core.ComponentSpec{
	TagName:   "mj-button",
	EndingTag: true,
	AllowedAttributes: map[string]string{
		"align":                      "enum(left,center,right)",
		"background-color":           "color",
		"border-bottom":              "string",
		"border-left":                "string",
		"border-radius":              "string",
		"border-right":               "string",
		"border-top":                 "string",
		"border":                     "string",
		"color":                      "color",
		"container-background-color": "color",
		"font-family":                "string",
		"font-size":                  "unit(px)",
		"font-style":                 "string",
		"font-weight":                "string",
		"height":                     "unit(px,%)",
		"href":                       "string",
		"name":                       "string",
		"title":                      "string",
		"inner-padding":              "unit(px,%){1,4}",
		"letter-spacing":             "unitWithNegative(px,em)",
		"line-height":                "unit(px,%,)",
		"padding-bottom":             "unit(px,%)",
		"padding-left":               "unit(px,%)",
		"padding-right":              "unit(px,%)",
		"padding-top":                "unit(px,%)",
		"padding":                    "unit(px,%){1,4}",
		"rel":                        "string",
		"target":                     "string",
		"text-decoration":            "string",
		"text-transform":             "string",
		"vertical-align":             "enum(top,bottom,middle)",
		"text-align":                 "enum(left,right,center)",
		"width":                      "unit(px,%)",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "align", Value: "center"},
		{Key: "background-color", Value: "#414141"},
		{Key: "border", Value: "none"},
		{Key: "border-radius", Value: "3px"},
		{Key: "color", Value: "#ffffff"},
		{Key: "font-family", Value: "Ubuntu, Helvetica, Arial, sans-serif"},
		{Key: "font-size", Value: "13px"},
		{Key: "font-weight", Value: "normal"},
		{Key: "inner-padding", Value: "10px 25px"},
		{Key: "line-height", Value: "120%"},
		{Key: "padding", Value: "10px 25px"},
		{Key: "target", Value: "_blank"},
		{Key: "text-decoration", Value: "none"},
		{Key: "text-transform", Value: "none"},
		{Key: "vertical-align", Value: "middle"},
	},
	Create: NewMjButton,
}

type MjButton struct {
	*core.BodyComponent
}

func NewMjButton(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	button := &MjButton{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	button.BodyComponent.Component = button

	return button, nil
}

func (b *MjButton) GetTagName() string {
	return b.Spec.TagName
}

func (b *MjButton) GetStyles(element string) []*core.Style {
	switch element {
	case "table":
		return []*core.Style{
			{Name: "border-collapse", Value: "separate"},
			{Name: "width", Value: b.GetAttributeOr("width", "")},
			{Name: "line-height", Value: "100%"},
		}
	case "td":
		return []*core.Style{
			{Name: "border", Value: b.GetAttributeOr("border", "")},
			{Name: "border-bottom", Value: b.GetAttributeOr("border-bottom", "")},
			{Name: "border-left", Value: b.GetAttributeOr("border-left", "")},
			{Name: "border-radius", Value: b.GetAttributeOr("border-radius", "")},
			{Name: "border-right", Value: b.GetAttributeOr("border-right", "")},
			{Name: "border-top", Value: b.GetAttributeOr("border-top", "")},
			{Name: "cursor", Value: "auto"},
			{Name: "font-style", Value: b.GetAttributeOr("font-style", "")},
			{Name: "height", Value: b.GetAttributeOr("height", "")},
			{Name: "mso-padding-alt", Value: b.GetAttributeOr("inner-padding", "")},
			{Name: "text-align", Value: b.GetAttributeOr("text-align", "")},
			{Name: "background", Value: b.GetAttributeOr("background-color", "")},
		}
	case "content":
		return []*core.Style{
			{Name: "display", Value: "inline-block"},
			{Name: "width", Value: b.calculateAWidth(b.GetAttributeOr("width", ""))},
			{Name: "background", Value: b.GetAttributeOr("background-color", "")},
			{Name: "color", Value: b.GetAttributeOr("color", "")},
			{Name: "font-family", Value: b.GetAttributeOr("font-family", "")},
			{Name: "font-size", Value: b.GetAttributeOr("font-size", "")},
			{Name: "font-style", Value: b.GetAttributeOr("font-style", "")},
			{Name: "font-weight", Value: b.GetAttributeOr("font-weight", "")},
			{Name: "line-height", Value: b.GetAttributeOr("line-height", "")},
			{Name: "letter-spacing", Value: b.GetAttributeOr("letter-spacing", "")},
			{Name: "margin", Value: "0"},
			{Name: "text-decoration", Value: b.GetAttributeOr("text-decoration", "")},
			{Name: "text-transform", Value: b.GetAttributeOr("text-transform", "")},
			{Name: "padding", Value: b.GetAttributeOr("inner-padding", "")},
			{Name: "mso-padding-alt", Value: "0px"},
			{Name: "border-radius", Value: b.GetAttributeOr("border-radius", "")},
		}
	default:
		return nil
	}
}

func (b *MjButton) GetChildContext() *core.MJMLContext {
	return b.Context
}

func (b *MjButton) calculateAWidth(width string) string {
	if width == "" {
		return ""
	}

	parsedWidth, unit := helpers.WidthParser(width, true)

	// impossible to handle percents because it depends on padding and text width
	if unit != "px" {
		return ""
	}

	boxWidths := b.GetBoxWidths()

	innerPaddings := b.GetShorthandAttrValue("inner-padding", "left") +
		b.GetShorthandAttrValue("inner-padding", "right")

	aWidth := parsedWidth - float64(innerPaddings) - float64(boxWidths.Borders)

	return strconv.FormatFloat(aWidth, 'f', -1, 64) + "px"
}

func (b *MjButton) Render(w core.MJMLWriter) error {
	href := b.GetAttribute("href")

	var tag string
	if href != nil && *href != "" {
		tag = "a"
	} else {
		tag = "p"
		href = nil
	}

	if _, err := w.WriteString("<table "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(b).
		Add("border", "0").
		Add("cellpadding", "0").
		Add("cellspacing", "0").
		Add("role", "presentation").
		Add("style", "table").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if _, err := w.WriteString("<tbody><tr><td "); err != nil {
		return err
	}

	bgColor := b.GetAttribute("background-color")
	if bgColor != nil && *bgColor == "none" {
		bgColor = nil
	}

	if err := NewAttributesBuilder(b).
		Add("align", "center").
		AddNullable("bgcolor", bgColor).
		Add("role", "presentation").
		Add("style", "td").
		AddNullable("valign", b.GetAttribute("vertical-align")).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(fmt.Sprintf("><%s ", tag)); err != nil {
		return err
	}

	var target *string
	if tag == "a" {
		target = b.GetAttribute("target")
	}

	if err := NewAttributesBuilder(b).
		AddNullable("href", href).
		AddNullable("name", b.GetAttribute("name")).
		AddNullable("rel", b.GetAttribute("rel")).
		AddNullable("title", b.GetAttribute("title")).
		Add("style", "content").
		AddNullable("target", target).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if _, err := w.WriteString(b.GetContent()); err != nil {
		return err
	}

	if _, err := w.WriteString(fmt.Sprintf("</%s>", tag)); err != nil {
		return err
	}

	if _, err := w.WriteString("</td></tr></tbody></table>"); err != nil {
		return err
	}

	return nil
}
