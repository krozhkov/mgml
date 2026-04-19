package components

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/core/helpers"
	"github.com/krozhkov/mgml/parser"
)

var percentageRegex = regexp.MustCompile(`^(\d+(?:\.\d+)?)%$`)

func isPercentage(str string) bool {
	return percentageRegex.MatchString(str)
}

var MjSectionSpec = &core.ComponentSpec{
	TagName: "mj-section",
	AllowedAttributes: map[string]string{
		"background-color":      "color",
		"background-url":        "string",
		"background-repeat":     "enum(repeat,no-repeat)",
		"background-size":       "string",
		"background-position":   "string",
		"background-position-x": "string",
		"background-position-y": "string",
		"border":                "string",
		"border-bottom":         "string",
		"border-left":           "string",
		"border-radius":         "string",
		"border-right":          "string",
		"border-top":            "string",
		"direction":             "enum(ltr,rtl)",
		"full-width":            "enum(full-width,false,)",
		"padding":               "unit(px,%){1,4}",
		"padding-top":           "unit(px,%)",
		"padding-bottom":        "unit(px,%)",
		"padding-left":          "unit(px,%)",
		"padding-right":         "unit(px,%)",
		"text-align":            "enum(left,center,right)",
		"text-padding":          "unit(px,%){1,4}",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "background-repeat", Value: "repeat"},
		{Key: "background-size", Value: "auto"},
		{Key: "background-position", Value: "top center"},
		{Key: "direction", Value: "ltr"},
		{Key: "padding", Value: "20px 0"},
		{Key: "text-align", Value: "center"},
		{Key: "text-padding", Value: "4px 4px 4px 0"},
	},
	Create: NewMjSection,
}

type MjSection struct {
	*core.BodyComponent
	renderWrappedChildren func(s *MjSection, w core.MJMLWriter) error
}

func NewMjSection(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	section := &MjSection{
		BodyComponent:         core.NewBodyComponent(node, spec, props, context),
		renderWrappedChildren: renderSectionWrappedChildren,
	}

	section.BodyComponent.Component = section

	children, err := section.CreateChildren(nil, &core.CreateChildrenOptions{})
	if err != nil {
		return section, err
	}

	section.Children = children

	return section, nil
}

func (s *MjSection) GetTagName() string {
	return s.Spec.TagName
}

