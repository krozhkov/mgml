package components

import (
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/parser"
)

var MjSocialSpec = &core.ComponentSpec{
	TagName: "mj-social",
	AllowedAttributes: map[string]string{
		"align":                      "enum(left,right,center)",
		"border-radius":              "unit(px,%)",
		"container-background-color": "color",
		"color":                      "color",
		"font-family":                "string",
		"font-size":                  "unit(px)",
		"font-style":                 "string",
		"font-weight":                "string",
		"icon-size":                  "unit(px,%)",
		"icon-height":                "unit(px,%)",
		"icon-padding":               "unit(px,%){1,4}",
		"inner-padding":              "unit(px,%){1,4}",
		"line-height":                "unit(px,%,)",
		"mode":                       "enum(horizontal,vertical)",
		"padding-bottom":             "unit(px,%)",
		"padding-left":               "unit(px,%)",
		"padding-right":              "unit(px,%)",
		"padding-top":                "unit(px,%)",
		"padding":                    "unit(px,%){1,4}",
		"table-layout":               "enum(auto,fixed)",
		"text-padding":               "unit(px,%){1,4}",
		"text-decoration":            "string",
		"vertical-align":             "enum(top,bottom,middle)",
	},
	DefaultAttributes: []*core.Attribute{
		{Key: "align", Value: "center"},
		{Key: "border-radius", Value: "3px"},
		{Key: "color", Value: "#333333"},
		{Key: "font-family", Value: "Ubuntu, Helvetica, Arial, sans-serif"},
		{Key: "font-size", Value: "13px"},
		{Key: "icon-size", Value: "20px"},
		{Key: "line-height", Value: "22px"},
		{Key: "mode", Value: "horizontal"},
		{Key: "padding", Value: "10px 25px"},
		{Key: "text-decoration", Value: "none"},
	},
	Create: NewMjSocial,
}

type MjSocial struct {
	*core.BodyComponent
}

func NewMjSocial(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	social := &MjSocial{
		BodyComponent: core.NewBodyComponent(node, spec, props, context),
	}

	social.BodyComponent.Component = social

	attributes := social.getSocialAttributes()
	children, err := social.CreateChildren(nil, &core.CreateChildrenOptions{
		Attributes: attributes,
	})

	if err != nil {
		return social, err
	}

	social.Children = children

	return social, nil
}

func (s *MjSocial) GetTagName() string {
	return s.Spec.TagName
}

func (s *MjSocial) GetStyles(element string) []*core.Style {
	switch element {
	case "tableVertical":
		return []*core.Style{
			{Name: "margin", Value: "0px"},
		}
	default:
		return nil
	}
}

func (s *MjSocial) GetChildContext() *core.MJMLContext {
	return s.Context
}

func (s *MjSocial) getSocialAttributes() []*core.Attribute {
	attrs := make([]*core.Attribute, 0)

	innerPadding := s.GetAttribute("inner-padding")
	if innerPadding != nil && *innerPadding != "" {
		attrs = append(attrs, &core.Attribute{Key: "padding", Value: *innerPadding})
	}

	attrNames := []string{
		"border-radius",
		"color",
		"font-family",
		"font-size",
		"font-weight",
		"font-style",
		"icon-size",
		"icon-height",
		"icon-padding",
		"text-padding",
		"line-height",
		"text-decoration",
	}

	for _, attr := range attrNames {
		value := s.GetAttribute(attr)
		if value != nil {
			attrs = append(attrs, &core.Attribute{Key: attr, Value: *value})
		}
	}

	return attrs
}

func (s *MjSocial) renderHorizontal(w core.MJMLWriter) error {
	if _, err := w.WriteString("<!--[if mso | IE]>"); err != nil {
		return err
	}

	if _, err := w.WriteString("<table "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(s).
		AddNullable("align", s.GetAttribute("align")).
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

	for _, child := range s.Children {
		if child.IsRawElement() {
			if err := child.Render(w); err != nil {
				return err
			}
		} else {
			if _, err := w.WriteString("<!--[if mso | IE]><td><![endif]--><table "); err != nil {
				return err
			}

			if err := NewAttributesBuilder(child).
				AddNullable("align", s.GetAttribute("align")).
				Add("border", "0").
				Add("cellpadding", "0").
				Add("cellspacing", "0").
				Add("role", "presentation").
				Write(w); err != nil {
				return err
			}

			if _, err := w.WriteString(" style=\""); err != nil {
				return err
			}

			if err := FormatCssStyles(w, []*core.Style{
				{Name: "float", Value: "none"},
				{Name: "display", Value: "inline-table"},
			}); err != nil {
				return err
			}

			if _, err := w.WriteString("\"><tbody>"); err != nil {
				return err
			}

			if err := child.Render(w); err != nil {
				return err
			}

			if _, err := w.WriteString("</tbody></table><!--[if mso | IE]></td><![endif]-->"); err != nil {
				return err
			}
		}
	}

	if _, err := w.WriteString("<!--[if mso | IE]></tr></table><![endif]-->\n"); err != nil {
		return err
	}

	return nil
}

func (s *MjSocial) renderVertical(w core.MJMLWriter) error {
	if _, err := w.WriteString("<table "); err != nil {
		return err
	}

	if err := NewAttributesBuilder(s).
		Add("border", "0").
		Add("cellpadding", "0").
		Add("cellspacing", "0").
		Add("role", "presentation").
		Add("style", "tableVertical").
		Write(w); err != nil {
		return err
	}

	if _, err := w.WriteString("><tbody>"); err != nil {
		return err
	}

	for _, child := range s.Children {
		if err := child.Render(w); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("</tbody></table>\n"); err != nil {
		return err
	}

	return nil
}

func (s *MjSocial) Render(w core.MJMLWriter) error {
	if s.GetAttributeOr("mode", "") == "horizontal" {
		if err := s.renderHorizontal(w); err != nil {
			return err
		}
	} else {
		if err := s.renderVertical(w); err != nil {
			return err
		}
	}

	return nil
}
