package helpers

import "regexp"

var outlookRegex = regexp.MustCompile(`(?im)<!\[endif\]-->\s*?<!--\[if mso \| IE\]>`)

// # OPTIMIZE ME: — check if previous conditionnal is `<!--[if mso | I`]>` too
func MergeOutlookConditionals(content string) string {
	return outlookRegex.ReplaceAllString(content, "")
}