func (s *MjSection) GetStyles(element string) []*core.Style {
	containerWidth := s.Context.ContainerWidth
	fullWidth := s.isFullWidth()
	hasBorderRadius := s.hasBorderRadius()
	hasBackground := s.hasBackground()

	var background []*core.Style
	if hasBackground {
		background = []*core.Style{
			{Name: "background", Value: s.getBackground()},
			// background size, repeat and position has to be seperate since yahoo does not support shorthand background css property
			{Name: "background-position", Value: s.getBackgroundString()},
			{Name: "background-repeat", Value: s.GetAttributeOr("background-repeat", "")},
			{Name: "background-size", Value: s.GetAttributeOr("background-size", "")},
		}
	} else {
		background = []*core.Style{
			{Name: "background", Value: s.GetAttributeOr("background-color", "")},
			{Name: "background-color", Value: s.GetAttributeOr("background-color", "")},
		}
	}

	switch element {
	case "tableFullwidth":
		var attrs []*core.Style
		if fullWidth {
			attrs = background
		}
		attrs = append(attrs, &core.Style{Name: "width", Value: "100%"})
		return attrs
	case "table":
		var attrs []*core.Style
		if !fullWidth {
			attrs = background
		}
		attrs = append(attrs, &core.Style{Name: "width", Value: "100%"})
		if hasBorderRadius {
			attrs = append(attrs, &core.Style{Name: "border-collapse", Value: "separate"})
		}
		return attrs
	case "td":
		return []*core.Style{
			{Name: "border", Value: s.GetAttributeOr("border", "")},
			{Name: "border-bottom", Value: s.GetAttributeOr("border-bottom", "")},
			{Name: "border-left", Value: s.GetAttributeOr("border-left", "")},
			{Name: "border-right", Value: s.GetAttributeOr("border-right", "")},
			{Name: "border-top", Value: s.GetAttributeOr("border-top", "")},
			{Name: "border-radius", Value: s.GetAttributeOr("border-radius", "")},
			{Name: "direction", Value: s.GetAttributeOr("direction", "")},
			{Name: "font-size", Value: "0px"},
			{Name: "padding", Value: s.GetAttributeOr("padding", "")},
			{Name: "padding-bottom", Value: s.GetAttributeOr("padding-bottom", "")},
			{Name: "padding-left", Value: s.GetAttributeOr("padding-left", "")},
			{Name: "padding-right", Value: s.GetAttributeOr("padding-right", "")},
			{Name: "padding-top", Value: s.GetAttributeOr("padding-top", "")},
			{Name: "text-align", Value: s.GetAttributeOr("text-align", "")},
		}
	case "div":
		var attrs []*core.Style
		if !fullWidth {
			attrs = background
		}
		attrs = append(attrs, &core.Style{Name: "margin", Value: "0px auto"})
		attrs = append(attrs, &core.Style{Name: "max-width", Value: strconv.Itoa(containerWidth) + "px"})
		attrs = append(attrs, &core.Style{Name: "border-radius", Value: s.GetAttributeOr("border-radius", "")})
		if hasBorderRadius {
			attrs = append(attrs, &core.Style{Name: "overflow", Value: "hidden"})
		}
		return attrs
	case "innerDiv":
		return []*core.Style{
			{Name: "line-height", Value: "0"},
			{Name: "font-size", Value: "0"},
		}
	case "beforeSection":
		return []*core.Style{
			{Name: "width", Value: strconv.Itoa(containerWidth) + "px"},
		}
	case "vRect":
		if fullWidth {
			return []*core.Style{
				{Name: "mso-width-percent", Value: "1000"},
			}
		} else {
			return []*core.Style{
				{Name: "width", Value: strconv.Itoa(containerWidth) + "px"},
			}
		}
	default:
		return nil
	}
}

func (s *MjSection) GetChildContext() *core.MJMLContext {
	box := s.GetBoxWidths().Box

	copy := *s.Context
	copy.ContainerWidth = box

	return &copy
}

func (s *MjSection) getBackground() string {
	background := []string{s.GetAttributeOr("background-color", "")}
	if s.hasBackground() {
		background = append(background, fmt.Sprintf("url('%s')", s.GetAttributeOr("background-url", "")))
		background = append(background, s.getBackgroundString())
		background = append(background, fmt.Sprintf("/ %s", s.GetAttributeOr("background-size", "")))
		background = append(background, s.GetAttributeOr("background-repeat", ""))
	}

	var sb = new(strings.Builder)

	for _, p := range background {
		if p != "" {
			if sb.Len() > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(p)
		}
	}

	return sb.String()
}

func (s *MjSection) getBackgroundString() string {
	posX, posY := s.getBackgroundPosition()
	return posX + " " + posY
}

func (s *MjSection) getBackgroundPosition() (string, string) {
	x, y := s.parseBackgroundPosition()

	posX := s.GetAttributeOr("background-position-x", "")
	if posX == "" {
		posX = x
	}

	posY := s.GetAttributeOr("background-position-y", "")
	if posY == "" {
		posY = y
	}

	return posX, posY
}

func (s *MjSection) parseBackgroundPosition() (string, string) {
	posSplit := strings.Split(s.GetAttributeOr("background-position", ""), " ")

	if len(posSplit) == 1 {
		val := posSplit[0]
		// here we must determine if x or y was provided ; other will be center
		if val == "top" || val == "bottom" {
			return "center", val
		}

		return val, "center"
	}

	if len(posSplit) == 2 {
		// x and y can be put in any order in background-position so we need to determine that based on values
		val1 := posSplit[0]
		val2 := posSplit[1]

		if val1 == "top" || val1 == "bottom" || (val1 == "center" && (val2 == "left" || val2 == "right")) {
			return val2, val1
		}

		return val1, val2
	}

	// more than 2 values is not supported, let's treat as default value
	return "center", "top"
}

