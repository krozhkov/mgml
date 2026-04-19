package components

import (
	"strconv"

	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/core/helpers"
	"github.com/krozhkov/mgml/parser"
)

var MjGroupSpec = &core.ComponentSpec{
	TagName: "mj-group",
	AllowedAttributes: map[string]string{
		"background-color": "color",
		"direction":        "enum(ltr,rtl)",
		"vertical-align":   "enum(top,bottom,middle)",
		"width":            "unit(px,%)",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "direction", Value: "ltr"},
	},
	Create: NewMjGroup,
}

type MjGroup struct {
	*core.BodyComponent
}

func NewMjGroup(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	group := &MjGroup{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	group.BodyComponent.Component = group

	children, err := group.CreateChildren(nil, &core.CreateChildrenOptions{
		Attributes: []*core.Attribute{
			{Key: "mobileWidth", Value: "mobileWidth"},
		},
	})
	if err != nil {
		return group, err
	}

	group.Children = children

	return group, nil
}

func (g *MjGroup) GetTagName() string {
	return g.Spec.TagName
}

func (g *MjGroup) GetStyles(element string) []*core.Style {
	switch element {
	case "div":
		return []*core.Style{
			{Name: "font-size", Value: "0px"},
			{Name: "line-height", Value: "0"},
			{Name: "text-align", Value: "left"},
			{Name: "display", Value: "inline-block"},
			{Name: "width", Value: "100%"},
			{Name: "direction", Value: g.GetAttributeOr("direction", "")},
			{Name: "vertical-align", Value: g.GetAttributeOr("vertical-align", "")},
			{Name: "background-color", Value: g.GetAttributeOr("background-color", "")},
		}
	case "tdOutlook":
		return []*core.Style{
			{Name: "vertical-align", Value: g.GetAttributeOr("vertical-align", "")},
			{Name: "width", Value: g.GetWidthAsPixel()},
		}
	default:
		return nil
	}
}

func (g *MjGroup) GetChildContext() *core.MJMLContext {
	parentWidth := g.Context.ContainerWidth
	nonRawSiblings := g.Props.NonRawSiblings
	paddingSize := g.GetShorthandAttrValue("padding", "left") +
		g.GetShorthandAttrValue("padding", "right")

	containerWidth := g.GetAttributeOr("width", "")
	if containerWidth == "" {
		value := float64(parentWidth) / float64(nonRawSiblings)
		containerWidth = strconv.FormatFloat(value, 'f', -1, 64) + "px"
	}

	parsedWidth, unit := helpers.WidthParser(containerWidth, false)

	var calculatedWidth int
	if unit == "%" {
		calculatedWidth = int((float64(parentWidth)*parsedWidth)/100.0 - float64(paddingSize))
	} else {
		calculatedWidth = int(parsedWidth) - paddingSize
	}

	copy := *g.Context
	copy.ContainerWidth = calculatedWidth

	return &copy
}

func (g *MjGroup) Render(w core.MJMLWriter) error {
	groupWidth := g.GetChildContext().ContainerWidth

	classesName := g.getColumnClass() + " mj-outlook-group-fix"

	cssClass := g.GetAttributeOr("css-class", "")
	if cssClass != "" {
		classesName += " " + cssClass
	}

	if _, err := w.WriteString("<div "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(g).
		Add("class", classesName).
		Add("style", "div").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if _, err := w.WriteString("<!--[if mso | IE]>"); err != nil {
		return err
	}

	if _, err := w.WriteString("<table "); err != nil {
		return err
	}

	var bgColor *string
	if attr := g.GetAttribute("background-color"); attr != nil && *attr != "none" {
		bgColor = attr
	}

	if err := NewAttributesBuilder(g).
		AddNullable("bgcolor", bgColor).
		Add("border", "0").
		Add("cellpadding", "0").
		Add("cellspacing", "0").
		Add("role", "presentation").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><tr>"); err != nil {
		return err
	}

	if _, err := w.WriteString("<![endif]-->"); err != nil {
		return err
	}

	for _, child := range g.Children {
		if child.IsRawElement() {
			if err := child.Render(w); err != nil {
				return err
			}
		} else {
			if _, err := w.WriteString("<!--[if mso | IE]>"); err != nil {
				return err
			}

			if _, err := w.WriteString("<td style=\""); err != nil {
				return err
			}

			if err := FormatCssStyles(w, g.getCellStyles(child, groupWidth)); err != nil {
				return err
			}

			if _, err := w.WriteString("\">"); err != nil {
				return err
			}

			if _, err := w.WriteString("<![endif]-->"); err != nil {
				return err
			}

			if err := child.Render(w); err != nil {
				return err
			}

			if _, err := w.WriteString("<!--[if mso | IE]>"); err != nil {
				return err
			}

			if _, err := w.WriteString("</td>"); err != nil {
				return err
			}

			if _, err := w.WriteString("<![endif]-->"); err != nil {
				return err
			}
		}
	}

	if _, err := w.WriteString("<!--[if mso | IE]>"); err != nil {
		return err
	}

	if _, err := w.WriteString("</tr></table>"); err != nil {
		return err
	}

	if _, err := w.WriteString("<![endif]-->"); err != nil {
		return err
	}

	if _, err := w.WriteString("</div>\n"); err != nil {
		return err
	}

	return nil
}

func (g *MjGroup) GetWidthAsPixel() string {
	containerWidth := g.Context.ContainerWidth
	width := g.getWidth()

	parsedWidth, unit := helpers.WidthParser(width, false)

	if unit == "%" {
		pxWidth := (float64(containerWidth) * parsedWidth) / 100.0
		return strconv.FormatFloat(pxWidth, 'f', -1, 64) + "px"
	}

	return width
}

func (g *MjGroup) getWidth() string {
	nonRawSiblings := g.Props.NonRawSiblings
	width := g.GetAttributeOr("width", "")

	if width == "" {
		value := 100.0 / float64(nonRawSiblings)
		return strconv.FormatFloat(value, 'f', -1, 64) + "%"
	}

	return width
}

func (g *MjGroup) getColumnClass() string {
	var className string

	width := g.getWidth()
	parsedWidth, unit := helpers.WidthParser(width, false)

	switch unit {
	case "%":
		className = "mj-column-per-" + strconv.Itoa(int(parsedWidth))
	case "px":
		className = "mj-column-px-" + strconv.Itoa(int(parsedWidth))
	default:
		className = "mj-column-px-" + strconv.Itoa(int(parsedWidth))
	}

	// Add className to media queries
	g.Context.AddMediaQuery(className, parsedWidth, unit)

	return className
}

func (g *MjGroup) getElementWidth(width string, groupWidth int) string {
	nonRawSiblings := g.Props.NonRawSiblings
	containerWidth := g.Context.ContainerWidth

	if width == "" {
		pxWidth := float64(containerWidth) * float64(nonRawSiblings)
		return strconv.FormatFloat(pxWidth, 'f', -1, 64) + "px"
	}

	parsedWidth, unit := helpers.WidthParser(width, false)

	if unit == "%" {
		pxWidth := (100.0 * parsedWidth) / float64(groupWidth)
		return strconv.FormatFloat(pxWidth, 'f', -1, 64) + "px"
	}

	return strconv.FormatFloat(parsedWidth, 'f', -1, 64) + unit
}

func (g *MjGroup) getCellStyles(component core.Component, groupWidth int) []*core.Style {
	var width string
	if upgrade, ok := component.(interface {
		GetWidthAsPixel() string
	}); ok {
		width = upgrade.GetWidthAsPixel()
	} else {
		width = component.GetAttributeOr("width", "")
	}

	return []*core.Style{
		{Name: "align", Value: component.GetAttributeOr("align", "")},
		{Name: "vertical-align", Value: component.GetAttributeOr("vertical-align", "")},
		{Name: "width", Value: g.getElementWidth(width, groupWidth)},
	}
}
