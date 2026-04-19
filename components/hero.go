package components

import (
	"strconv"
	"strings"

	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/core/helpers"
	"github.com/krozhkov/mgml/parser"
)

var MjHeroSpec = &core.ComponentSpec{
	TagName: "mj-hero",
	AllowedAttributes: map[string]string{
		"mode":                       "string",
		"height":                     "unit(px,%)",
		"background-url":             "string",
		"background-width":           "unit(px,%)",
		"background-height":          "unit(px,%)",
		"background-position":        "string",
		"border-radius":              "string",
		"container-background-color": "color",
		"inner-background-color":     "color",
		"inner-padding":              "unit(px,%){1,4}",
		"inner-padding-top":          "unit(px,%)",
		"inner-padding-left":         "unit(px,%)",
		"inner-padding-right":        "unit(px,%)",
		"inner-padding-bottom":       "unit(px,%)",
		"padding":                    "unit(px,%){1,4}",
		"padding-bottom":             "unit(px,%)",
		"padding-left":               "unit(px,%)",
		"padding-right":              "unit(px,%)",
		"padding-top":                "unit(px,%)",
		"background-color":           "color",
		"vertical-align":             "enum(top,bottom,middle)",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "mode", Value: "fixed-height"},
		{Key: "height", Value: "0px"},
		{Key: "background-position", Value: "center center"},
		{Key: "padding", Value: "0px"},
		{Key: "background-color", Value: "#ffffff"},
		{Key: "vertical-align", Value: "top"},
	},
	Create: NewMjHero,
}

type MjHero struct {
	*core.BodyComponent
}

