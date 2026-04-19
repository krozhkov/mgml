package components

import (
	"strconv"

	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/core/helpers"
	"github.com/krozhkov/mgml/parser"
)

var MjWrapperSpec = &core.ComponentSpec{
	TagName:           "mj-wrapper",
	AllowedAttributes: MjSectionSpec.AllowedAttributes,
	DefaultAttributes: MjSectionSpec.DefaultAttributes,
	Create:            NewMjWrapper,
}

func NewMjWrapper(node *parser.MJMLNode, spec *core.ComponentSpec, props *core.ComponentProps, context *core.MJMLContext) (core.Component, error) {
	wrapper := &MjSection{
		BodyComponent:         core.NewBodyComponent(node, spec, props, context),
		renderWrappedChildren: renderWrapperWrappedChildren,
	}

	wrapper.BodyComponent.Component = wrapper

	children, err := wrapper.CreateChildren(nil, &core.CreateChildrenOptions{})
	if err != nil {
		return wrapper, err
	}

	wrapper.Children = children

	return wrapper, nil
}

func renderWrapperWrappedChildren(s *MjSection, w core.MJMLWriter) error {
	containerWidth := s.Context.ContainerWidth

	for _, child := range s.Children {
		if child.IsRawElement() {
			if err := child.Render(w); err != nil {
				return err
			}
		} else {
			if _, err := w.WriteString("<!--[if mso | IE]>"); err != nil {
				return err
			}

			if _, err := w.WriteString("<tr><td "); err != nil {
				return err
			}

			if err := NewAttributesBuilder(child).
				AddNullable("align", child.GetAttribute("align")).
				AddIf("class", helpers.SuffixCssClasses(child.GetAttributeOr("css-class", ""), "outlook"), child.GetAttribute("css-class") != nil).
				Add("width", strconv.Itoa(containerWidth)+"px").
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

			if _, err := w.WriteString("</td></tr>"); err != nil {
				return err
			}

			if _, err := w.WriteString("<![endif]-->"); err != nil {
				return err
			}
		}
	}

	return nil
}
