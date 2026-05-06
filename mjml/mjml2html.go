package mjml

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
	"unsafe"

	"github.com/elliotchance/orderedmap/v3"
	"github.com/krozhkov/go-css-select/query"
	"github.com/krozhkov/go-htmlparser2/dom"
	htmlparser2 "github.com/krozhkov/go-htmlparser2/parser"
	"github.com/krozhkov/go-htmlparser2/serializer"
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/core/helpers"
	"github.com/krozhkov/mgml/cssparser"
	"github.com/krozhkov/mgml/internal/utils"
	"github.com/krozhkov/mgml/parser"
	"github.com/krozhkov/mgml/validator"
	"github.com/krozhkov/mgml/validator/rules"
)

func MJML2Html(mjml string, options *core.MJMLOptions) (string, error) {
	if options == nil {
		options = &core.MJMLOptions{}
	}

	fonts := []*helpers.FontDeclaration{
		{Name: "Open Sans", Href: "https://fonts.googleapis.com/css?family=Open+Sans:300,400,500,700"},
		{Name: "Droid Sans", Href: "https://fonts.googleapis.com/css?family=Droid+Sans:300,400,500,700"},
		{Name: "Lato", Href: "https://fonts.googleapis.com/css?family=Lato:300,400,500,700"},
		{Name: "Roboto", Href: "https://fonts.googleapis.com/css?family=Roboto:300,400,500,700"},
		{Name: "Ubuntu", Href: "https://fonts.googleapis.com/css?family=Ubuntu:300,400,500,700"},
	}
	keepComments := options.KeepComments
	ignoreIncludes := options.IgnoreIncludes
	validationLevel := options.ValidationLevel
	data := options.Data
	if data == nil {
		data = make(map[string]any)
	}

	defaultStyles := options.InlineStyles
	defaultAttributes := make(map[string]*orderedmap.OrderedMap[string, string])
	for name, attrs := range options.DefaultAttributes {
		defaultAttributes[name] = orderedmap.NewOrderedMapWithElements(attrs...)
	}

	components := make(map[string]*core.ComponentSpec)
	maps.Copy(components, GlobalComponents)
	maps.Copy(components, options.Components)

	endingComponents := make([]parser.MJMLComponent, 0, len(components))
	for _, c := range components {
		endingComponents = append(endingComponents, parser.MJMLComponent{TagName: c.TagName, EndingTag: c.EndingTag})
	}

	node, err := parser.MJMLParser(unsafe.Slice(unsafe.StringData(mjml), len(mjml)), parser.MJMLParserOptions{
		Components:     endingComponents,
		KeepComments:   keepComments,
		FilePath:       ".",
		ActualPath:     ".",
		IgnoreIncludes: ignoreIncludes,
	}, nil)

	if err != nil {
		return "", err
	}

	var lang = "und"
	if node.Attributes != nil && node.Attributes.Has("lang") {
		lang = node.Attributes.GetOrDefault("lang", "")
	}

	var dir = "auto"
	if node.Attributes != nil && node.Attributes.Has("dir") {
		dir = node.Attributes.GetOrDefault("dir", "")
	}

	var owa = "mobile"
	if node.Attributes != nil && node.Attributes.Has("owa") {
		owa = node.Attributes.GetOrDefault("owa", "")
	}

	validatorOptions := &core.MJMLOptions{
		Components:   components,
		Dependencies: options.Dependencies,
		SkipElements: options.SkipElements,
	}

	var ruleErrors []*rules.RuleError
	switch validationLevel {
	case core.ValidationLevelSkip:
		break
	case core.ValidationLevelStrict:
		ruleErrors = validator.MJMLValidator(node, validatorOptions)
		if len(ruleErrors) > 0 {
			formattedMessages := utils.MapFunc(ruleErrors, func(e *rules.RuleError) string { return e.FormattedMessage })
			return "", fmt.Errorf("ValidationError: \n %s", strings.Join(formattedMessages, "\n"))
		}
	case core.ValidationLevelSoft:
		ruleErrors = validator.MJMLValidator(node, validatorOptions)
	default:
		ruleErrors = validator.MJMLValidator(node, validatorOptions)
	}

	bodyIndex := slices.IndexFunc(node.Children, func(n *parser.MJMLNode) bool { return n.TagName == "mj-body" })
	headIndex := slices.IndexFunc(node.Children, func(n *parser.MJMLNode) bool { return n.TagName == "mj-head" })
	var mjBody *parser.MJMLNode
	if bodyIndex != -1 {
		mjBody = node.Children[bodyIndex]
	}
	var mjHead *parser.MJMLNode
	if headIndex != -1 {
		mjHead = node.Children[headIndex]
	}

	//mjOutsideRaws := utils.FilterFunc(node.Children, func(n *parser.MJMLNode) bool { return n.TagName == "mj-raw" })

	rawComponents := utils.MapFunc(
		utils.FilterFunc(slices.Collect(maps.Values(components)), func(s *core.ComponentSpec) bool { return s.IsRawElement }),
		func(c *core.ComponentSpec) string { return c.TagName },
	)

	context := &core.MJMLContext{
		Components:     components,
		RawComponents:  rawComponents,
		ContainerWidth: 0,
		PrinterSupport: false,
		GlobalData: &core.GlobalData{
			BeforeDoctype:     "",
			Breakpoint:        "480px",
			Classes:           make(map[string]*orderedmap.OrderedMap[string, string]),
			DefaultAttributes: defaultAttributes,
			InlineStyles:      defaultStyles,
			Fonts:             fonts,
			MediaQueries:      orderedmap.NewOrderedMap[string, string](),
			HeadStyles:        make(map[string]func(breakpoint string) string),
			ForceOWADesktop:   owa == "desktop",
			Lang:              lang,
			Dir:               dir,
			Data:              data,
		},
		Processing: func(w core.MJMLWriter, xml string, ctx *core.MJMLContext) error {
			// supports returning siblings elements from a custom component
			wrapped := "<fragment>" + xml + "</fragment>"
			partialMjml, err := parser.MJMLParser(unsafe.Slice(unsafe.StringData(wrapped), len(wrapped)), parser.MJMLParserOptions{
				Components:     endingComponents,
				KeepComments:   keepComments,
				FilePath:       ".",
				ActualPath:     ".",
				IgnoreIncludes: true,
			}, nil)

			if err != nil {
				return err
			}

			for _, child := range partialMjml.Children {
				attributes := ctx.GetAttributes(child)
				err = processing(w, child, &core.ComponentProps{Attributes: utils.ToEntries(attributes)}, ctx)

				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	sb := new(core.MJMLBuilder)
	err = processing(sb, mjHead, &core.ComponentProps{}, context)

	if err != nil {
		return "", err
	}

	headRaw := sb.String()

	sb = new(core.MJMLBuilder)

	attributes := context.GetAttributes(mjBody)
	err = processing(sb, mjBody, &core.ComponentProps{Attributes: utils.ToEntries(attributes)}, context)

	if err != nil {
		return "", err
	}

	if sb.Len() <= 0 {
		return "", errors.New("Malformed MJML. Check that your structure is correct and enclosed in <mjml> tags.")
	}

	content := helpers.MinifyOutlookConditionals(sb.String())

	content, err = core.DefaultSkeleton(headRaw, content, context)
	if err != nil {
		return "", err
	}

	var dom []*dom.Node

	if len(context.GlobalData.HtmlAttributes) > 0 {
		dom = parseHtml(content, dom)

		for _, rule := range context.GlobalData.HtmlAttributes {
			matches, err := query.SelectAll(rule.Path, dom, nil)
			if err != nil {
				return "", nil
			}

			for _, node := range matches {
				for _, attr := range rule.Attributes {
					node.Attribs.Set(attr.Name, attr.Value)
				}
			}
		}
	}

	if len(context.GlobalData.InlineStyles) > 0 {
		dom = parseHtml(content, dom)
		parser := cssparser.NewCssParser()
		for _, css := range context.GlobalData.InlineStyles {
			styles, err := parser.Parse(css)
			if err != nil {
				return "", err
			}

			err = inlineStyles(dom, styles)
			if err != nil {
				return "", err
			}
		}
	}

	if len(dom) > 0 {
		content = serializer.Render(dom, &serializer.DomSerializerOptions{XmlMode: ptr("false"), SelfClosingTags: ptr(false), DecodeEntities: ptr(false)})
	}

	content = helpers.MergeOutlookConditionals(content)

	return content, nil
}

func ptr[T any](v T) *T {
	return &v
}

func processing(w core.MJMLWriter, node *parser.MJMLNode, props *core.ComponentProps, context *core.MJMLContext) error {
	if node == nil {
		return nil
	}

	component, err := core.InitComponent(node, props, context)

	if err != nil {
		return err
	}

	if component != nil {
		return component.Render(w)
	}

	return nil
}

func parseHtml(content string, nodes []*dom.Node) []*dom.Node {
	if nodes != nil {
		return nodes
	}

	return dom.ParseDOM(content, &htmlparser2.ParserOptions{
		XmlMode:                 false,
		LowerCaseAttributeNames: true,
		DecodeEntities:          false,
		RecognizeSelfClosing:    true,
	})
}

func inlineStyles(dom []*dom.Node, styles []*cssparser.Rule) error {
	for _, rule := range styles {
		matches, err := query.SelectAll(rule.Selector, dom, nil)
		if err != nil {
			return err
		}

		const important = "!important"
		var prefix = make([]core.Style, 0, len(rule.Declarations))
		var suffix = make([]core.Style, 0, len(rule.Declarations))
		var suffixNames = make([]string, 0, len(rule.Declarations))
		for _, d := range rule.Declarations {
			lenValue := len(d.Value)
			if lenValue > len(important) && strings.EqualFold(d.Value[lenValue-len(important):], important) {
				suffix = append(suffix, core.Style{Name: d.Property, Value: d.Value[:lenValue-len(important)]})
				suffixNames = append(suffixNames, d.Property)
			} else {
				prefix = append(prefix, core.Style{Name: d.Property, Value: d.Value})
			}
		}

		for _, node := range matches {
			decls := strings.Split(strings.TrimSpace(node.Attribs.GetOrDefault("style", "")), ";")
			oldStyle := make([]core.Style, 0, len(decls))
			oldStyleNames := make([]string, 0, len(decls))
			for _, d := range decls {
				if d == "" {
					continue
				}

				idx := strings.Index(d, ":")
				name := strings.TrimSpace(d[:idx])
				value := strings.TrimSpace(d[idx+1:])
				oldStyle = append(oldStyle, core.Style{Name: name, Value: value})
				oldStyleNames = append(oldStyleNames, name)
			}
			style := make([]string, 0, len(prefix)+len(suffix)+len(oldStyle))

			for _, s := range prefix {
				if !slices.Contains(oldStyleNames, s.Name) {
					style = append(style, s.Name+":"+s.Value)
				}
			}

			for _, s := range oldStyle {
				if !slices.Contains(suffixNames, s.Name) {
					style = append(style, s.Name+":"+s.Value)
				}
			}

			for _, s := range suffix {
				style = append(style, s.Name+":"+s.Value)
			}

			node.Attribs.Set("style", strings.Join(style, ";"))

		}
	}

	return nil
}
