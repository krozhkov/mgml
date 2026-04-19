package components

import (
	"fmt"
	"maps"
	"slices"

	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/core/helpers"
	"github.com/krozhkov/mgml/parser"
)

var IMG_BASE_URL = "https://www.mailjet.com/images/theme/v1/icons/ico-social/"

type SocialNetwork struct {
	Src             string
	BackgroundColor string
	ShareUrl        string
}

var defaultSocialNetworks = map[string]*SocialNetwork{
	"facebook": {
		ShareUrl:        "https://www.facebook.com/sharer/sharer.php?u=%s",
		BackgroundColor: "#3b5998",
		Src:             IMG_BASE_URL + "facebook.png",
	},
	"twitter": {
		ShareUrl:        "https://twitter.com/intent/tweet?url=%s",
		BackgroundColor: "#55acee",
		Src:             IMG_BASE_URL + "twitter.png",
	},
	"x": {
		ShareUrl:        "https://twitter.com/intent/tweet?url=%s",
		BackgroundColor: "#000000",
		Src:             IMG_BASE_URL + "twitter-x.png",
	},
	"google": {
		ShareUrl:        "https://plus.google.com/share?url=%s",
		BackgroundColor: "#dc4e41",
		Src:             IMG_BASE_URL + "google-plus.png",
	},
	"pinterest": {
		ShareUrl:        "https://pinterest.com/pin/create/button/?url=%s&media=&description=",
		BackgroundColor: "#bd081c",
		Src:             IMG_BASE_URL + "pinterest.png",
	},
	"linkedin": {
		ShareUrl:        "https://www.linkedin.com/shareArticle?mini=true&url=%s&title=&summary=&source=",
		BackgroundColor: "#0077b5",
		Src:             IMG_BASE_URL + "linkedin.png",
	},
	"instagram": {
		BackgroundColor: "#3f729b",
		Src:             IMG_BASE_URL + "instagram.png",
	},
	"web": {
		Src:             IMG_BASE_URL + "web.png",
		BackgroundColor: "#4BADE9",
	},
	"snapchat": {
		Src:             IMG_BASE_URL + "snapchat.png",
		BackgroundColor: "#FFFA54",
	},
	"youtube": {
		Src:             IMG_BASE_URL + "youtube.png",
		BackgroundColor: "#EB3323",
	},
	"tumblr": {
		Src:             IMG_BASE_URL + "tumblr.png",
		ShareUrl:        "https://www.tumblr.com/widgets/share/tool?canonicalUrl=%s",
		BackgroundColor: "#344356",
	},
	"github": {
		Src:             IMG_BASE_URL + "github.png",
		BackgroundColor: "#000000",
	},
	"xing": {
		Src:             IMG_BASE_URL + "xing.png",
		ShareUrl:        "https://www.xing.com/app/user?op=share&url=%s",
		BackgroundColor: "#296366",
	},
	"vimeo": {
		Src:             IMG_BASE_URL + "vimeo.png",
		BackgroundColor: "#53B4E7",
	},
	"medium": {
		Src:             IMG_BASE_URL + "medium.png",
		BackgroundColor: "#000000",
	},
	"soundcloud": {
		Src:             IMG_BASE_URL + "soundcloud.png",
		BackgroundColor: "#EF7F31",
	},
	"dribbble": {
		Src:             IMG_BASE_URL + "dribbble.png",
		BackgroundColor: "#D95988",
	},
}

func init() {
	keys := slices.Collect(maps.Keys(defaultSocialNetworks))

	for _, key := range keys {
		network := defaultSocialNetworks[key]
		noshareKey := key + "-noshare"
		defaultSocialNetworks[noshareKey] = &SocialNetwork{
			Src:             network.Src,
			BackgroundColor: network.BackgroundColor,
			ShareUrl:        "%s",
		}
	}
}

