package components

import (
	"strconv"

	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/core/helpers"
	"github.com/krozhkov/mgml/parser"
)

func ptr[T any](v T) *T {
	return &v
}

var MjCarouselImageSpec = &core.ComponentSpec{
	TagName:   "mj-carousel-image",
	EndingTag: true,
	AllowedAttributes: map[string]string{
		"alt":              "string",
		"href":             "string",
		"rel":              "string",
		"target":           "string",
		"title":            "string",
		"src":              "string",
		"thumbnails-src":   "string",
		"border-radius":    "unit(px,%){1,4}",
		"tb-border":        "string",
		"tb-border-radius": "unit(px,%){1,4}",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "alt", Value: ""},
		{Key: "target", Value: "_blank"},
	},
	Create: NewMjCarouselImage,
}

type MjCarouselImage struct {
	*core.BodyComponent
}

func NewMjCarouselImage(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	image := &MjCarouselImage{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	image.BodyComponent.Component = image

	return image, nil
}

func (m *MjCarouselImage) GetTagName() string {
	return m.Spec.TagName
}

func (m *MjCarouselImage) GetStyles(element string) []*core.Style {
	switch element {
	case "images.img":
		return []*core.Style{
			{Name: "border-radius", Value: m.GetAttributeOr("border-radius", "")},
			{Name: "display", Value: "block"},
			{Name: "width", Value: strconv.Itoa(m.Context.ContainerWidth) + "px"},
			{Name: "max-width", Value: "100%"},
			{Name: "height", Value: "auto"},
		}
	case "images.firstImageDiv":
		return nil
	case "images.otherImageDiv":
		return []*core.Style{
			{Name: "display", Value: "none"},
			{Name: "mso-hide", Value: "all"},
		}
	case "radio.input":
		return []*core.Style{
			{Name: "display", Value: "none"},
			{Name: "mso-hide", Value: "all"},
		}
	case "thumbnails.a":
		display := "inline-block"
		if m.hasThumbnailsSupported() {
			display = "none"
		}

		return []*core.Style{
			{Name: "border", Value: m.GetAttributeOr("tb-border", "")},
			{Name: "border-radius", Value: m.GetAttributeOr("tb-border-radius", "")},
			{Name: "display", Value: display},
			{Name: "overflow", Value: "hidden"},
			{Name: "width", Value: helpers.TrimNumber(m.GetAttributeOr("tb-width", "")) + "px"},
		}
	case "thumbnails.img":
		return []*core.Style{
			{Name: "display", Value: "block"},
			{Name: "width", Value: "100%"},
			{Name: "height", Value: "auto"},
		}
	default:
		return nil
	}
}

func (m *MjCarouselImage) GetChildContext() *core.MJMLContext {
	return m.Context
}

func (m *MjCarouselImage) hasThumbnailsSupported() bool {
	thumbnails := m.GetAttributeOr("thumbnails", "")
	return thumbnails == "supported"
}

func (m *MjCarouselImage) renderThumbnail(w core.MJMLWriter) error {
	carouselId := m.GetAttributeOr("carouselId", "")
	cssClass := helpers.SuffixCssClasses(m.GetAttributeOr("css-class", ""), "thumbnail")
	imgIndex := strconv.Itoa(m.Props.Index + 1)
	className := "mj-carousel-thumbnail mj-carousel-" + carouselId + "-thumbnail mj-carousel-" + carouselId + "-thumbnail-" + imgIndex + " " + cssClass

	thumbnailsSrc := m.GetAttribute("thumbnails-src")
	if thumbnailsSrc == nil {
		thumbnailsSrc = m.GetAttribute("src")
	}

	if _, err := w.WriteString("<a "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(m).
		Add("style", "thumbnails.a").
		Add("href", "#"+imgIndex).
		AddNullable("target", m.GetAttribute("target")).
		Add("class", className).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><label "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(m).
		Add("for", "mj-carousel-"+carouselId+"-radio-"+imgIndex).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><img "); err != nil {
		return err
	}

	tbWidth := m.GetAttribute("tb-width")
	if tbWidth != nil && *tbWidth != "" {
		tmp := helpers.TrimNumber(*tbWidth)
		tbWidth = &tmp
	}

	if err := NewAttributesBuilder(m).
		Add("style", "thumbnails.img").
		AddNullable("src", thumbnailsSrc).
		AddNullable("alt", m.GetAttribute("alt")).
		AddNullable("width", tbWidth).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("/></label></a>"); err != nil {
		return err
	}

	return nil
}

func (m *MjCarouselImage) renderRadio(w core.MJMLWriter) error {
	index := strconv.Itoa(m.Props.Index + 1)
	carouselId := m.GetAttributeOr("carouselId", "")
	var checked *string
	if m.Props.Index == 0 {
		checked = ptr("checked")
	}

	if _, err := w.WriteString("<input "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(m).
		Add("class", "mj-carousel-radio mj-carousel-"+carouselId+"-radio mj-carousel-"+carouselId+"-radio-"+index).
		AddNullable("checked", checked).
		Add("type", "radio").
		Add("name", "mj-carousel-radio-"+carouselId).
		Add("id", "mj-carousel-"+carouselId+"-radio-"+index).
		Add("style", "radio.input").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("/>"); err != nil {
		return err
	}

	return nil
}

func (m *MjCarouselImage) renderImage(w core.MJMLWriter) error {
	containerWidth := strconv.Itoa(m.Context.ContainerWidth)

	if _, err := w.WriteString("<img "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(m).
		AddNullable("title", m.GetAttribute("title")).
		AddNullable("src", m.GetAttribute("src")).
		AddNullable("alt", m.GetAttribute("alt")).
		Add("style", "images.img").
		Add("width", containerWidth).
		Add("border", "0").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("/>"); err != nil {
		return err
	}

	return nil
}

func (m *MjCarouselImage) Render(w core.MJMLWriter) error {
	index := strconv.Itoa(m.Props.Index + 1)
	cssClass := m.GetAttributeOr("css-class", "")
	style := "images.otherImageDiv"
	if m.Props.Index == 0 {
		style = "images.firstImageDiv"
	}

	if _, err := w.WriteString("<div "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(m).
		Add("class", "mj-carousel-image mj-carousel-image-"+index+" "+cssClass).
		Add("style", style).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	href := m.GetAttribute("href")
	if href != nil && *href != "" {
		if _, err := w.WriteString("<a "); err != nil {
			return err
		}

		if err := NewAttributesBuilder(m).
			AddNullable("href", href).
			AddNullable("rel", m.GetAttribute("rel")).
			Add("target", "_blank").
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
	} else {
		if err := m.renderImage(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("</div>"); err != nil {
		return err
	}

	return nil
}
