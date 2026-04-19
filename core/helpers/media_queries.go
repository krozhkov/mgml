package helpers

import (
	"fmt"
	"io"

	"github.com/krozhkov/mgml/parser"
)

type MediaQueriesOptions struct {
	ForceOWADesktop bool
	PrinterSupport  bool
}

func BuildMediaQueriesTags(w io.StringWriter, breakpoint string, mediaQueries *parser.Attributes, options MediaQueriesOptions) error {
	if mediaQueries == nil || mediaQueries.Len() == 0 {
		return nil
	}

	forceOWADesktop := options.ForceOWADesktop
	printerSupport := options.PrinterSupport

	var baseMediaQueries = make([]string, 0, mediaQueries.Len())
	for className, mediaQuery := range mediaQueries.AllFromFront() {
		baseMediaQueries = append(baseMediaQueries, fmt.Sprintf(".%s %s", className, mediaQuery))
	}

	var thunderbirdMediaQueries = make([]string, 0, mediaQueries.Len())
	for className, mediaQuery := range mediaQueries.AllFromFront() {
		thunderbirdMediaQueries = append(thunderbirdMediaQueries, fmt.Sprintf(".moz-text-html .%s %s", className, mediaQuery))
	}

	var owaQueries = make([]string, 0, len(baseMediaQueries))
	for _, mq := range baseMediaQueries {
		owaQueries = append(owaQueries, fmt.Sprintf("[owa] %s", mq))
	}

	if _, err := w.WriteString("<style type=\"text/css\">\n"); err != nil {
		return err
	}

	if _, err := w.WriteString(fmt.Sprintf("@media only screen and (min-width:%s) {", breakpoint)); err != nil {
		return err
	}

	for _, mq := range baseMediaQueries {
		if _, err := w.WriteString(mq); err != nil {
			return err
		}
		if _, err := w.WriteString("\n"); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("}\n</style>\n"); err != nil {
		return err
	}

	if _, err := w.WriteString(fmt.Sprintf("<style media=\"screen and (min-width:%s)\">", breakpoint)); err != nil {
		return err
	}

	for _, mq := range thunderbirdMediaQueries {
		if _, err := w.WriteString(mq); err != nil {
			return err
		}
		if _, err := w.WriteString("\n"); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("</style>\n"); err != nil {
		return err
	}

	if printerSupport {
		if _, err := w.WriteString("<style type=\"text/css\">\n@media only print {\n"); err != nil {
			return err
		}

		for _, mq := range baseMediaQueries {
			if _, err := w.WriteString(mq); err != nil {
				return err
			}
			if _, err := w.WriteString("\n"); err != nil {
				return err
			}
		}

		if _, err := w.WriteString("}\n</style>\n"); err != nil {
			return err
		}
	}

	if forceOWADesktop {
		if _, err := w.WriteString("<style type=\"text/css\">\n"); err != nil {
			return err
		}

		for _, mq := range owaQueries {
			if _, err := w.WriteString(mq); err != nil {
				return err
			}
			if _, err := w.WriteString("\n"); err != nil {
				return err
			}
		}

		if _, err := w.WriteString("</style>\n"); err != nil {
			return err
		}
	}

	return nil
}