var MjSocialElementSpec = &core.ComponentSpec{
	TagName:   "mj-social-element",
	EndingTag: true,
	AllowedAttributes: map[string]string{
		"align":            "enum(left,center,right)",
		"icon-position":    "enum(left,right)",
		"background-color": "color",
		"color":            "color",
		"border-radius":    "unit(px)",
		"font-family":      "string",
		"font-size":        "unit(px)",
		"font-style":       "string",
		"font-weight":      "string",
		"href":             "string",
		"icon-size":        "unit(px,%)",
		"icon-height":      "unit(px,%)",
		"icon-padding":     "unit(px,%){1,4}",
		"line-height":      "unit(px,%,)",
		"name":             "string",
		"padding-bottom":   "unit(px,%)",
		"padding-left":     "unit(px,%)",
		"padding-right":    "unit(px,%)",
		"padding-top":      "unit(px,%)",
		"padding":          "unit(px,%){1,4}",
		"text-padding":     "unit(px,%){1,4}",
		"rel":              "string",
		"src":              "string",
		"srcset":           "string",
		"sizes":            "string",
		"alt":              "string",
		"title":            "string",
		"target":           "string",
		"text-decoration":  "string",
		"vertical-align":   "enum(top,middle,bottom)",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "alt", Value: ""},
		{Key: "align", Value: "left"},
		{Key: "icon-position", Value: "left"},
		{Key: "color", Value: "#000"},
		{Key: "border-radius", Value: "3px"},
		{Key: "font-family", Value: "Ubuntu, Helvetica, Arial, sans-serif"},
		{Key: "font-size", Value: "13px"},
		{Key: "line-height", Value: "1"},
		{Key: "padding", Value: "4px"},
		{Key: "text-padding", Value: "4px 4px 4px 0"},
		{Key: "target", Value: "_blank"},
		{Key: "text-decoration", Value: "none"},
		{Key: "vertical-align", Value: "middle"},
	},
	Create: NewMjSocialElement,
}

type MjSocialElement struct {
	*core.BodyComponent
}

