package mjml

import (
	"github.com/krozhkov/mgml/components"
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/validator"
)

func init() {
	RegisterComponent(components.MjHeadSpec)
	RegisterComponent(components.MjAttributesSpec)
	RegisterComponent(components.MjFontSpec)
	RegisterComponent(components.MjStyleSpec)
	RegisterComponent(components.MjBreakpointSpec)
	RegisterComponent(components.MjTitleSpec)
	RegisterComponent(components.MjPreviewSpec)
	RegisterComponent(components.MjHtmlAttributesSpec)
	RegisterComponent(components.MjBodySpec)
	RegisterComponent(components.MjButtonSpec)
	RegisterComponent(components.MjTextSpec)
	RegisterComponent(components.MjColumnSpec)
	RegisterComponent(components.MjGroupSpec)
	RegisterComponent(components.MjSectionSpec)
	RegisterComponent(components.MjWrapperSpec)
	RegisterComponent(components.MjRawSpec)
	RegisterComponent(components.MjImageSpec)
	RegisterComponent(components.MjSpacerSpec)
	RegisterComponent(components.MjDividerSpec)
	RegisterComponent(components.MjSocialElementSpec)
	RegisterComponent(components.MjSocialSpec)
	RegisterComponent(components.MjTableSpec)
	RegisterComponent(components.MjNavbarLinkSpec)
	RegisterComponent(components.MjNavbarSpec)
	RegisterComponent(components.MjAccordionTitleSpec)
	RegisterComponent(components.MjAccordionTextSpec)
	RegisterComponent(components.MjAccordionElementSpec)
	RegisterComponent(components.MjAccordionSpec)
	RegisterComponent(components.MjCarouselImageSpec)
	RegisterComponent(components.MjCarouselSpec)
	RegisterComponent(components.MjHeroSpec)
}

var GlobalComponents = make(map[string]*core.ComponentSpec)

func RegisterComponent(c *core.ComponentSpec) {
	GlobalComponents[c.TagName] = c

	for name, deps := range c.Dependencies {
		validator.RegisterDependency(name, deps)
	}
}
