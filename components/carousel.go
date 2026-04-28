package components

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/core/helpers"
	"github.com/krozhkov/mgml/parser"
)

var carouselChildrenAttributes = []string{
	"tb-border",
	"tb-border-radius",
	"border-radius",
}

var MjCarouselSpec = &core.ComponentSpec{
	TagName: "mj-carousel",
	AllowedAttributes: map[string]string{
		"align":                      "enum(left,center,right)",
		"border-radius":              "unit(px,%){1,4}",
		"container-background-color": "color",
		"icon-width":                 "unit(px,%)",
		"left-icon":                  "string",
		"padding":                    "unit(px,%){1,4}",
		"padding-top":                "unit(px,%)",
		"padding-bottom":             "unit(px,%)",
		"padding-left":               "unit(px,%)",
		"padding-right":              "unit(px,%)",
		"right-icon":                 "string",
		"thumbnails":                 "enum(visible,hidden,supported)",
		"tb-border":                  "string",
		"tb-border-radius":           "unit(px,%)",
		"tb-hover-border-color":      "color",
		"tb-selected-border-color":   "color",
		"tb-width":                   "unit(px,%)",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "align", Value: "center"},
		{Key: "border-radius", Value: "6px"},
		{Key: "icon-width", Value: "44px"},
		{Key: "left-icon", Value: "https://i.imgur.com/xTh3hln.png"},
		{Key: "right-icon", Value: "https://i.imgur.com/os7o9kz.png"},
		{Key: "thumbnails", Value: "visible"},
		{Key: "tb-border", Value: "2px solid transparent"},
		{Key: "tb-border-radius", Value: "6px"},
		{Key: "tb-hover-border-color", Value: "#fead0d"},
		{Key: "tb-selected-border-color", Value: "#ccc"},
	},
	Create: NewMjCarousel,
}

type MjCarousel struct {
	*core.BodyComponent
	carouselId string
}

func NewMjCarousel(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	carousel := &MjCarousel{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
		carouselId:    helpers.GenRandomHexString(16),
	}

	carousel.BodyComponent.Component = carousel

	childrenAttr := carousel.getChildrenAttr()
	children, err := carousel.CreateChildren(nil, &core.CreateChildrenOptions{
		Attributes: childrenAttr,
	})
	if err != nil {
		return carousel, err
	}

	carousel.Children = children

	return carousel, nil
}

func (c *MjCarousel) GetTagName() string {
	return c.Spec.TagName
}

