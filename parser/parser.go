package parser

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/elliotchance/orderedmap/v3"
	"github.com/krozhkov/go-htmlparser2/parser"
	"github.com/krozhkov/mgml/internal/utils"
)

type Attribute = orderedmap.Element[string, string]

type Attributes = orderedmap.OrderedMap[string, string]

func indexesForNewLine(data []byte) []int {
	indexes := make([]int, 1, 10)

	for i, b := range data {
		if b == '\n' {
			indexes = append(indexes, i+1)
		}
	}

	return indexes
}

func findTag(tagName string, tree *MJMLNode) *MJMLNode {
	for _, child := range tree.Children {
		if child.TagName == tagName {
			return child
		}
	}

	return nil
}

func bindToTree(children []*MJMLNode, tree *MJMLNode) {
	for _, child := range children {
		child.Parent = tree
	}
}

func cleanNode(node *MJMLNode) {
	node.Parent = nil

	// Delete children if needed
	if len(node.Children) > 0 {
		for _, child := range node.Children {
			cleanNode(child)
		}
	} else {
		node.Children = nil
	}

	// Delete attributes if needed
	if node.Attributes != nil && node.Attributes.Len() == 0 {
		node.Attributes = nil
	}
}

func setEmptyAttributes(node *MJMLNode) {
	if node.Attributes == nil {
		node.Attributes = orderedmap.NewOrderedMap[string, string]()
	}

	for _, child := range node.Children {
		setEmptyAttributes(child)
	}
}

type MJMLNode struct {
	File             string
	Line             int
	AbsoluteFilePath string
	TagName          string
	IncludedIn       []*MJMLIncludedIn
	Attributes       *Attributes
	Content          string

	Parent   *MJMLNode
	Children []*MJMLNode

	Err error
}

type MJMLIncludedIn struct {
	File string
	Line int
}

type MJMLComponent struct {
	TagName   string
	EndingTag bool
}

type MJMLParserOptions struct {
	AddEmptyAttributes bool
	Components         []MJMLComponent
	ConvertBooleans    bool
	KeepComments       bool
	FilePath           string
	ActualPath         string
	IgnoreIncludes     bool
	Preprocessors      []MJMLPreprocessor
}

type MJMLPreprocessor func([]byte) []byte

type mjmlParser struct {
	xml                []byte
	components         []MJMLComponent
	preprocessors      []MJMLPreprocessor
	addEmptyAttributes bool
	keepComments       bool
	convertBooleans    bool
	ignoreIncludes     bool
	filePath           string
	actualPath         string
	cwd                string

	includedIn []*MJMLIncludedIn
	inInclude  bool

	endingTags                 []string
	inEndingTag                int
	currentEndingTagStartIndex int
	currentEndingTagEndIndex   int
	lineIndexes                []int
	mjml                       *MJMLNode
	cur                        *MJMLNode
	cssIncludes                []*MJMLNode

	err    error
	parser *parser.Parser
}

