package components

import (
	"strconv"

	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/core/helpers"
	"github.com/krozhkov/mgml/parser"
)

var MjTableSpec = &core.ComponentSpec{
	TagName:   "mj-table",
	EndingTag: true,
	AllowedAttributes: map[string]string{
		"align":                      "enum(left,right,center)",
		"border":                     "string",
		"cellpadding":                "integer",
		"cellspacing":                "integer",
		"container-background-color": "color",
		"color":                      "color",
		"font-family":                "string",
		"font-size":                  "unit(px)",
		"font-weight":                "string",
		"line-height":                "unit(px,%,)",
		"padding-bottom":             "unit(px,%)",
		"padding-left":               "unit(px,%)",
		"padding-right":              "unit(px,%)",
		"padding-top":                "unit(px,%)",
		"padding":                    "unit(px,%){1,4}",
		"role":                       "enum(none,presentation)",
		"table-layout":               "enum(auto,fixed,initial,inherit)",
		"vertical-align":             "enum(top,bottom,middle)",
		"width":                      "unit(px,%,auto)",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "align", Value: "left"},
		{Key: "border", Value: "none"},
		{Key: "cellpadding", Value: "0"},
		{Key: "cellspacing", Value: "0"},
		{Key: "color", Value: "#000000"},
		{Key: "font-family", Value: "Ubuntu, Helvetica, Arial, sans-serif"},
		{Key: "font-size", Value: "13px"},
		{Key: "line-height", Value: "22px"},
		{Key: "padding", Value: "10px 25px"},
		{Key: "table-layout", Value: "auto"},
		{Key: "width", Value: "100%"},
	},
	Create: NewMjTable,
}

type MjTable struct {
	*core.BodyComponent
}

func NewMjTable(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	table := &MjTable{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	table.BodyComponent.Component = table

	return table, nil
}

func (t *MjTable) GetTagName() string {
	return t.Spec.TagName
}

func (t *MjTable) GetStyles(element string) []*core.Style {
	switch element {
	case "table":
		var borderCollapse string
		if t.hasCellspacing() {
			borderCollapse = "separate"
		}

		return []*core.Style{
			{Name: "color", Value: t.GetAttributeOr("color", "")},
			{Name: "font-family", Value: t.GetAttributeOr("font-family", "")},
			{Name: "font-size", Value: t.GetAttributeOr("font-size", "")},
			{Name: "line-height", Value: t.GetAttributeOr("line-height", "")},
			{Name: "table-layout", Value: t.GetAttributeOr("table-layout", "")},
			{Name: "width", Value: t.GetAttributeOr("width", "")},
			{Name: "border", Value: t.GetAttributeOr("border", "")},
			{Name: "border-collapse", Value: borderCollapse},
		}
	default:
		return nil
	}
}

func (t *MjTable) GetChildContext() *core.MJMLContext {
	return t.Context
}

func (t *MjTable) getWidth() string {
	width := t.GetAttributeOr("width", "")

	if width == "auto" {
		return width
	}

	parsedWidth, unit := helpers.WidthParser(width, true)

	switch unit {
	case "%":
		return width
	default:
		return strconv.FormatFloat(parsedWidth, 'f', -1, 64)
	}
}

func (t *MjTable) hasCellspacing() bool {
	cellspacing := t.GetAttributeOr("cellspacing", "")
	numericValue, err := helpers.ParseFloatLoose(cellspacing)
	return err == nil && numericValue > 0
}

func (t *MjTable) Render(w core.MJMLWriter) error {
	if _, err := w.WriteString("<table "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(t).
		AddNullable("cellpadding", t.GetAttribute("cellpadding")).
		AddNullable("cellspacing", t.GetAttribute("cellspacing")).
		AddNullable("role", t.GetAttribute("role")).
		Add("width", t.getWidth()).
		Add("border", "0").
		Add("style", "table").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if _, err := w.WriteString(t.GetContent()); err != nil {
		return err
	}

	if _, err := w.WriteString("</table>\n"); err != nil {
		return err
	}

	return nil
}