func (c *MjCarousel) ComponentHeadStyle(breakpoint string) string {
	length := len(c.Node.Children)
	carouselId := c.carouselId
	sb := new(strings.Builder)

	sb.WriteString(`
	.mj-carousel {
      -webkit-user-select: none;
      -moz-user-select: none;
      user-select: none;
    }
	`)

	sb.WriteString(fmt.Sprintf(`
	.mj-carousel-%s-icons-cell {
      display: table-cell !important;
      width: %s !important;
    }
	`, carouselId, c.GetAttributeOr("icon-width", "")))

	sb.WriteString(`
	.mj-carousel-radio,
    .mj-carousel-next,
    .mj-carousel-previous {
      display: none !important;
    }

    .mj-carousel-thumbnail,
    .mj-carousel-next,
    .mj-carousel-previous {
      touch-action: manipulation;
    }
	`)

	for i := 0; i < length; i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(".mj-carousel-" + carouselId + "-radio:checked " + strings.Repeat("+ * ", i) + "+ .mj-carousel-content .mj-carousel-image")
	}

	sb.WriteString(` {
      display: none !important;
    }
	`)

	for i := 0; i < length; i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(".mj-carousel-" + carouselId + "-radio-" + strconv.Itoa(i+1) + ":checked " + strings.Repeat("+ * ", length-i-1) + "+ .mj-carousel-content .mj-carousel-image-" + strconv.Itoa(i+1))
	}

	sb.WriteString(` {
      display: block !important;
    }
	.mj-carousel-previous-icons,
    .mj-carousel-next-icons,
	`)

	for i := 0; i < length; i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(".mj-carousel-" + carouselId + "-radio-" + strconv.Itoa(i+1) + ":checked " + strings.Repeat("+ * ", length-i-1) + "+ .mj-carousel-content .mj-carousel-next-" + strconv.Itoa(((i+(1%length)+length)%length)+1))
	}

	for i := 0; i < length; i++ {
		sb.WriteString(",")
		sb.WriteString(".mj-carousel-" + carouselId + "-radio-" + strconv.Itoa(i+1) + ":checked " + strings.Repeat("+ * ", length-i-1) + "+ .mj-carousel-content .mj-carousel-previous-" + strconv.Itoa(((i-(1%length)+length)%length)+1))
	}

	sb.WriteString(` {
      display: block !important;
    }
	`)

	for i := 0; i < length; i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(".mj-carousel-" + carouselId + "-radio-" + strconv.Itoa(i+1) + ":checked " + strings.Repeat("+ * ", length-i-1) + "+ .mj-carousel-content .mj-carousel-" + carouselId + "-thumbnail-" + strconv.Itoa(i+1))
	}

	sb.WriteString(fmt.Sprintf(` {
      border-color: %s !important;
    }
	`, c.GetAttributeOr("tb-selected-border-color", "")))

	for i := 0; i < length; i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(".mj-carousel-" + carouselId + "-radio-" + strconv.Itoa(i+1) + ":checked " + strings.Repeat("+ * ", length-i-1) + "+ .mj-carousel-content .mj-carousel-" + carouselId + "-thumbnail")
	}

	sb.WriteString(` {
      display: inline-block !important;
    }
	`)

	sb.WriteString(`
	.mj-carousel-image img + div,
    .mj-carousel-thumbnail img + div {
      display: none !important;
    }
	`)

	for i := 0; i < length; i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(".mj-carousel-" + carouselId + "-thumbnail:hover " + strings.Repeat("+ * ", length-i-1) + "+ .mj-carousel-main .mj-carousel-image")
	}

	sb.WriteString(` {
      display: none !important;
    }
	`)

	sb.WriteString(fmt.Sprintf(`
	.mj-carousel-thumbnail:hover {
      border-color: %s !important;
    }
	`, c.GetAttributeOr("tb-hover-border-color", "")))

	for i := 0; i < length; i++ {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(".mj-carousel-" + carouselId + "-thumbnail-" + strconv.Itoa(i+1) + ":hover " + strings.Repeat("+ * ", length-i-1) + "+ .mj-carousel-main .mj-carousel-image-" + strconv.Itoa(i+1))
	}

	sb.WriteString(` {
      display: block !important;
    }

	.mj-carousel noinput { display:block !important; }
	.mj-carousel noinput .mj-carousel-image-1 { display: block !important;  }
	.mj-carousel noinput .mj-carousel-arrows,
	.mj-carousel noinput .mj-carousel-thumbnails { display: none !important; }

	[owa] .mj-carousel-thumbnail { display: none !important; }
	`)

	sb.WriteString(`
	@media screen yahoo {
          .mj-carousel-` + carouselId + `-icons-cell,
          .mj-carousel-previous-icons,
          .mj-carousel-next-icons {
              display: none !important;
          }
	`)

	sb.WriteString(".mj-carousel-" + carouselId + "-radio-1:checked " + strings.Repeat("+ *", length-1) + "+ .mj-carousel-content .mj-carousel-" + carouselId + "-thumbnail-1")

	sb.WriteString(` {
              border-color: transparent;
          }
      }
	`)

	return sb.String()
}

func (c *MjCarousel) GetStyles(element string) []*core.Style {
	switch element {
	case "carousel.div":
		return []*core.Style{
			{Name: "display", Value: "table"},
			{Name: "width", Value: "100%"},
			{Name: "table-layout", Value: "fixed"},
			{Name: "text-align", Value: "center"},
			{Name: "font-size", Value: "0px"},
		}
	case "carousel.table":
		return []*core.Style{
			{Name: "caption-side", Value: "top"},
			{Name: "display", Value: "table-caption"},
			{Name: "table-layout", Value: "fixed"},
			{Name: "width", Value: "100%"},
		}
	case "images.td":
		return []*core.Style{
			{Name: "padding", Value: "0px"},
		}
	case "controls.div":
		return []*core.Style{
			{Name: "display", Value: "none"},
			{Name: "mso-hide", Value: "all"},
		}
	case "controls.img":
		return []*core.Style{
			{Name: "display", Value: "block"},
			{Name: "width", Value: c.GetAttributeOr("icon-width", "")},
			{Name: "height", Value: "auto"},
		}
	case "controls.td":
		return []*core.Style{
			{Name: "font-size", Value: "0px"},
			{Name: "display", Value: "none"},
			{Name: "mso-hide", Value: "all"},
			{Name: "padding", Value: "0px"},
		}
	default:
		return nil
	}
}

func (c *MjCarousel) GetChildContext() *core.MJMLContext {
	return c.Context
}

func (c *MjCarousel) getChildrenAttr() []*core.Attribute {
	attr := make([]*core.Attribute, 0, len(carouselChildrenAttributes)+2) // + carouselId and thumbnailsWidth

	for _, name := range carouselChildrenAttributes {
		value := c.GetAttribute(name)
		if value != nil {
			attr = append(attr, &core.Attribute{Key: name, Value: *value})
		}
	}

	attr = append(attr, &core.Attribute{Key: "carouselId", Value: c.carouselId})
	attr = append(attr, &core.Attribute{Key: "tb-width", Value: c.thumbnailsWidth()})

	thumbnails := c.GetAttribute("thumbnails")
	if thumbnails != nil {
		attr = append(attr, &core.Attribute{Key: "thumbnails", Value: *thumbnails})
	}

	return attr
}

