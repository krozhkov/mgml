package components

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/core/helpers"
	"github.com/krozhkov/mgml/parser"
)

var MjImageSpec = &core.ComponentSpec{
	TagName: "mj-image",
	AllowedAttributes: map[string]string{
		"alt":                        "string",
		"href":                       "string",
		"name":                       "string",
		"src":                        "string",
		"srcset":                     "string",
		"sizes":                      "string",
		"title":                      "string",
		"rel":                        "string",
		"align":                      "enum(left,center,right)",
		"border":                     "string",
		"border-bottom":              "string",
		"border-left":                "string",
		"border-right":               "string",
		"border-top":                 "string",
		"border-radius":              "unit(px,%){1,4}",
		"container-background-color": "color",
		"fluid-on-mobile":            "boolean",
		"padding":                    "unit(px,%){1,4}",
		"padding-bottom":             "unit(px,%)",
		"padding-left":               "unit(px,%)",
		"padding-right":              "unit(px,%)",
		"padding-top":                "unit(px,%)",
		"target":                     "string",
		"width":                      "unit(px)",
		"height":                     "unit(px,auto)",
		"max-height":                 "unit(px,%)",
		"font-size":                  "unit(px)",
		"usemap":                     "string",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "alt", Value: ""},
		{Key: "align", Value: "center"},
		{Key: "border", Value: "0"},
		{Key: "height", Value: "auto"},
		{Key: "padding", Value: "10px 25px"},
		{Key: "target", Value: "_blank"},
		{Key: "font-size", Value: "13px"},
	},
	Create: NewMjImage,
}

type MjImage struct {
	*core.BodyComponent
}

func NewMjImage(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	image := &MjImage{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	image.BodyComponent.Component = image

	return image, nil
}

func (m *MjImage) GetTagName() string {
	return m.Spec.TagName
}

func (m *MjImage) HeadStyle(breakpoint string) string {
	return fmt.Sprintf(`
    @media only screen and (max-width:%s) {
      table.mj-full-width-mobile { width: 100%% !important; }
      td.mj-full-width-mobile { width: auto !important; }
    }
	`,
		helpers.MakeLowerBreakpoint(breakpoint),
	)
}

func (m *MjImage) GetStyles(element string) []*core.Style {
	width := m.getContentWidth()
	fullWidth := m.GetAttributeOr("full-width", "") == "full-width"

	switch element {
	case "img":
		attrs := []*core.Style{
			{Name: "border", Value: m.GetAttributeOr("border", "")},
			{Name: "border-left", Value: m.GetAttributeOr("border-left", "")},
			{Name: "border-right", Value: m.GetAttributeOr("border-right", "")},
			{Name: "border-top", Value: m.GetAttributeOr("border-top", "")},
			{Name: "border-bottom", Value: m.GetAttributeOr("border-bottom", "")},
			{Name: "border-radius", Value: m.GetAttributeOr("border-radius", "")},
			{Name: "display", Value: "block"},
			{Name: "outline", Value: "none"},
			{Name: "text-decoration", Value: "none"},
			{Name: "height", Value: m.GetAttributeOr("height", "")},
			{Name: "max-height", Value: m.GetAttributeOr("max-height", "")},
		}
		if fullWidth {
			attrs = append(attrs, &core.Style{Name: "min-width", Value: "100%"})
		}
		attrs = append(attrs, &core.Style{Name: "width", Value: "100%"})
		if fullWidth {
			attrs = append(attrs, &core.Style{Name: "max-width", Value: "100%"})
		}
		attrs = append(attrs, &core.Style{Name: "font-size", Value: m.GetAttributeOr("font-size", "")})
		return attrs
	case "td":
		if !fullWidth {
			return []*core.Style{
				{Name: "width", Value: strconv.Itoa(width) + "px"},
			}
		}
		return nil
	case "table":
		attrs := []*core.Style{}
		if fullWidth {
			attrs = append(attrs, &core.Style{Name: "min-width", Value: "100%"})
		}
		if fullWidth {
			attrs = append(attrs, &core.Style{Name: "max-width", Value: "100%"})
		}
		if fullWidth {
			attrs = append(attrs, &core.Style{Name: "width", Value: strconv.Itoa(width) + "px"})
		}
		attrs = append(attrs, &core.Style{Name: "border-collapse", Value: "collapse"})
		attrs = append(attrs, &core.Style{Name: "border-spacing", Value: "0px"})
		return attrs
	default:
		return nil
	}
}

func (m *MjImage) GetChildContext() *core.MJMLContext {
	return m.Context
}

func (m *MjImage) getContentWidth() int {
	var width int = -1
	if m.Attributes.Has("width") {
		var err error
		width, err = helpers.ParseIntLoose(m.GetAttributeOr("width", ""))
		if err != nil {
			width = -1
		}
	}

	box := m.GetBoxWidths().Box

	if width == -1 || box < width {
		return box
	}

	return width
}

func (m *MjImage) renderImage(w io.StringWriter) error {
	height := m.GetAttributeOr("height", "")
	if height != "auto" && height != "" {
		height = helpers.TrimNumber(height)
	}

	if _, err := w.WriteString("<img "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(m).
		AddNullable("alt", m.GetAttribute("alt")).
		AddNullable("src", m.GetAttribute("src")).
		AddNullable("srcset", m.GetAttribute("srcset")).
		AddNullable("sizes", m.GetAttribute("sizes")).
		Add("style", "img").
		AddNullable("title", m.GetAttribute("title")).
		Add("width", strconv.Itoa(m.getContentWidth())).
		AddNullable("usemap", m.GetAttribute("usemap")).
		Add("height", height).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(" />"); err != nil {
		return err
	}

	return nil
}

func (m *MjImage) renderWrappedImage(w io.StringWriter) error {
	href := m.GetAttribute("href")

	if href != nil && *href != "" {
		if _, err := w.WriteString("<a "); err != nil {
			return err
		}

		if err := NewAttributesBuilder(m).
			AddNullable("href", href).
			AddNullable("target", m.GetAttribute("target")).
			AddNullable("rel", m.GetAttribute("rel")).
			AddNullable("name", m.GetAttribute("name")).
			AddNullable("title", m.GetAttribute("title")).
			Write(w); err != nil {
			return err
		}

		if _, err := w.WriteString(">"); err != nil {
			return err
		}

		if err := m.renderImage(w); err != nil {
			return err
		}

		if _, err := w.WriteString("</a>"); err != nil {
			return err
		}

		return nil
	}

	if err := m.renderImage(w); err != nil {
		return err
	}

	return nil
}

func (m *MjImage) Render(w core.MJMLWriter) error {
	if _, err := w.WriteString("<table "); err != nil {
		return err
	}

	fluidOnMobile := m.GetAttributeOr("fluid-on-mobile", "")

	if err := NewAttributesBuilder(m).
		Add("border", "0").
		Add("cellpadding", "0").
		Add("cellspacing", "0").
		Add("role", "presentation").
		Add("style", "table").
		AddIf("class", "mj-full-width-mobile", fluidOnMobile != "" && !strings.EqualFold(fluidOnMobile, "false")).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><tbody><tr><td "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(m).
		Add("style", "td").
		AddIf("class", "mj-full-width-mobile", fluidOnMobile != "" && !strings.EqualFold(fluidOnMobile, "false")).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if err := m.renderWrappedImage(w); err != nil {
		return err
	}

	if _, err := w.WriteString("</td></tr></tbody></table>\n"); err != nil {
		return err
	}

	return nil
}
