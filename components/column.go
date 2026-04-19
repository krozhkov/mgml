package components

import (
	"strconv"
	"strings"

	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/core/helpers"
	"github.com/krozhkov/mgml/parser"
)

var MjColumnSpec = &core.ComponentSpec{
	TagName: "mj-column",
	AllowedAttributes: map[string]string{
		"background-color":       "color",
		"border":                 "string",
		"border-bottom":          "string",
		"border-left":            "string",
		"border-radius":          "unit(px,%){1,4}",
		"border-right":           "string",
		"border-top":             "string",
		"direction":              "enum(ltr,rtl)",
		"inner-background-color": "color",
		"padding-bottom":         "unit(px,%)",
		"padding-left":           "unit(px,%)",
		"padding-right":          "unit(px,%)",
		"padding-top":            "unit(px,%)",
		"inner-border":           "string",
		"inner-border-bottom":    "string",
		"inner-border-left":      "string",
		"inner-border-radius":    "unit(px,%){1,4}",
		"inner-border-right":     "string",
		"inner-border-top":       "string",
		"padding":                "unit(px,%){1,4}",
		"vertical-align":         "enum(top,bottom,middle)",
		"width":                  "unit(px,%)",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "direction", Value: "ltr"},
		{Key: "vertical-align", Value: "top"},
	},
	Create: NewMjColumn,
}

type MjColumn struct {
	*core.BodyComponent
}