func (c *MjCarousel) thumbnailsWidth() string {
	childrenLength := len(c.Node.Children)
	if childrenLength == 0 {
		return "0"
	}

	tbWidth := c.GetAttribute("tb-width")
	if tbWidth != nil && *tbWidth != "" {
		return *tbWidth
	}

	width := math.Min(float64(c.Context.ContainerWidth)/float64(childrenLength), 110.0)

	return strconv.FormatFloat(width, 'f', -1, 64)
}

func (c *MjCarousel) generateRadios(w core.MJMLWriter) error {
	for _, child := range c.Children {
		if upgrade, ok := child.(interface {
			renderRadio(w core.MJMLWriter) error
		}); ok {
			if err := upgrade.renderRadio(w); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *MjCarousel) generateThumbnails(w core.MJMLWriter) error {
	thumbnails := c.GetAttributeOr("thumbnails", "")
	if thumbnails != "visible" && thumbnails != "supported" {
		return nil
	}

	for _, child := range c.Children {
		if upgrade, ok := child.(interface {
			renderThumbnail(w core.MJMLWriter) error
		}); ok {
			if err := upgrade.renderThumbnail(w); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *MjCarousel) generateControls(w core.MJMLWriter, direction string, icon *string) error {
	length := len(c.Node.Children)
	iconWidth, err := helpers.ParseIntLoose(c.GetAttributeOr("icon-width", ""))
	if err != nil {
		return err
	}

	if _, err := w.WriteString("<td "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(c).
		Add("class", "mj-carousel-"+c.carouselId+"-icons-cell").
		Add("style", "controls.td").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><div "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(c).
		Add("class", "mj-carousel-"+direction+"-icons").
		Add("style", "controls.div").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	for i := 0; i < length; i++ {
		if _, err := w.WriteString("<label "); err != nil {
			return err
		}

		if err := NewAttributesBuilder(c).
			Add("for", "mj-carousel-"+c.carouselId+"-radio-"+strconv.Itoa(i+1)).
			Add("class", "mj-carousel-"+direction+" mj-carousel-"+direction+"-"+strconv.Itoa(i+1)).
			Write(w); err != nil {
			return err
		}

		if _, err := w.WriteString("><img "); err != nil {
			return err
		}

		if err := NewAttributesBuilder(c).
			AddNullable("src", icon).
			Add("alt", direction).
			Add("style", "controls.img").
			Add("width", strconv.Itoa(iconWidth)).
			Write(w); err != nil {
			return err
		}

		if _, err := w.WriteString("/></label>"); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("</div></td>"); err != nil {
		return err
	}

	return nil
}

func (c *MjCarousel) generateImages(w core.MJMLWriter) error {
	if _, err := w.WriteString("<td "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(c).
		Add("style", "images.td").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><div "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(c).
		Add("class", "mj-carousel-images").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	for _, child := range c.Children {
		if err := child.Render(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("</div></td>"); err != nil {
		return err
	}

	return nil
}

func (c *MjCarousel) generateCarousel(w core.MJMLWriter) error {
	if _, err := w.WriteString("<table "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(c).
		Add("style", "carousel.table").
		Add("border", "0").
		Add("cellpadding", "0").
		Add("cellspacing", "0").
		Add("width", "100%").
		Add("role", "presentation").
		Add("class", "mj-carousel-main").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><tbody><tr>"); err != nil {
		return err
	}

	if err := c.generateControls(w, "previous", c.GetAttribute("left-icon")); err != nil {
		return err
	}

	if err := c.generateImages(w); err != nil {
		return err
	}

	if err := c.generateControls(w, "next", c.GetAttribute("right-icon")); err != nil {
		return err
	}

	if _, err := w.WriteString("</tr></tbody></table>"); err != nil {
		return err
	}

	return nil
}

func (c *MjCarousel) renderFallback(w core.MJMLWriter) error {
	if len(c.Children) == 0 {
		return nil
	}

	child := c.Children[0]

	if _, err := w.WriteString("<!--[if mso]>"); err != nil {
		return err
	}

	if err := child.Render(w); err != nil {
		return err
	}

	if _, err := w.WriteString("<![endif]-->\n"); err != nil {
		return err
	}

	return nil
}

func (c *MjCarousel) Render(w core.MJMLWriter) error {
	if _, err := w.WriteString("<!--[if !mso]><!-->"); err != nil {
		return err
	}

	if _, err := w.WriteString("<div "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(c).
		Add("class", "mj-carousel").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if err := c.generateRadios(w); err != nil {
		return err
	}

	if _, err := w.WriteString("<div "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(c).
		Add("class", "mj-carousel-content mj-carousel-"+c.carouselId+"-content").
		Add("style", "carousel.div").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if err := c.generateThumbnails(w); err != nil {
		return err
	}

	if err := c.generateCarousel(w); err != nil {
		return err
	}

	if _, err := w.WriteString("</div></div>"); err != nil {
		return err
	}

	if _, err := w.WriteString("<!--<![endif]-->"); err != nil {
		return err
	}

	if err := c.renderFallback(w); err != nil {
		return err
	}

	return nil
}
