package components

import (
	"fmt"
	"io"

	"github.com/krozhkov/mgml/core"
)

type AttributesBuilder struct {
	attrs []*core.Attribute
	comp  core.Component
}

func NewAttributesBuilder(component core.Component) *AttributesBuilder {
	return &AttributesBuilder{
		attrs: make([]*core.Attribute, 0),
		comp:  component,
	}
}

func (b *AttributesBuilder) Add(name string, value string) *AttributesBuilder {
	b.attrs = append(b.attrs, &core.Attribute{Key: name, Value: value})
	return b
}

func (b *AttributesBuilder) AddNullable(name string, value *string) *AttributesBuilder {
	if value != nil {
		b.attrs = append(b.attrs, &core.Attribute{Key: name, Value: *value})
	}
	return b
}

func (b *AttributesBuilder) AddIf(name string, value string, condition bool) *AttributesBuilder {
	if condition {
		b.attrs = append(b.attrs, &core.Attribute{Key: name, Value: value})
	}
	return b
}

func (b *AttributesBuilder) Write(w io.StringWriter) error {
	return b.htmlAttributes(w, b.attrs)
}

func (b *AttributesBuilder) htmlAttributes(w io.StringWriter, attributes []*core.Attribute) error {
	var appendSpace bool

	for _, attr := range attributes {
		name := attr.Key
		value := attr.Value

		if name == "style" {
			var stylesObject []*core.Style
			if b.comp != nil {
				stylesObject = b.comp.GetStyles(value)
			}

			if appendSpace {
				w.WriteString(" ")
			}

			if _, err := w.WriteString("style=\""); err != nil {
				return err
			}

			if err := FormatCssStyles(w, stylesObject); err != nil {
				return err
			}

			if _, err := w.WriteString("\""); err != nil {
				return err
			}

			appendSpace = true

			continue
		}

		if appendSpace {
			w.WriteString(" ")
		}

		if _, err := w.WriteString(fmt.Sprintf("%s=\"%s\"", name, value)); err != nil {
			return err
		}

		appendSpace = true
	}

	return nil
}

func FormatCssStyles(w io.StringWriter, stylesObject []*core.Style) error {
	for _, pair := range stylesObject {
		name := pair.Name
		value := pair.Value

		if value == "" {
			continue
		}

		if _, err := w.WriteString(fmt.Sprintf("%s:%s;", name, value)); err != nil {
			return err
		}
	}

	return nil
}