func NewMjHero(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	hero := &MjHero{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	hero.BodyComponent.Component = hero

	children, err := hero.CreateChildren(nil, &core.CreateChildrenOptions{})
	if err != nil {
		return hero, err
	}

	hero.Children = children

	return hero, nil
}

func (h *MjHero) GetTagName() string {
	return h.Spec.TagName
}

func (h *MjHero) GetStyles(element string) []*core.Style {
	switch element {
	case "div":
		return []*core.Style{
			{Name: "margin", Value: "0 auto"},
			{Name: "max-width", Value: strconv.Itoa(h.Context.ContainerWidth) + "px"},
		}
	case "table":
		return []*core.Style{
			{Name: "width", Value: "100%"},
		}
	case "tr":
		return []*core.Style{
			{Name: "vertical-align", Value: "top"},
		}
	case "td-fluid":
		return []*core.Style{
			{Name: "width", Value: "0.01%"},
			{Name: "padding-bottom", Value: h.getBackgroundRatio()},
			{Name: "mso-padding-bottom-alt", Value: "0"},
		}
	case "outlook-table":
		return []*core.Style{
			{Name: "width", Value: strconv.Itoa(h.Context.ContainerWidth) + "px"},
		}
	case "outlook-td":
		return []*core.Style{
			{Name: "line-height", Value: "0"},
			{Name: "font-size", Value: "0"},
			{Name: "mso-line-height-rule", Value: "exactly"},
		}
	case "outlook-inner-table":
		return []*core.Style{
			{Name: "width", Value: strconv.Itoa(h.Context.ContainerWidth) + "px"},
		}
	case "outlook-image":
		return []*core.Style{
			{Name: "border", Value: "0"},
			{Name: "height", Value: h.GetAttributeOr("background-height", "")},
			{Name: "mso-position-horizontal", Value: "center"},
			{Name: "position", Value: "absolute"},
			{Name: "top", Value: "0"},
			{Name: "width", Value: h.getOutlookImageWidth()},
			{Name: "z-index", Value: "-3"},
		}
	case "outlook-inner-td":
		return []*core.Style{
			{Name: "background-color", Value: h.GetAttributeOr("inner-background-color", "")},
			{Name: "padding", Value: h.GetAttributeOr("inner-padding", "")},
			{Name: "padding-top", Value: h.GetAttributeOr("inner-padding-top", "")},
			{Name: "padding-left", Value: h.GetAttributeOr("inner-padding-left", "")},
			{Name: "padding-right", Value: h.GetAttributeOr("inner-padding-right", "")},
			{Name: "padding-bottom", Value: h.GetAttributeOr("inner-padding-bottom", "")},
		}
	case "inner-table":
		return []*core.Style{
			{Name: "width", Value: "100%"},
			{Name: "margin", Value: "0px"},
		}
	case "inner-div":
		return []*core.Style{
			{Name: "background-color", Value: h.GetAttributeOr("inner-background-color", "")},
			{Name: "float", Value: h.GetAttributeOr("align", "")},
			{Name: "margin", Value: "0px auto"},
			{Name: "width", Value: h.GetAttributeOr("width", "")},
		}
	case "fluid-height":
		return []*core.Style{
			{Name: "background", Value: h.getBackground()},
			{Name: "background-position", Value: h.GetAttributeOr("background-position", "")},
			{Name: "background-repeat", Value: "no-repeat"},
			{Name: "border-radius", Value: h.GetAttributeOr("border-radius", "")},
			{Name: "padding", Value: h.GetAttributeOr("padding", "")},
			{Name: "padding-top", Value: h.GetAttributeOr("padding-top", "")},
			{Name: "padding-left", Value: h.GetAttributeOr("padding-left", "")},
			{Name: "padding-right", Value: h.GetAttributeOr("padding-right", "")},
			{Name: "padding-bottom", Value: h.GetAttributeOr("padding-bottom", "")},
			{Name: "vertical-align", Value: h.GetAttributeOr("vertical-align", "")},
		}
	case "fixed-height":
		return []*core.Style{
			{Name: "background", Value: h.getBackground()},
			{Name: "background-position", Value: h.GetAttributeOr("background-position", "")},
			{Name: "background-repeat", Value: "no-repeat"},
			{Name: "border-radius", Value: h.GetAttributeOr("border-radius", "")},
			{Name: "padding", Value: h.GetAttributeOr("padding", "")},
			{Name: "padding-top", Value: h.GetAttributeOr("padding-top", "")},
			{Name: "padding-left", Value: h.GetAttributeOr("padding-left", "")},
			{Name: "padding-right", Value: h.GetAttributeOr("padding-right", "")},
			{Name: "padding-bottom", Value: h.GetAttributeOr("padding-bottom", "")},
			{Name: "vertical-align", Value: h.GetAttributeOr("vertical-align", "")},
			{Name: "height", Value: strconv.Itoa(h.getInnerHeight()) + "px"},
		}
	default:
		return nil
	}
}

func (h *MjHero) getChildStyles(component core.Component) []*core.Style {
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

func (h *MjHero) GetChildContext() *core.MJMLContext {
	// Refactor -- removePaddingFor(width, ['padding', 'inner-padding'])
	containerWidth := h.Context.ContainerWidth
	paddingSize := h.GetShorthandAttrValue("padding", "left") +
		h.GetShorthandAttrValue("padding", "right")

	currentContainerWidth := containerWidth - paddingSize

	copy := *h.Context
	copy.ContainerWidth = currentContainerWidth

	return &copy
}

func (h *MjHero) getBackgroundRatio() string {
	backgroundHeight := h.GetAttributeOr("background-height", "")
	backgroundWidth := h.GetAttributeOr("background-width", "")

	height, err := helpers.ParseIntLoose(backgroundHeight)
	if err != nil {
		height = 0
	}

	width, err := helpers.ParseIntLoose(backgroundWidth)
	if err != nil {
		width = h.Context.ContainerWidth
	}

	ratio := (float64(height) / float64(width)) * 100

	return strconv.Itoa(int(ratio)) + "%"
}

func (h *MjHero) getInnerHeight() int {
	height, _ := helpers.ParseIntLoose(h.GetAttributeOr("height", ""))
	paddingTop := h.GetShorthandAttrValue("padding", "top")
	paddingBottom := h.GetShorthandAttrValue("padding", "bottom")

	return height - paddingTop - paddingBottom
}

func (h *MjHero) getOutlookImageWidth() string {
	width := h.GetAttribute("background-width")
	if width != nil && *width != "" {
		return *width
	}
	return strconv.Itoa(h.Context.ContainerWidth) + "px"
}

func (h *MjHero) getBackground() string {
	attrs := []string{h.GetAttributeOr("background-color", "")}
	url := h.GetAttributeOr("background-url", "")
	if url != "" {
		attrs = append(attrs, "url('"+url+"')")
		attrs = append(attrs, "no-repeat")
		attrs = append(attrs, h.GetAttributeOr("background-position", "")+" / cover")
	}

	return strings.Join(attrs, " ")
}

func (h *MjHero) renderChild(w core.MJMLWriter, child core.Component) error {
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
			AddNullable("background", child.GetAttribute("container-background-color")).
			AddNullable("class", child.GetAttribute("css-class")).
			Write(w); err != nil {
			return err
		}

		if _, err := w.WriteString(" style=\""); err != nil {
			return err
		}

		if err := FormatCssStyles(w, h.getChildStyles(child)); err != nil {
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

	return nil
}

func (h *MjHero) renderContent(w core.MJMLWriter) error {
	if _, err := w.WriteString("<!--[if mso | IE]><table "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(h).
		AddNullable("align", h.GetAttribute("align")).
		Add("border", "0").
		Add("cellpadding", "0").
		Add("cellspacing", "0").
		Add("style", "outlook-inner-table").
		Add("width", strconv.Itoa(h.Context.ContainerWidth)).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><tr><td "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(h).
		Add("style", "outlook-inner-td").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><![endif]--><div "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(h).
		AddNullable("align", h.GetAttribute("align")).
		Add("class", "mj-hero-content").
		Add("style", "inner-div").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><table "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(h).
		Add("border", "0").
		Add("cellpadding", "0").
		Add("cellspacing", "0").
		Add("role", "presentation").
		Add("style", "inner-table").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><tbody><tr><td "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(h).
		Add("style", "inner-td").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><table "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(h).
		Add("border", "0").
		Add("cellpadding", "0").
		Add("cellspacing", "0").
		Add("role", "presentation").
		Add("style", "inner-table").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><tbody>"); err != nil {
		return err
	}

	for _, child := range h.Children {
		if err := h.renderChild(w, child); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("</tbody></table></td></tr></tbody></table></div><!--[if mso | IE]></td></tr></table><![endif]-->"); err != nil {
		return err
	}

	return nil
}

func (h *MjHero) renderMode(w core.MJMLWriter) error {
	mode := h.GetAttributeOr("mode", "")

	if mode == "fluid-height" {
		if _, err := w.WriteString("<td "); err != nil {
			return err
		}

		if err := NewAttributesBuilder(h).
			Add("style", "td-fluid").
			Write(w); err != nil {
			return err
		}

		if _, err := w.WriteString("/>"); err != nil {
			return err
		}

		if _, err := w.WriteString("<td "); err != nil {
			return err
		}

		if err := NewAttributesBuilder(h).
			AddNullable("background", h.GetAttribute("background-url")).
			Add("style", "fluid-height").
			Write(w); err != nil {
			return err
		}

		if _, err := w.WriteString(">"); err != nil {
			return err
		}

		if err := h.renderContent(w); err != nil {
			return err
		}

		if _, err := w.WriteString("</td>"); err != nil {
			return err
		}

		if _, err := w.WriteString("<td "); err != nil {
			return err
		}

		if err := NewAttributesBuilder(h).
			Add("style", "td-fluid").
			Write(w); err != nil {
			return err
		}

		if _, err := w.WriteString("/>"); err != nil {
			return err
		}

		return nil
	} else {
		if _, err := w.WriteString("<td "); err != nil {
			return err
		}

		if err := NewAttributesBuilder(h).
			AddNullable("background", h.GetAttribute("background-url")).
			Add("style", "fixed-height").
			Add("height", strconv.Itoa(h.getInnerHeight())).
			Write(w); err != nil {
			return err
		}

		if _, err := w.WriteString(">"); err != nil {
			return err
		}

		if err := h.renderContent(w); err != nil {
			return err
		}

		if _, err := w.WriteString("</td>"); err != nil {
			return err
		}
	}

	return nil
}

func (h *MjHero) Render(w core.MJMLWriter) error {
	if _, err := w.WriteString("<!--[if mso | IE]><table "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(h).
		Add("align", "center").
		Add("border", "0").
		Add("cellpadding", "0").
		Add("cellspacing", "0").
		Add("role", "presentation").
		Add("style", "outlook-table").
		Add("width", strconv.Itoa(h.Context.ContainerWidth)).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><tr><td "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(h).
		Add("style", "outlook-td").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><v:image "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(h).
		Add("style", "outlook-image").
		AddNullable("src", h.GetAttribute("background-url")).
		Add("xmlns:v", "urn:schemas-microsoft-com:vml").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("/><![endif]--><div "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(h).
		AddNullable("align", h.GetAttribute("align")).
		AddNullable("class", h.GetAttribute("css-class")).
		Add("style", "div").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><table "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(h).
		Add("border", "0").
		Add("cellpadding", "0").
		Add("cellspacing", "0").
		Add("role", "presentation").
		Add("style", "table").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><tbody><tr "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(h).
		Add("style", "tr").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if err := h.renderMode(w); err != nil {
		return err
	}

	if _, err := w.WriteString("</tr></tbody></table></div><!--[if mso | IE]></td></tr></table><![endif]-->\n"); err != nil {
		return err
	}

	return nil
}
