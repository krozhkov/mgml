package core

import (
	"fmt"
	"slices"
	"strings"

	"github.com/elliotchance/orderedmap/v3"
	"github.com/krozhkov/mgml/core/helpers"
	"github.com/krozhkov/mgml/internal/utils"
	"github.com/krozhkov/mgml/parser"
)

func InitComponent(node *parser.MJMLNode, props *ComponentProps, context *MJMLContext) (Component, error) {
	spec := context.Components[node.TagName]

	if spec != nil {
		component, err := spec.Create(node, spec, props, context)

		if err != nil {
			return nil, err
		}

		if upgrade, ok := component.(interface {
			HeadStyle(breakpoint string) string
		}); ok {
			context.AddHeadStyle(node.TagName, upgrade.HeadStyle)
		}

		if upgrade, ok := component.(interface {
			ComponentHeadStyle(breakpoint string) string
		}); ok {
			context.AddComponentHeadSyle(upgrade.ComponentHeadStyle)
		}

		return component, nil
	}

	return nil, fmt.Errorf("no matching component for tag : %s", node.TagName)
}

type BaseComponent struct {
	Node    *parser.MJMLNode
	Spec    *ComponentSpec
	Props   *ComponentProps
	Context *MJMLContext

	Attributes *Attributes
	Children   []Component
}

func NewBaseComponent(node *parser.MJMLNode, spec *ComponentSpec, props *ComponentProps, context *MJMLContext) *BaseComponent {
	attributes := helpers.FormatAttributes(
		toOrderedMap(spec.DefaultAttributes, props.Attributes, utils.ToEntries(node.Attributes)),
		spec.AllowedAttributes,
	)

	return &BaseComponent{
		Node:       node,
		Spec:       spec,
		Props:      props,
		Context:    context,
		Attributes: attributes,
	}
}

func (c *BaseComponent) GetAttribute(name string) *string {
	if c.Attributes != nil && c.Attributes.Has(name) {
		if attr, ok := c.Attributes.Get(name); ok {
			return &attr
		}
	}

	return nil
}

func (c *BaseComponent) GetAttributeOr(name string, orValue string) string {
	if c.Attributes != nil && c.Attributes.Has(name) {
		if attr, ok := c.Attributes.Get(name); ok {
			return attr
		}
	}

	return orValue
}

func (c *BaseComponent) GetContent() string {
	return strings.TrimSpace(c.Node.Content)
}

func (c *BaseComponent) GetChildren() []Component {
	return c.Children
}

func (c *BaseComponent) IsRawElement() bool {
	return c.Spec.IsRawElement
}

func (c *BaseComponent) RenderMJML(w MJMLWriter, xml string) error {
	return c.Context.Processing(w, xml, c.Context)
}

type BodyComponent struct {
	*BaseComponent
	Component Component
}

func NewBodyComponent(node *parser.MJMLNode, spec *ComponentSpec, props *ComponentProps, context *MJMLContext) *BodyComponent {
	return &BodyComponent{
		BaseComponent: NewBaseComponent(node, spec, props, context),
	}
}

func (c *BodyComponent) GetShorthandAttrValue(attribute string, direction string) int {
	mjAttributeDirection := c.GetAttribute(attribute + "-" + direction)
	mjAttribute := c.GetAttribute(attribute)

	if mjAttributeDirection != nil && *mjAttributeDirection != "" {
		val, err := helpers.ParseIntLoose(*mjAttributeDirection)
		if err != nil {
			return 0
		}

		return val
	}

	if mjAttribute == nil || *mjAttribute == "" {
		return 0
	}

	return helpers.ShorthandParser(*mjAttribute, direction)
}

func (c *BodyComponent) GetShorthandBorderValue(attribute string, direction string) int {
	if attribute == "" {
		attribute = "border"
	}

	var borderDirection string
	if direction != "" {
		attr := c.GetAttribute(attribute + "-" + direction)
		if attr != nil {
			borderDirection = *attr
		}
	}

	border := c.GetAttribute(attribute)

	if borderDirection != "" {
		return helpers.BorderParser(borderDirection)
	} else if border != nil && *border != "" {
		return helpers.BorderParser(*border)
	} else {
		return 0
	}
}

type BoxWidths struct {
	TotalWidth int
	Borders    int
	Paddings   int
	Box        int
}

func (c *BodyComponent) GetBoxWidths() BoxWidths {
	context := c.Context

	containerWidth := context.ContainerWidth

	paddings := c.GetShorthandAttrValue("padding", "right") +
		c.GetShorthandAttrValue("padding", "left")

	borders := c.GetShorthandBorderValue("border", "right") +
		c.GetShorthandBorderValue("border", "left")

	return BoxWidths{
		TotalWidth: containerWidth,
		Borders:    borders,
		Paddings:   paddings,
		Box:        containerWidth - paddings - borders,
	}
}

type CreateChildrenOptions struct {
	Attributes []*Attribute
}

func (c *BodyComponent) CreateChildren(children []*parser.MJMLNode, options *CreateChildrenOptions) ([]Component, error) {
	if children == nil {
		children = c.Node.Children
	}

	context := c.Component.GetChildContext()
	sibling := len(children)

	nonRawSiblings := len(utils.FilterFunc(children, func(node *parser.MJMLNode) bool { return !slices.Contains(context.RawComponents, node.TagName) }))

	var components = make([]Component, 0, len(children))

	for index, node := range children {
		attributes := context.GetAttributes(node)
		for _, el := range options.Attributes {
			if !attributes.Has(el.Key) {
				attributes.Set(el.Key, el.Value)
			}
		}

		component, err := InitComponent(node, &ComponentProps{
			Index:          index,
			First:          index == 0,
			Last:           (index + 1) == sibling,
			Sibling:        sibling,
			NonRawSiblings: nonRawSiblings,
			Attributes:     utils.ToEntries(attributes),
		}, context)

		if err != nil {
			return nil, err
		}

		if component == nil {
			return nil, fmt.Errorf("no matching component for tag : %s", node.TagName)
		}

		components = append(components, component)
	}

	return components, nil
}

type HeadComponent struct {
	*BaseComponent
	Component Component
}

func NewHeadComponent(node *parser.MJMLNode, spec *ComponentSpec, props *ComponentProps, context *MJMLContext) *HeadComponent {
	return &HeadComponent{
		BaseComponent: NewBaseComponent(node, spec, props, context),
	}
}

func (c *HeadComponent) CreateChildren() ([]Component, error) {
	context := c.Component.GetChildContext()
	children := c.Node.Children
	sibling := len(children)

	var components = make([]Component, 0, len(children))

	for index, node := range children {
		component, err := InitComponent(node, &ComponentProps{
			Index:   index,
			First:   index == 0,
			Last:    (index + 1) == sibling,
			Sibling: sibling,
		}, context)

		if err != nil {
			return nil, err
		}

		if component == nil {
			return nil, fmt.Errorf("no matching component for tag : %s", node.TagName)
		}

		components = append(components, component)
	}

	return components, nil
}

func toOrderedMap[K comparable, V any](sources ...[]*orderedmap.Element[K, V]) *orderedmap.OrderedMap[K, V] {
	dst := orderedmap.NewOrderedMap[K, V]()

	for _, src := range sources {
		for _, el := range src {
			dst.Set(el.Key, el.Value)
		}
	}

	return dst
}