func NewMjSocialElement(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	el := &MjSocialElement{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	el.BodyComponent.Component = el

	return el, nil
}

func (e *MjSocialElement) GetTagName() string {
	return e.Spec.TagName
}

func (e *MjSocialElement) GetStyles(element string) []*core.Style {
	attrs := e.getSocialAttributes()
	var iconSize string
	if attr, ok := attrs["icon-size"]; ok {
		iconSize = attr
	}
	var iconHeight string
	if attr, ok := attrs["icon-height"]; ok {
		iconHeight = attr
	}
	var backgroundColor string
	if attr, ok := attrs["background-color"]; ok {
		backgroundColor = attr
	}

	switch element {
	case "td":
		return []*core.Style{
			{Name: "padding", Value: e.GetAttributeOr("padding", "")},
			{Name: "padding-top", Value: e.GetAttributeOr("padding-top", "")},
			{Name: "padding-right", Value: e.GetAttributeOr("padding-right", "")},
			{Name: "padding-bottom", Value: e.GetAttributeOr("padding-bottom", "")},
			{Name: "padding-left", Value: e.GetAttributeOr("padding-left", "")},
			{Name: "vertical-align", Value: e.GetAttributeOr("vertical-align", "")},
		}
	case "table":
		return []*core.Style{
			{Name: "background", Value: backgroundColor},
			{Name: "border-radius", Value: e.GetAttributeOr("border-radius", "")},
			{Name: "width", Value: iconSize},
		}
	case "icon":
		height := iconHeight
		if height == "" {
			height = iconSize
		}
		return []*core.Style{
			{Name: "padding", Value: e.GetAttributeOr("icon-padding", "")},
			{Name: "font-size", Value: "0"},
			{Name: "height", Value: height},
			{Name: "vertical-align", Value: "middle"},
			{Name: "width", Value: iconSize},
		}
	case "img":
		return []*core.Style{
			{Name: "border-radius", Value: e.GetAttributeOr("border-radius", "")},
			{Name: "display", Value: "block"},
		}
	case "tdText":
		return []*core.Style{
			{Name: "vertical-align", Value: "middle"},
			{Name: "padding", Value: e.GetAttributeOr("text-padding", "")},
		}
	case "text":
		return []*core.Style{
			{Name: "color", Value: e.GetAttributeOr("color", "")},
			{Name: "font-size", Value: e.GetAttributeOr("font-size", "")},
			{Name: "font-weight", Value: e.GetAttributeOr("font-weight", "")},
			{Name: "font-style", Value: e.GetAttributeOr("font-style", "")},
			{Name: "font-family", Value: e.GetAttributeOr("font-family", "")},
			{Name: "line-height", Value: e.GetAttributeOr("line-height", "")},
			{Name: "text-decoration", Value: e.GetAttributeOr("text-decoration", "")},
		}
	default:
		return nil
	}
}

func (e *MjSocialElement) GetChildContext() *core.MJMLContext {
	return e.Context
}

func (e *MjSocialElement) getSocialAttributes() map[string]string {
	name := e.GetAttributeOr("name", "")
	var socialNetwork *SocialNetwork
	if sn, ok := defaultSocialNetworks[name]; ok {
		socialNetwork = sn
	}

	href := e.GetAttributeOr("href", "")

	if href != "" && socialNetwork.ShareUrl != "" {
		href = fmt.Sprintf(socialNetwork.ShareUrl, href)
	}

	attrNames := []string{
		"icon-size",
		"icon-height",
		"srcset",
		"sizes",
		"src",
		"background-color",
	}

	attrs := map[string]string{
		"href": href,
	}

	for _, attr := range attrNames {
		value := e.GetAttribute(attr)
		if value != nil {
			attrs[attr] = *value
			continue
		}

		switch attr {
		case "src":
			attrs[attr] = socialNetwork.Src
		case "background-color":
			attrs[attr] = socialNetwork.BackgroundColor
		}
	}

	return attrs
}

func (e *MjSocialElement) renderIcon(w core.MJMLWriter, socialAttrs map[string]string) error {
	var src *string
	if attr, ok := socialAttrs["src"]; ok {
		src = &attr
	}
	var srcset *string
	if attr, ok := socialAttrs["srcset"]; ok {
		srcset = &attr
	}
	var sizes *string
	if attr, ok := socialAttrs["sizes"]; ok {
		sizes = &attr
	}
	var href *string
	if attr, ok := socialAttrs["href"]; ok {
		href = &attr
	}
	var iconSize *string
	if attr, ok := socialAttrs["icon-size"]; ok {
		attr = helpers.TrimNumber(attr)
		iconSize = &attr
	}
	var iconHeight *string
	if attr, ok := socialAttrs["icon-height"]; ok {
		attr = helpers.TrimNumber(attr)
		iconHeight = &attr
	}

	hasLink := href != nil && *href != ""

	if _, err := w.WriteString("<td "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(e).
		Add("style", "td").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><table "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(e).
		Add("border", "0").
		Add("cellpadding", "0").
		Add("cellspacing", "0").
		Add("role", "presentation").
		Add("style", "table").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><tbody><tr><td "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(e).
		Add("style", "icon").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if hasLink {
		if _, err := w.WriteString("<a "); err != nil {
			return err
		}

		if err := NewAttributesBuilder(e).
			AddNullable("href", href).
			AddNullable("rel", e.GetAttribute("rel")).
			AddNullable("target", e.GetAttribute("target")).
			Write(w); err != nil {
			return err
		}

		if _, err := w.WriteString(">"); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("<img "); err != nil {
		return err
	}

	height := iconHeight
	if height == nil {
		height = iconSize
	}
	if err := NewAttributesBuilder(e).
		AddNullable("alt", e.GetAttribute("alt")).
		AddNullable("title", e.GetAttribute("title")).
		AddNullable("height", height).
		AddNullable("src", src).
		Add("style", "img").
		AddNullable("width", iconSize).
		AddNullable("sizes", sizes).
		AddNullable("srcset", srcset).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("/>"); err != nil {
		return err
	}

	if hasLink {
		if _, err := w.WriteString("</a>"); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("</td></tr></tbody></table></td>"); err != nil {
		return err
	}

	return nil
}

func (e *MjSocialElement) renderContent(w core.MJMLWriter, socialAttrs map[string]string) error {
	var href *string
	if attr, ok := socialAttrs["href"]; ok {
		href = &attr
	}

	hasLink := href != nil && *href != ""

	content := e.GetContent()

	if content == "" {
		return nil
	}

	if _, err := w.WriteString("<td "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(e).
		Add("style", "tdText").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if hasLink {
		if _, err := w.WriteString("<a "); err != nil {
			return err
		}

		if err := NewAttributesBuilder(e).
			AddNullable("href", href).
			Add("style", "text").
			AddNullable("rel", e.GetAttribute("rel")).
			AddNullable("target", e.GetAttribute("target")).
			Write(w); err != nil {
			return err
		}

		if _, err := w.WriteString(">"); err != nil {
			return err
		}
	} else {
		if _, err := w.WriteString("<span "); err != nil {
			return err
		}

		if err := NewAttributesBuilder(e).
			Add("style", "text").
			Write(w); err != nil {
			return err
		}

		if _, err := w.WriteString(">"); err != nil {
			return err
		}
	}

	if _, err := w.WriteString(content); err != nil {
		return err
	}

	if hasLink {
		if _, err := w.WriteString("</a>"); err != nil {
			return err
		}
	} else {
		if _, err := w.WriteString("</span>"); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("</td>"); err != nil {
		return err
	}

	return nil
}

func (e *MjSocialElement) Render(w core.MJMLWriter) error {
	socialAttrs := e.getSocialAttributes()

	if _, err := w.WriteString("<tr "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(e).
		AddNullable("class", e.GetAttribute("css-class")).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	iconPosition := e.GetAttributeOr("icon-position", "")

	if iconPosition == "left" {
		if err := e.renderIcon(w, socialAttrs); err != nil {
			return err
		}

		if err := e.renderContent(w, socialAttrs); err != nil {
			return err
		}
	} else {
		if err := e.renderContent(w, socialAttrs); err != nil {
			return err
		}

		if err := e.renderIcon(w, socialAttrs); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("</tr>\n"); err != nil {
		return err
	}

	return nil
}
