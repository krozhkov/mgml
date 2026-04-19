package helpers

import "regexp"

var outlookConditionalRegex = regexp.MustCompile(`(?m)(<!--\[if\s[^\]]+]>)([\s\S]*?)(<!\[endif]-->)`)
var spaceBetweenTagsRegex = regexp.MustCompile(`(?m)(^|>)(\s+)(<|$)`)
var sequentialSpacesRegex = regexp.MustCompile(`(?m)\s{2,}`)

func MinifyOutlookConditionals(content string) string {
	return outlookConditionalRegex.ReplaceAllStringFunc(content, func(match string) string {
		return sequentialSpacesRegex.ReplaceAllLiteralString(
			spaceBetweenTagsRegex.ReplaceAllLiteralString(match, "><"), " ",
		)
	})
}
