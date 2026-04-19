package core

import (
	"fmt"
	"strings"

	"github.com/krozhkov/mgml/core/helpers"
)

func DefaultSkeleton(headRaw string, content string, ctx *MJMLContext) (string, error) {
	w := new(strings.Builder)

	if ctx.GlobalData.BeforeDoctype != "" {
		if _, err := w.WriteString(ctx.GlobalData.BeforeDoctype); err != nil {
			return "", err
		}
	}

	if _, err := w.WriteString("<!doctype html>\n"); err != nil {
		return "", err
	}

	if _, err := w.WriteString(fmt.Sprintf("<html lang=\"%s\" dir=\"%s\" xmlns=\"http://www.w3.org/1999/xhtml\" xmlns:v=\"urn:schemas-microsoft-com:vml\" xmlns:o=\"urn:schemas-microsoft-com:office:office\">\n", ctx.GlobalData.Lang, ctx.GlobalData.Dir)); err != nil {
		return "", err
	}

	if _, err := w.WriteString("<head>\n"); err != nil {
		return "", err
	}

	var title string
	if ctx.GlobalData.Title != nil {
		title = *ctx.GlobalData.Title
	}
	if _, err := w.WriteString(fmt.Sprintf("<title>%s</title>\n", title)); err != nil {
		return "", err
	}

	if _, err := w.WriteString(`<!--[if !mso]><!-->
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <!--<![endif]-->
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style type="text/css">
      #outlook a { padding:0; }
      body { margin:0;padding:0;-webkit-text-size-adjust:100%;-ms-text-size-adjust:100%; }
      table, td { border-collapse:collapse;mso-table-lspace:0pt;mso-table-rspace:0pt; }
      img { border:0;height:auto;line-height:100%; outline:none;text-decoration:none;-ms-interpolation-mode:bicubic; }
      p { display:block;margin:13px 0; }
    </style>
    <!--[if mso]>
    <noscript>
    <xml>
    <o:OfficeDocumentSettings>
      <o:AllowPNG/>
      <o:PixelsPerInch>96</o:PixelsPerInch>
    </o:OfficeDocumentSettings>
    </xml>
    </noscript>
    <![endif]-->
    <!--[if lte mso 11]>
    <style type="text/css">
      .mj-outlook-group-fix { width:100% !important; }
    </style>
    <![endif]-->
	`); err != nil {
		return "", err
	}

	if err := helpers.BuildFontsTags(w, content, ctx.GlobalData.InlineStyles, ctx.GlobalData.Fonts); err != nil {
		return "", err
	}

	if err := helpers.BuildMediaQueriesTags(w, ctx.GlobalData.Breakpoint, ctx.GlobalData.MediaQueries, helpers.MediaQueriesOptions{ForceOWADesktop: ctx.GlobalData.ForceOWADesktop, PrinterSupport: ctx.PrinterSupport}); err != nil {
		return "", err
	}

	if err := helpers.BuildStyleFromComponents(w, ctx.GlobalData.Breakpoint, ctx.GlobalData.ComponentsHeadStyles, ctx.GlobalData.HeadStyles); err != nil {
		return "", err
	}

	if err := helpers.BuildStyleFromTags(w, ctx.GlobalData.Breakpoint, ctx.GlobalData.Styles); err != nil {
		return "", err
	}

	if _, err := w.WriteString(headRaw); err != nil {
		return "", err
	}

	if _, err := w.WriteString("</head>\n"); err != nil {
		return "", err
	}

	if _, err := w.WriteString("<body style=\"word-spacing:normal;"); err != nil {
		return "", err
	}

	if ctx.GlobalData.BackgroundColor != nil {
		if _, err := w.WriteString(fmt.Sprintf("background-color:%s;", *ctx.GlobalData.BackgroundColor)); err != nil {
			return "", err
		}
	}

	if _, err := w.WriteString("\">\n"); err != nil {
		return "", err
	}

	if ctx.GlobalData.Preview != nil {
		if err := helpers.BuildPreview(w, *ctx.GlobalData.Preview); err != nil {
			return "", err
		}
	}

	if _, err := w.WriteString(content); err != nil {
		return "", err
	}

	if _, err := w.WriteString("</body>\n</html>"); err != nil {
		return "", err
	}

	return w.String(), nil
}
