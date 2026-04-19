package helpers

import (
	"fmt"
	"io"
	"regexp"
	"slices"
	"strings"

	"github.com/krozhkov/mgml/internal/utils"
)

var fontRegex = regexp.MustCompile(`(?im)[^"]*font-family:([^";]+)`)
var fontInlineRegex = regexp.MustCompile(`(?im)font-family:([^;}]+)`)

type FontDeclaration struct {
	Name string
	Href string
}

func BuildFontsTags(w io.StringWriter, content string, inlineStyles []string, fontDeclarations []*FontDeclaration) error {
	usages := fontRegex.FindAllStringSubmatch(content, -1)
	for _, style := range inlineStyles {
		inlineUsages := fontInlineRegex.FindAllStringSubmatch(style, -1)
		if inlineUsages != nil {
			usages = append(usages, inlineUsages...)
		}
	}

	var fontNames []string
	for _, usage := range usages {
		names := strings.Split(usage[1], ",")
		for _, name := range names {
			name = strings.TrimSpace(name)
			if name[0] == '\'' && name[len(name)-1] == '\'' {
				name = name[1 : len(name)-1]
			}
			fontNames = append(fontNames, strings.ToLower(name))
		}
	}

	fonts := utils.NewSet(fontNames...)

	if slices.IndexFunc(fontDeclarations, func(decl *FontDeclaration) bool { return fonts.Has(strings.ToLower(decl.Name)) }) == -1 {
		return nil
	}

	if _, err := w.WriteString("<!--[if !mso]><!-->"); err != nil {
		return err
	}

	for _, decl := range fontDeclarations {
		if fonts.Has(strings.ToLower(decl.Name)) {
			w.WriteString(fmt.Sprintf("\n<link href=\"%s\" rel=\"stylesheet\" type=\"text/css\">", decl.Href))
		}
	}

	if _, err := w.WriteString("\n<style type=\"text/css\">"); err != nil {
		return err
	}

	for _, decl := range fontDeclarations {
		if fonts.Has(strings.ToLower(decl.Name)) {
			w.WriteString(fmt.Sprintf("\n@import url(%s);", decl.Href))
		}
	}

	if _, err := w.WriteString("\n</style>"); err != nil {
		return err
	}

	if _, err := w.WriteString("\n<!--<![endif]-->\n"); err != nil {
		return err
	}

	return nil
}
