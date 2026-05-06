package core

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/elliotchance/orderedmap/v3"
	"github.com/krozhkov/mgml/core/helpers"
	"github.com/krozhkov/mgml/parser"
)

type Style struct {
	Name  string
	Value string
}

type Attribute = orderedmap.Element[string, string]
type Attributes = orderedmap.OrderedMap[string, string]

type ValidationLevel string

const (
	ValidationLevelStrict ValidationLevel = "strict"
	ValidationLevelSoft   ValidationLevel = "soft"
	ValidationLevelSkip   ValidationLevel = "skip"
)

type MJMLWriter interface {
	io.StringWriter
	WriteConditionalTag(s string, negation bool) (int, error)
	WriteMsoConditionalTag(s string, negation bool) (int, error)
}

type Component interface {
	GetStyles(element string) []*Style
	GetAttribute(name string) *string
	GetAttributeOr(name string, orValue string) string
	GetChildContext() *MJMLContext
	IsRawElement() bool
	Render(w MJMLWriter) error
}

type ComponentProps struct {
	Index          int
	First          bool
	Last           bool
	Sibling        int
	NonRawSiblings int
	Attributes     []*Attribute
}

type ComponentSpec struct {
	TagName           string
	EndingTag         bool
	AllowedAttributes map[string]string
	DefaultAttributes []*Attribute
	Dependencies      map[string][]string
	IsRawElement      bool
	Create            func(node *parser.MJMLNode, spec *ComponentSpec, props *ComponentProps, context *MJMLContext) (Component, error)
}

type MJMLOptions struct {
	Components        map[string]*ComponentSpec
	Dependencies      map[string][]string
	SkipElements      []string
	DefaultAttributes map[string][]*Attribute
	InlineStyles      []string
	KeepComments      bool
	IgnoreIncludes    bool
	ValidationLevel   ValidationLevel
	Data              map[string]any
}

type HtmlAttribute struct {
	Name  string
	Value string
}

type HtmlAttributes struct {
	Path       string
	Attributes []HtmlAttribute
}

type GlobalData struct {
	BackgroundColor      *string
	BeforeDoctype        string
	Breakpoint           string
	Classes              map[string]*Attributes
	DefaultAttributes    map[string]*Attributes
	HtmlAttributes       []*HtmlAttributes
	Fonts                []*helpers.FontDeclaration
	InlineStyles         []string
	HeadStyles           map[string]func(breakpoint string) string
	ComponentsHeadStyles []func(breakpoint string) string
	MediaQueries         *Attributes
	Preview              *string
	Styles               []string
	Title                *string
	ForceOWADesktop      bool
	Lang                 string
	Dir                  string
	Data                 map[string]any
}

type MJMLContext struct {
	Components     map[string]*ComponentSpec
	RawComponents  []string
	ContainerWidth int
	PrinterSupport bool
	GlobalData     *GlobalData
	Processing     func(w MJMLWriter, xml string, context *MJMLContext) error
}

func (c *MJMLContext) AddMediaQuery(className string, value float64, unit string) {
	formatted := strconv.FormatFloat(value, 'f', -1, 64)
	c.GlobalData.MediaQueries.Set(className, fmt.Sprintf("{ width:%s%s !important; max-width: %s%s; }", formatted, unit, formatted, unit))
}

func (c *MJMLContext) AddFont(name string, href string) {
	index := slices.IndexFunc(c.GlobalData.Fonts, func(f *helpers.FontDeclaration) bool { return f.Name == name })
	if index == -1 {
		c.GlobalData.Fonts = append(c.GlobalData.Fonts, &helpers.FontDeclaration{Name: name, Href: href})
	} else {
		c.GlobalData.Fonts[index].Href = href
	}
}

func (c *MJMLContext) AddPreview(preview string) {
	if preview != "" {
		c.GlobalData.Preview = &preview
	}
}

func (c *MJMLContext) AddStyle(css string, inline bool) {
	if inline {
		c.GlobalData.InlineStyles = append(c.GlobalData.InlineStyles, css)
	} else {
		c.GlobalData.Styles = append(c.GlobalData.Styles, css)
	}
}

func (c *MJMLContext) AddHtmlAttributes(path string, attributes []HtmlAttribute) {
	c.GlobalData.HtmlAttributes = append(c.GlobalData.HtmlAttributes, &HtmlAttributes{Path: path, Attributes: attributes})
}

func (c *MJMLContext) AddHeadStyle(name string, styleFunc func(breakpoint string) string) {
	c.GlobalData.HeadStyles[name] = styleFunc
}

func (c *MJMLContext) AddComponentHeadSyle(styleFunc func(breakpoint string) string) {
	c.GlobalData.ComponentsHeadStyles = append(c.GlobalData.ComponentsHeadStyles, styleFunc)
}

func (c *MJMLContext) GetAttributes(node *parser.MJMLNode) *Attributes {
	if node == nil {
		return nil
	}

	attrs := orderedmap.NewOrderedMap[string, string]()

	if attributes, ok := c.GlobalData.DefaultAttributes["mj-all"]; ok {
		for attrName, val := range attributes.AllFromFront() {
			attrs.Set(attrName, val)
		}
	}

	if attributes, ok := c.GlobalData.DefaultAttributes[node.TagName]; ok {
		for attrName, val := range attributes.AllFromFront() {
			attrs.Set(attrName, val)
		}
	}

	if node.Attributes != nil && node.Attributes.Has("mj-class") {
		mjClass := node.Attributes.GetOrDefault("mj-class", "")
		classes := strings.Split(mjClass, " ")

		for _, value := range classes {
			if value == "" {
				continue
			}

			if mjClassValues, ok := c.GlobalData.Classes[value]; ok {
				for attrName, val := range mjClassValues.AllFromFront() {
					if attrName == "css-class" && attrs.Has("css-class") {
						oldValue, _ := attrs.Get("css-class")
						attrs.Set(attrName, oldValue+" "+val)
					} else {
						attrs.Set(attrName, val)
					}
				}
			}
		}
	}

	if node.Attributes != nil {
		for attrName, val := range node.Attributes.AllFromFront() {
			if attrName == "mj-class" {
				continue
			}

			attrs.Set(attrName, val)
		}
	}

	return attrs
}