func (s *MjSection) hasBackground() bool {
	return s.GetAttributeOr("background-url", "") != ""
}

func (s *MjSection) isFullWidth() bool {
	return s.GetAttributeOr("full-width", "") == "full-width"
}

func (s *MjSection) hasBorderRadius() bool {
	return s.GetAttributeOr("border-radius", "") != ""
}

func (s *MjSection) renderBefore(w core.MJMLWriter) error {
	containerWidth := s.Context.ContainerWidth

	if _, err := w.WriteString("<!--[if mso | IE]>"); err != nil {
		return err
	}

	if _, err := w.WriteString("<table "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(s).
		Add("align", "center").
		Add("border", "0").
		Add("cellpadding", "0").
		Add("cellspacing", "0").
		AddIf("class", helpers.SuffixCssClasses(s.GetAttributeOr("css-class", ""), "outlook"), s.GetAttribute("css-class") != nil).
		Add("role", "presentation").
		Add("style", "beforeSection").
		Add("width", strconv.Itoa(containerWidth)).
		AddNullable("bgcolor", s.GetAttribute("background-color")).
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if _, err := w.WriteString("<tr>"); err != nil {
		return err
	}

	if _, err := w.WriteString("<td style=\"line-height:0px;font-size:0px;mso-line-height-rule:exactly;\">"); err != nil {
		return err
	}

	if _, err := w.WriteString("<![endif]-->"); err != nil {
		return err
	}

	return nil
}

func (s *MjSection) renderAfter(w core.MJMLWriter) error {
	if _, err := w.WriteString("<!--[if mso | IE]>"); err != nil {
		return err
	}

	if _, err := w.WriteString("</td></tr></table>"); err != nil {
		return err
	}

	if _, err := w.WriteString("<![endif]-->"); err != nil {
		return err
	}

	return nil
}

func renderSectionWrappedChildren(s *MjSection, w core.MJMLWriter) error {
	if _, err := w.WriteString("<!--[if mso | IE]>"); err != nil {
		return err
	}

	if _, err := w.WriteString("<tr>"); err != nil {
		return err
	}

	if _, err := w.WriteString("<![endif]-->"); err != nil {
		return err
	}

	for _, child := range s.Children {
		if child.IsRawElement() {
			if err := child.Render(w); err != nil {
				return err
			}
		} else {
			if _, err := w.WriteString("<!--[if mso | IE]>"); err != nil {
				return err
			}

			if _, err := w.WriteString("<td "); err != nil {
				return err
			}

			if err := NewAttributesBuilder(child).
				AddNullable("align", child.GetAttribute("align")).
				AddIf("class", helpers.SuffixCssClasses(child.GetAttributeOr("css-class", ""), "outlook"), child.GetAttribute("css-class") != nil).
				Add("style", "tdOutlook").
				Write(w); err != nil {
				return err
			}

			if _, err := w.WriteString(">"); err != nil {
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

	if _, err := w.WriteString("</tr>"); err != nil {
		return err
	}

	if _, err := w.WriteString("<![endif]-->"); err != nil {
		return err
	}

	return nil
}

func (s *MjSection) adjustPosition(coordinate, bgPosX, bgPosY string) (string, string) {
	isX := coordinate == "x"
	bgRepeat := s.GetAttributeOr("background-repeat", "") == "repeat"
	pos := bgPosY
	if isX {
		pos = bgPosX
	}
	origin := pos

	if isPercentage(pos) {
		// Should be percentage at this point
		percentageValue := percentageRegex.FindString(pos)
		value, _ := helpers.ParseIntLoose(percentageValue) // ignore error
		decimal := float64(value) / 100.0

		if bgRepeat {
			pos = strconv.FormatFloat(decimal, 'f', -1, 64)
			origin = pos
		} else {
			pos = strconv.FormatFloat((-50+decimal*100)/100, 'f', -1, 64)
			origin = pos
		}
	} else if bgRepeat {
		// top (y) or center (x)
		if isX {
			origin = "0.5"
			pos = "0.5"
		} else {
			origin = "0"
			pos = "0"
		}
	} else {
		if isX {
			origin = "0"
			pos = "0"
		} else {
			origin = "-0.5"
			pos = "-0.5"
		}
	}

	return origin, pos
}

func (s *MjSection) renderWithBackground(w core.MJMLWriter, renderContent func(w core.MJMLWriter) error) error {
	bgPosX, bgPosY := s.getBackgroundPosition()

	switch bgPosX {
	case "left":
		bgPosX = "0%"
	case "center":
		bgPosX = "50%"
	case "right":
		bgPosX = "100%"
	default:
		if !isPercentage(bgPosX) {
			bgPosX = "50%"
		}
	}
	switch bgPosY {
	case "top":
		bgPosY = "0%"
	case "center":
		bgPosY = "50%"
	case "bottom":
		bgPosY = "100%"
	default:
		if !isPercentage(bgPosY) {
			bgPosY = "0%"
		}
	}

	// this logic is different when using repeat or no-repeat
	vOriginX, vPosX := s.adjustPosition("x", bgPosX, bgPosY)
	vOriginY, vPosY := s.adjustPosition("y", bgPosX, bgPosY)

	var vSizeAttributes []*core.Attribute
	// If background size is either cover or contain, we tell VML to keep the aspect
	// and fill the entire element.
	if s.GetAttributeOr("background-size", "") == "cover" || s.GetAttributeOr("background-size", "") == "contain" {
		vSizeAttributes = append(vSizeAttributes, &core.Attribute{Key: "size", Value: "1,1"})
		if s.GetAttributeOr("background-size", "") == "cover" {
			vSizeAttributes = append(vSizeAttributes, &core.Attribute{Key: "aspect", Value: "atleast"})
		} else {
			vSizeAttributes = append(vSizeAttributes, &core.Attribute{Key: "aspect", Value: "atmost"})
		}
	} else if s.GetAttributeOr("background-size", "") != "auto" {
		bgSplit := strings.Split(s.GetAttributeOr("background-size", ""), " ")

		if len(bgSplit) == 1 {
			vSizeAttributes = append(vSizeAttributes, &core.Attribute{Key: "size", Value: bgSplit[0]})
			vSizeAttributes = append(vSizeAttributes, &core.Attribute{Key: "aspect", Value: "atmost"}) // reproduces height auto
		} else {
			vSizeAttributes = append(vSizeAttributes, &core.Attribute{Key: "size", Value: strings.Join(bgSplit, ",")})
		}
	}

	vmlType := "tile"
	if s.GetAttributeOr("background-repeat", "") == "no-repeat" {
		vmlType = "frame"
	}

	if s.GetAttributeOr("background-size", "") == "auto" {
		vmlType = "tile" // if no size provided, keep old behavior because outlook can't use original image size with "frame"
		vOriginX = "0.5"
		vPosX = "0.5"
		vOriginY = "0"
		vPosY = "0" // also ensure that images are still cropped the same way
	}

	if _, err := w.WriteString("<!--[if mso | IE]>"); err != nil {
		return err
	}

	if _, err := w.WriteString("<v:rect "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(s).
		Add("style", "vRect").
		Add("xmlns:v", "urn:schemas-microsoft-com:vml").
		Add("fill", "true").
		Add("stroke", "false").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if _, err := w.WriteString("<v:fill "); err != nil {
		return err
	}

	vFillAttrs := NewAttributesBuilder(s).
		Add("origin", vOriginX+", "+vOriginY).
		Add("position", vPosX+", "+vPosY).
		AddNullable("src", s.GetAttribute("background-url")).
		AddNullable("color", s.GetAttribute("background-color")).
		Add("type", vmlType)

	for _, attr := range vSizeAttributes {
		vFillAttrs.Add(attr.Key, attr.Value)
	}

	if err := vFillAttrs.Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(" />"); err != nil {
		return err
	}

	if _, err := w.WriteString("<v:textbox style=\"mso-fit-shape-to-text:true\" inset=\"0,0,0,0\">"); err != nil {
		return err
	}

	if _, err := w.WriteString("<![endif]-->"); err != nil {
		return err
	}

	if err := renderContent(w); err != nil {
		return err
	}

	if _, err := w.WriteString("<!--[if mso | IE]>"); err != nil {
		return err
	}

	if _, err := w.WriteString("</v:textbox></v:rect>"); err != nil {
		return err
	}

	if _, err := w.WriteString("<![endif]-->"); err != nil {
		return err
	}

	return nil
}

func (s *MjSection) renderSection(w core.MJMLWriter) error {
	hasBackground := s.hasBackground()
	fullWidth := s.isFullWidth()

	if _, err := w.WriteString("<div "); err != nil {
		return err
	}

	var cssClass *string
	if !fullWidth {
		cssClass = s.GetAttribute("css-class")
	}
	if err := NewAttributesBuilder(s).
		AddNullable("class", cssClass).
		Add("style", "div").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if hasBackground {
		if _, err := w.WriteString("<div "); err != nil {
			return err
		}

		if err := NewAttributesBuilder(s).
			Add("style", "innerDiv").
			Write(w); err != nil {
			return err
		}

		if _, err := w.WriteString(">"); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("<table "); err != nil {
		return err
	}

	var background *string
	if !fullWidth {
		background = s.GetAttribute("background-url")
	}

	if err := NewAttributesBuilder(s).
		Add("align", "center").
		AddNullable("background", background).
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

	if err := NewAttributesBuilder(s).
		Add("style", "td").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString(">"); err != nil {
		return err
	}

	if _, err := w.WriteString("<!--[if mso | IE]>"); err != nil {
		return err
	}

	if _, err := w.WriteString("<table role=\"presentation\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\">"); err != nil {
		return err
	}

	if _, err := w.WriteString("<![endif]-->"); err != nil {
		return err
	}

	if err := s.renderWrappedChildren(s, w); err != nil {
		return err
	}

	if _, err := w.WriteString("<!--[if mso | IE]>"); err != nil {
		return err
	}

	if _, err := w.WriteString("</table>"); err != nil {
		return err
	}

	if _, err := w.WriteString("<![endif]-->"); err != nil {
		return err
	}

	if _, err := w.WriteString("</td></tr></tbody></table>"); err != nil {
		return err
	}

	if hasBackground {
		if _, err := w.WriteString("</div>"); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("</div>"); err != nil {
		return err
	}

	return nil
}

func (s *MjSection) renderContent(w core.MJMLWriter) error {
	if err := s.renderBefore(w); err != nil {
		return err
	}

	if err := s.renderSection(w); err != nil {
		return err
	}

	if err := s.renderAfter(w); err != nil {
		return err
	}

	return nil
}

func (s *MjSection) renderFullWidth(w core.MJMLWriter) error {
	if _, err := w.WriteString("<table "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(s).
		Add("align", "center").
		AddNullable("class", s.GetAttribute("css-class")).
		AddNullable("background", s.GetAttribute("background-url")).
		Add("border", "0").
		Add("cellpadding", "0").
		Add("cellspacing", "0").
		Add("role", "presentation").
		Add("style", "tableFullwidth").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><tbody><tr><td>"); err != nil {
		return err
	}

	if s.hasBackground() {
		if err := s.renderWithBackground(w, s.renderContent); err != nil {
			return err
		}
	} else {
		if err := s.renderContent(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("</td></tr></tbody></table>"); err != nil {
		return err
	}

	return nil
}

func (s *MjSection) renderSimple(w core.MJMLWriter) error {
	if err := s.renderBefore(w); err != nil {
		return err
	}

	if s.hasBackground() {
		if err := s.renderWithBackground(w, s.renderSection); err != nil {
			return err
		}
	} else {
		if err := s.renderSection(w); err != nil {
			return err
		}
	}

	if err := s.renderAfter(w); err != nil {
		return err
	}

	return nil
}

func (s *MjSection) Render(w core.MJMLWriter) error {
	if s.isFullWidth() {
		if err := s.renderFullWidth(w); err != nil {
			return err
		}
	} else {
		if err := s.renderSimple(w); err != nil {
			return err
		}
	}

	return nil
}