func NewMjColumn(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	column := &MjColumn{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	column.BodyComponent.Component = column

	children, err := column.CreateChildren(nil, &core.CreateChildrenOptions{})
	if err != nil {
		return column, err
	}

	column.Children = children

	return column, nil
}

func (c *MjColumn) GetTagName() string {
	return c.Spec.TagName
}

func (c *MjColumn) GetStyles(element string) []*core.Style {
	switch element {
	case "div":
		return []*core.Style{
			{Name: "font-size", Value: "0px"},
			{Name: "text-align", Value: "left"},
			{Name: "direction", Value: c.GetAttributeOr("direction", "")},
			{Name: "display", Value: "inline-block"},
			{Name: "vertical-align", Value: c.GetAttributeOr("vertical-align", "")},
			{Name: "width", Value: c.getMobileWidth()},
		}
	case "table":
		if c.hasGutter() {
			return []*core.Style{
				{Name: "background-color", Value: c.GetAttributeOr("inner-background-color", "")},
				{Name: "border", Value: c.GetAttributeOr("inner-border", "")},
				{Name: "border-bottom", Value: c.GetAttributeOr("inner-border-bottom", "")},
				{Name: "border-left", Value: c.GetAttributeOr("inner-border-left", "")},
				{Name: "border-radius", Value: c.GetAttributeOr("inner-border-radius", "")},
				{Name: "border-right", Value: c.GetAttributeOr("inner-border-right", "")},
				{Name: "border-top", Value: c.GetAttributeOr("inner-border-top", "")},
			}
		}

		return []*core.Style{
			{Name: "background-color", Value: c.GetAttributeOr("background-color", "")},
			{Name: "border", Value: c.GetAttributeOr("border", "")},
			{Name: "border-bottom", Value: c.GetAttributeOr("border-bottom", "")},
			{Name: "border-left", Value: c.GetAttributeOr("border-left", "")},
			{Name: "border-radius", Value: c.GetAttributeOr("border-radius", "")},
			{Name: "border-right", Value: c.GetAttributeOr("border-right", "")},
			{Name: "border-top", Value: c.GetAttributeOr("border-top", "")},
			{Name: "vertical-align", Value: c.GetAttributeOr("vertical-align", "")},
		}
	case "tdOutlook":
		return []*core.Style{
			{Name: "vertical-align", Value: c.GetAttributeOr("vertical-align", "")},
			{Name: "width", Value: c.GetWidthAsPixel()},
		}
	case "gutter":
		return []*core.Style{
			{Name: "background-color", Value: c.GetAttributeOr("background-color", "")},
			{Name: "border", Value: c.GetAttributeOr("border", "")},
			{Name: "border-bottom", Value: c.GetAttributeOr("border-bottom", "")},
			{Name: "border-left", Value: c.GetAttributeOr("border-left", "")},
			{Name: "border-radius", Value: c.GetAttributeOr("border-radius", "")},
			{Name: "border-right", Value: c.GetAttributeOr("border-right", "")},
			{Name: "border-top", Value: c.GetAttributeOr("border-top", "")},
			{Name: "vertical-align", Value: c.GetAttributeOr("vertical-align", "")},
			{Name: "padding", Value: c.GetAttributeOr("padding", "")},
			{Name: "padding-top", Value: c.GetAttributeOr("padding-top", "")},
			{Name: "padding-right", Value: c.GetAttributeOr("padding-right", "")},
			{Name: "padding-bottom", Value: c.GetAttributeOr("padding-bottom", "")},
			{Name: "padding-left", Value: c.GetAttributeOr("padding-left", "")},
		}
	default:
		return nil
	}
}

func (c *MjColumn) GetChildContext() *core.MJMLContext {
	parentWidth := c.Context.ContainerWidth
	nonRawSiblings := c.Props.NonRawSiblings
	boxWidths := c.GetBoxWidths()
	borders := boxWidths.Borders
	paddings := boxWidths.Paddings

	innerBorders :=
		c.GetShorthandBorderValue("left", "inner-border") +
			c.GetShorthandBorderValue("right", "inner-border")

	allPaddings := paddings + borders + innerBorders

	containerWidth := c.GetAttributeOr("width", "")
	if containerWidth == "" {
		containerWidth = strconv.Itoa(parentWidth/nonRawSiblings) + "px"
	}

	parsedWidth, unit := helpers.WidthParser(containerWidth, false)

	var calculatedWidth int
	if unit == "%" {
		calculatedWidth = int((float64(parentWidth)*parsedWidth)/100.0 - float64(allPaddings))
	} else {
		calculatedWidth = int(parsedWidth) - allPaddings
	}

	copy := *c.Context
	copy.ContainerWidth = calculatedWidth

	return &copy
}

func (c *MjColumn) Render(w core.MJMLWriter) error {
	classesName := c.getColumnClass() + " mj-outlook-group-fix"

	cssClass := c.GetAttributeOr("css-class", "")
	if cssClass != "" {
		classesName += " " + cssClass
	}

	if _, err := w.WriteString("<div "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(c).
		Add("class", classesName).
		Add("style", "div").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if c.hasGutter() {
		if err := c.renderGutter(w); err != nil {
			return err
		}
	} else {
		if err := c.renderColumn(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("</div>"); err != nil {
		return err
	}

	return nil
}

func (c *MjColumn) hasGutter() bool {
	return c.GetAttribute("padding") != nil ||
		c.GetAttribute("padding-bottom") != nil ||
		c.GetAttribute("padding-left") != nil ||
		c.GetAttribute("padding-right") != nil ||
		c.GetAttribute("padding-top") != nil
}

func (c *MjColumn) getMobileWidth() string {
	context := c.Context
	containerWidth := context.ContainerWidth
	nonRawSiblings := c.Props.NonRawSiblings
	width := c.GetAttributeOr("width", "")
	mobileWidth := c.GetAttributeOr("mobileWidth", "")

	if mobileWidth != "mobileWidth" {
		return "100%"
	}
	if width == "" {
		return strconv.Itoa(100/nonRawSiblings) + "%"
	}

	parsedWidth, unit := helpers.WidthParser(width, false)

	switch unit {
	case "%":
		return width
	case "px":
		return strconv.Itoa((int(parsedWidth)/containerWidth)*100) + "%"
	default:
		return strconv.Itoa((int(parsedWidth)/containerWidth)*100) + "%"
	}
}

func (c *MjColumn) GetWidthAsPixel() string {
	context := c.Context
	containerWidth := context.ContainerWidth
	width := c.getWidth()

	parsedWidth, unit := helpers.WidthParser(width, false)

	if unit == "%" {
		pxWidth := (float64(containerWidth) * parsedWidth) / 100.0
		return strconv.FormatFloat(pxWidth, 'f', -1, 64) + "px"
	}

	return width
}

func (c *MjColumn) getWidth() string {
	nonRawSiblings := c.Props.NonRawSiblings
	width := c.GetAttribute("width")

	if width == nil || *width == "" {
		value := 100.0 / float64(nonRawSiblings)
		return strconv.FormatFloat(value, 'f', -1, 64) + "%"
	}

	return *width
}

func (c *MjColumn) getColumnClass() string {
	var className string

	width := c.getWidth()
	parsedWidth, unit := helpers.WidthParser(width, false)

	formattedClassNb := strings.Replace(strconv.FormatFloat(parsedWidth, 'f', -1, 64), ".", "-", 1)

	switch unit {
	case "%":
		className = "mj-column-per-" + formattedClassNb
	case "px":
		className = "mj-column-px-" + formattedClassNb
	default:
		className = "mj-column-px-" + formattedClassNb
	}

	// Add className to media queries
	c.Context.AddMediaQuery(className, parsedWidth, unit)

	return className
}

func (c *MjColumn) getCellStyles(component core.Component) []*core.Style {
	return []*core.Style{
		{Name: "background", Value: component.GetAttributeOr("container-background-color", "")},
		{Name: "font-size", Value: "0px"},
		{Name: "padding", Value: component.GetAttributeOr("padding", "")},
		{Name: "padding-top", Value: component.GetAttributeOr("padding-top", "")},
		{Name: "padding-right", Value: component.GetAttributeOr("padding-right", "")},
		{Name: "padding-bottom", Value: component.GetAttributeOr("padding-bottom", "")},
		{Name: "padding-left", Value: component.GetAttributeOr("padding-left", "")},
		{Name: "word-break", Value: "break-word"},
	}
}

func (c *MjColumn) renderColumn(w core.MJMLWriter) error {
	if _, err := w.WriteString("<table "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(c).
		Add("border", "0").
		Add("cellpadding", "0").
		Add("cellspacing", "0").
		Add("role", "presentation").
		Add("style", "table").
		Add("width", "100%").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><tbody>"); err != nil {
		return err
	}

	for _, child := range c.Children {
		if child.IsRawElement() {
			if err := child.Render(w); err != nil {
				return err
			}
		} else {
			if _, err := w.WriteString("<tr><td "); err != nil {
				return err
			}

			if err := NewAttributesBuilder(child).
				AddNullable("align", child.GetAttribute("align")).
				AddNullable("class", child.GetAttribute("css-class")).
				Write(w); err != nil {
				return err
			}

			if _, err := w.WriteString(" style=\""); err != nil {
				return err
			}

			if err := FormatCssStyles(w, c.getCellStyles(child)); err != nil {
				return err
			}

			if _, err := w.WriteString("\">"); err != nil {
				return err
			}

			if err := child.Render(w); err != nil {
				return err
			}

			if _, err := w.WriteString("</td></tr>"); err != nil {
				return err
			}
		}
	}

	if _, err := w.WriteString("</tbody></table>"); err != nil {
		return err
	}

	return nil
}

func (c *MjColumn) renderGutter(w core.MJMLWriter) error {
	if _, err := w.WriteString("<table "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(c).
		Add("border", "0").
		Add("cellpadding", "0").
		Add("cellspacing", "0").
		Add("role", "presentation").
		Add("width", "100%").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><tbody><tr><td "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(c).
		Add("style", "gutter").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if err := c.renderColumn(w); err != nil {
		return err
	}

	if _, err := w.WriteString("</td></tr></tbody></table>"); err != nil {
		return err
	}

	return nil
}
