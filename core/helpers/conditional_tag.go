package helpers

import "io"

var startConditionalTag = "<!--[if mso | IE]>"
var startMsoConditionalTag = "<!--[if mso]>"
var endConditionalTag = "<![endif]-->"
var startNegationConditionalTag = "<!--[if !mso | IE]><!-->"
var startMsoNegationConditionalTag = "<!--[if !mso]><!-->"
var endNegationConditionalTag = "<!--<![endif]-->"

func ConditionalTagStart(w io.StringWriter, negation bool) (int, error) {
	if negation {
		return w.WriteString(startNegationConditionalTag)
	} else {
		return w.WriteString(startConditionalTag)
	}
}

func ConditionalTagEnd(w io.StringWriter, negation bool) (int, error) {
	if negation {
		return w.WriteString(endNegationConditionalTag)
	} else {
		return w.WriteString(endConditionalTag)
	}
}

func MsoConditionalTagStart(w io.StringWriter, negation bool) (int, error) {
	if negation {
		return w.WriteString(startMsoNegationConditionalTag)
	} else {
		return w.WriteString(startMsoConditionalTag)
	}
}