func MJMLParser(xml []byte, options MJMLParserOptions, includedIn []*MJMLIncludedIn) (*MJMLNode, error) {
	p := &mjmlParser{
		xml:                xml,
		components:         options.Components,
		preprocessors:      options.Preprocessors,
		addEmptyAttributes: options.AddEmptyAttributes,
		keepComments:       options.KeepComments,
		convertBooleans:    options.ConvertBooleans,
		ignoreIncludes:     options.IgnoreIncludes,
		includedIn:         includedIn,
		inInclude:          len(includedIn) > 0,
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if options.FilePath == "" {
		p.filePath = "."
	} else {
		p.filePath = options.FilePath
		cwd, err = basePath(options.FilePath)
		if err != nil {
			return nil, err
		}
	}

	p.cwd = cwd

	if options.ActualPath == "" {
		p.actualPath = "."
	} else {
		p.actualPath = options.ActualPath
	}

	p.lineIndexes = indexesForNewLine(xml)

	p.endingTags = make([]string, 0, 10)
	for _, component := range options.Components {
		if component.EndingTag {
			p.endingTags = append(p.endingTags, component.TagName)
		}
	}

	p.parser = parser.NewParser(
		p,
		&parser.ParserOptions{
			XmlMode:                 false,
			LowerCaseTags:           true,
			RecognizeCDATA:          true,
			DecodeEntities:          false,
			RecognizeSelfClosing:    true,
			LowerCaseAttributeNames: false,
		},
	)

	mjml, err := p.Parse()
	if err != nil {
		return nil, err
	}

	if p.err != nil {
		return nil, p.err
	}

	return mjml, nil
}

func (p *mjmlParser) handleCssHtmlInclude(file string, attrs *Attributes, line int) error {
	partialPath, err := resolvePath(p.cwd, file)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(partialPath)
	if err != nil {
		absolutePath, err := resolvePath(p.cwd, p.actualPath)
		if err != nil {
			return err
		}

		newNode := &MJMLNode{
			File:             file,
			AbsoluteFilePath: absolutePath,
			Line:             line,
			Parent:           p.cur,
			TagName:          "mj-raw",
			Content:          fmt.Sprintf("<!-- mj-include fails to read file : %s at %s -->", file, partialPath),
			Err:              fmt.Errorf("mj-include fails to read file %s", partialPath),
		}

		p.cur.Children = append(p.cur.Children, newNode)

		return nil
	}

	if attrs.GetOrDefault("type", "") == "html" {
		absolutePath, err := resolvePath(p.cwd, p.actualPath)
		if err != nil {
			return err
		}

		newNode := &MJMLNode{
			File:             file,
			AbsoluteFilePath: absolutePath,
			Line:             line,
			Parent:           p.cur,
			TagName:          "mj-raw",
			Content:          string(content),
		}

		p.cur.Children = append(p.cur.Children, newNode)

		return nil
	}

	var attributes = orderedmap.NewOrderedMap[string, string]()
	if attrs.GetOrDefault("css-inline", "") == "inline" {
		attributes.Set("inline", "inline")
	}

	absolutePath, err := resolvePath(p.cwd, p.actualPath)
	if err != nil {
		return err
	}

	newNode := &MJMLNode{
		File:             file,
		AbsoluteFilePath: absolutePath,
		Line:             line,
		Parent:           p.cur,
		TagName:          "mj-style",
		Content:          string(content),
		Attributes:       attributes,
	}

	p.cssIncludes = append(p.cssIncludes, newNode)

	return nil
}

func (p *mjmlParser) handleInclude(file string, line int) error {
	partialPath, err := resolvePath(p.cwd, file)
	if err != nil {
		return err
	}

	if utils.LastIndexFunc(p.cur.IncludedIn, func(e *MJMLIncludedIn) bool { return e.File == partialPath }) != -1 {
		return fmt.Errorf("circular inclusion detected on file : %s", partialPath)
	}

	content, err := os.ReadFile(partialPath)

	if err != nil {
		absolutePath, err := resolvePath(p.cwd, p.actualPath)
		if err != nil {
			return err
		}

		newNode := &MJMLNode{
			File:             file,
			AbsoluteFilePath: absolutePath,
			Line:             line,
			Parent:           p.cur,
			TagName:          "mj-raw",
			Content:          fmt.Sprintf("<!-- mj-include fails to read file : %s at %s -->", file, partialPath),
			Err:              fmt.Errorf("mj-include fails to read file %s", partialPath),
		}

		p.cur.Children = append(p.cur.Children, newNode)

		return nil
	}

	if !bytes.Contains(content, []byte("<mjml>")) {
		content = bytes.Join([][]byte{[]byte("<mjml><mj-body>"), content, []byte("</mj-body></mjml>")}, nil)
	}

	partialMjml, err := MJMLParser(
		content,
		MJMLParserOptions{
			AddEmptyAttributes: p.addEmptyAttributes,
			Components:         p.components,
			ConvertBooleans:    p.convertBooleans,
			KeepComments:       p.keepComments,
			IgnoreIncludes:     p.ignoreIncludes,
			Preprocessors:      p.preprocessors,
			FilePath:           partialPath,
			ActualPath:         partialPath,
		},
		append(p.cur.IncludedIn, &MJMLIncludedIn{File: p.cur.AbsoluteFilePath, Line: line}),
	)
	if err != nil {
		return err
	}

	if partialMjml.TagName != "mjml" {
		return nil
	}

	body := findTag("mj-body", partialMjml)
	head := findTag("mj-head", partialMjml)

	if body != nil {
		bindToTree(body.Children, p.cur)
		p.cur.Children = slices.Concat(p.cur.Children, body.Children)
	}

	if head != nil {
		curHead := findTag("mj-head", p.mjml)

		if curHead == nil {
			absolutePath, err := resolvePath(p.cwd, p.actualPath)
			if err != nil {
				return err
			}

			curHead = &MJMLNode{
				File:             p.actualPath,
				AbsoluteFilePath: absolutePath,
				Parent:           p.cur,
				TagName:          "mj-head",
			}

			p.mjml.Children = append(p.mjml.Children, curHead)
		}

		bindToTree(head.Children, curHead)
		curHead.Children = slices.Concat(curHead.Children, head.Children)
	}

	return nil
}

func (p *mjmlParser) Parse() (*MJMLNode, error) {
	xml := p.xml
	// Apply preprocessors to raw xml
	for _, preprocessor := range p.preprocessors {
		xml = preprocessor(xml)
	}
	p.xml = xml

	p.parser.End(xml)

	if p.mjml == nil {
		return nil, errors.New("parsing failed. check your mjml")
	}

	cleanNode(p.mjml)

	// Assign "attributes" property if not set
	if p.addEmptyAttributes {
		setEmptyAttributes(p.mjml)
	}

	if len(p.cssIncludes) > 0 {
		head := findTag("mj-head", p.mjml)

		if head != nil {
			head.Children = slices.Concat(head.Children, p.cssIncludes)
		} else {
			newNode := &MJMLNode{
				File:     p.filePath,
				Line:     0,
				TagName:  "mj-head",
				Children: p.cssIncludes,
			}

			p.mjml.Children = append(p.mjml.Children, newNode)
		}
	}

	return p.mjml, nil
}

func (p *mjmlParser) OnOpenTag(name string, attrs []*parser.Attribute, isImplied bool) {
	isAnEndingTag := slices.Index(p.endingTags, name) != -1

	if p.inEndingTag > 0 {
		if isAnEndingTag {
			p.inEndingTag += 1
		}
		return
	}

	if isAnEndingTag {
		p.inEndingTag += 1

		if p.inEndingTag == 1 {
			// we're entering endingTag
			p.currentEndingTagStartIndex = p.parser.StartIndex
			p.currentEndingTagEndIndex = p.parser.EndIndex
		}
	}

	line := utils.LastIndexFunc(p.lineIndexes, func(i int) bool { return i <= p.parser.StartIndex }) + 1

	attributes := orderedmap.NewOrderedMapWithCapacity[string, string](len(attrs))
	for _, el := range attrs {
		attributes.Set(el.Name, el.Value)
	}

	if name == "mj-include" {
		if p.ignoreIncludes {
			return
		}

		if attributes.GetOrDefault("type", "") == "css" || attributes.GetOrDefault("type", "") == "html" {
			unescapedPath, err := url.PathUnescape(attributes.GetOrDefault("path", ""))
			if err != nil {
				p.err = err
				return
			}

			err = p.handleCssHtmlInclude(unescapedPath, attributes, line)
			if err != nil {
				p.err = err
			}

			return
		}

		p.inInclude = true

		unescapedPath, err := url.PathUnescape(attributes.GetOrDefault("path", ""))
		if err != nil {
			p.err = err
			return
		}

		err = p.handleInclude(unescapedPath, line)
		if err != nil {
			p.err = err
		}

		return
	}

	if p.convertBooleans {
		// "true" and "false" will be converted to bools
		// attrs = convertBooleansOnAttrs(attrs)
	}

	absolutePath, err := resolvePath(p.cwd, p.actualPath)
	if err != nil {
		absolutePath = p.actualPath
	}

	newNode := &MJMLNode{
		File:             p.actualPath,
		AbsoluteFilePath: absolutePath,
		Line:             line,
		IncludedIn:       p.includedIn,
		Parent:           p.cur,
		TagName:          name,
		Attributes:       attributes,
	}

	if p.cur != nil {
		p.cur.Children = append(p.cur.Children, newNode)
	} else {
		p.mjml = newNode
	}

	p.cur = newNode
}
func (p *mjmlParser) OnCloseTag(name string, isImplied bool) {
	if slices.Index(p.endingTags, name) != -1 {
		p.inEndingTag -= 1

		if p.inEndingTag == 0 {
			// we're getting out of endingTag
			// if self-closing tag we don't get the content
			if !p.isSelfClosing() {
				partialVal := bytes.TrimSpace(p.xml[p.currentEndingTagEndIndex+1 : p.parser.EndIndex])

				val := partialVal[0:bytes.LastIndex(partialVal, []byte("</"+name))]

				if len(val) > 0 {
					p.cur.Content = string(bytes.TrimSpace(val))
				}
			}
		}
	}

	if p.inEndingTag > 0 {
		return
	}

	if p.inInclude {
		p.inInclude = false
	}

	// for includes, setting cur is handled in handleInclude because when there is
	// only mj-head in include it doesn't create any elements, so setting back to parent is wrong
	if name != "mj-include" {
		if p.cur != nil && p.cur.Parent != nil {
			p.cur = p.cur.Parent
		} else {
			p.cur = nil
		}
	}
}
func (p *mjmlParser) OnText(text string) {
	if p.inEndingTag > 0 {
		return
	}

	text = strings.TrimSpace(text)

	if text != "" && p.cur != nil {
		p.cur.Content += text
	}
}
func (p *mjmlParser) OnComment(data string) {
	if p.inEndingTag > 0 {
		return
	}

	if p.cur != nil && p.keepComments {
		line := utils.LastIndexFunc(p.lineIndexes, func(i int) bool { return i <= p.parser.StartIndex }) + 1

		newNode := &MJMLNode{
			File:       p.actualPath,
			Line:       line,
			IncludedIn: p.includedIn,
			Parent:     p.cur,
			TagName:    "mj-raw",
			Content:    fmt.Sprintf("<!--%s-->", data),
		}

		p.cur.Children = append(p.cur.Children, newNode)
	}
}
func (p *mjmlParser) OnParserInit(parser *parser.Parser) {
}
func (p *mjmlParser) OnReset() {
}
func (p *mjmlParser) OnEnd() {
}
func (p *mjmlParser) OnError(e error) {
}
func (p *mjmlParser) OnOpenTagName(name string) {
}
func (p *mjmlParser) OnAttribute(name string, value string, quote parser.QuoteType) {
}
func (p *mjmlParser) OnCDataStart() {
}
func (p *mjmlParser) OnCDataEnd() {
}
func (p *mjmlParser) OnCommentEnd() {
}
func (p *mjmlParser) OnProcessingInstruction(name string, data string) {
}

func (p *mjmlParser) isSelfClosing() bool {
	return p.currentEndingTagStartIndex == p.parser.StartIndex && p.currentEndingTagEndIndex == p.parser.EndIndex
}
