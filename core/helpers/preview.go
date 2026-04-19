package helpers

import "io"

func BuildPreview(w io.StringWriter, content string) error {
	if content != "" {
		if _, err := w.WriteString("<div style=\"display:none;font-size:1px;color:#ffffff;line-height:1px;max-height:0px;max-width:0px;opacity:0;overflow:hidden;\">"); err != nil {
			return err
		}

		if _, err := w.WriteString(content); err != nil {
			return err
		}

		if _, err := w.WriteString("</div>\n"); err != nil {
			return err
		}

	}

	return nil
}
