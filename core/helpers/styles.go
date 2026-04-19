package helpers

import "io"

func BuildStyleFromComponents(
	w io.StringWriter,
	breakpoint string,
	componentsHeadStyles []func(breakpoint string) string,
	headStyles map[string]func(breakpoint string) string,
) error {
	if len(headStyles) == 0 && len(componentsHeadStyles) == 0 {
		return nil
	}

	if _, err := w.WriteString("<style type=\"text/css\">\n"); err != nil {
		return err
	}

	for _, styleFunction := range componentsHeadStyles {
		if _, err := w.WriteString(styleFunction(breakpoint)); err != nil {
			return err
		}
		if _, err := w.WriteString("\n"); err != nil {
			return err
		}
	}

	for _, styleFunction := range headStyles {
		if _, err := w.WriteString(styleFunction(breakpoint)); err != nil {
			return err
		}
		if _, err := w.WriteString("\n"); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("</style>\n"); err != nil {
		return err
	}

	return nil
}

func BuildStyleFromTags(w io.StringWriter, breakpoint string, styles []string) error {
	if len(styles) == 0 {
		return nil
	}

	if _, err := w.WriteString("<style type=\"text/css\">\n"); err != nil {
		return err
	}

	for _, style := range styles {
		if _, err := w.WriteString(style); err != nil {
			return err
		}
		if _, err := w.WriteString("\n"); err != nil {
			return err
		}
	}

	if _, err := w.WriteString("</style>\n"); err != nil {
		return err
	}

	return nil
}
