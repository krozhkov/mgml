package mjml_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/krozhkov/go-css-select/query"
	"github.com/krozhkov/go-htmlparser2/dom"
	"github.com/krozhkov/go-htmlparser2/parser"
	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/internal/utils"
	"github.com/krozhkov/mgml/mjml"
	"github.com/stretchr/testify/assert"
)

func parseDOM(str string) []*dom.Node {
	return dom.ParseDOM(str, &parser.ParserOptions{LowerCaseAttributeNames: true, DecodeEntities: true, RecognizeSelfClosing: true})
}

func TestRendering(t *testing.T) {
	t.Run("should set context", func(t *testing.T) {
		input := `<mjml>
			<mj-head>
				<mj-attributes>
					<mj-all font-family="Arial" />
					<mj-text color="#555" />
					<mj-class name="big-red" font-size="20px" color="red" />
				</mj-attributes>
				<mj-html-attributes>
					<mj-selector path=".custom > div">
						<mj-html-attribute name="data-id">42</mj-html-attribute>
					</mj-selector>
				</mj-html-attributes>
				<mj-style inline="inline">
					body > div > div {
						color: blue !important;
					}
				</mj-style>
			</mj-head>
			<mj-body background-color="green">
				<mj-group>
				<mj-column vertical-align="middle" width="65%">
					<mj-text css-class="custom" padding="0 20px 0 0" align="left">LOGO</mj-text>
				</mj-column>
				<mj-column vertical-align="middle" width="35%">
					<mj-text color="white" padding="0 20px 0 0" font-weight="500" font-size="20px" align="right">TESTING</mj-text>
				</mj-column>
				</mj-group>
				<mj-text mj-class="big-red">Hello World</mj-text>
			</mj-body>
		</mjml>`

		html, err := mjml.MJML2Html(input, &core.MJMLOptions{
			KeepComments:    true,
			IgnoreIncludes:  true,
			ValidationLevel: core.ValidationLevelSoft,
		})

		assert.Nil(t, err)

		snaps.MatchSnapshot(t, html)
	})

	t.Run("should match mjml output", func(t *testing.T) {
		xml := `
<mjml>
    <mj-head>
        <mj-attributes>
            <mj-all
                padding="0px"
            />
            <mj-wrapper
                background-color="yellow"
                padding="80px"
            />
        </mj-attributes>
    </mj-head>
    <mj-body>
        <mj-wrapper>
            <mj-section>
                <mj-column>
                    <mj-text>
                        lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem lorem
                    </mj-text>
                </mj-column>
            </mj-section>
        </mj-wrapper>
    </mj-body>
</mjml>
`

		html, err := mjml.MJML2Html(xml, &core.MJMLOptions{
			KeepComments:    true,
			IgnoreIncludes:  true,
			ValidationLevel: core.ValidationLevelSoft,
		})

		assert.Nil(t, err)

		snaps.MatchSnapshot(t, html)
	})

	t.Run("html-attributes should put the attributes at the right place", func(t *testing.T) {
		input := `
<mjml>
  <mj-head>
    <mj-html-attributes>
      <mj-selector path=".text div">
        <mj-html-attribute name="data-id">42</mj-html-attribute>
      </mj-selector>
      <mj-selector path=".image td">
        <mj-html-attribute name="data-name">43</mj-html-attribute>
      </mj-selector>
    </mj-html-attributes>
  </mj-head>
  <mj-body>
    <mj-raw>{ if item < 5 }</mj-raw>
    <mj-section css-class="section">
      <mj-column>
	  	<mj-image css-class="image2" src="https://placehold.co/150x40"/>
        <mj-raw>{ if item > 10 }</mj-raw>
        <mj-text css-class="text">
          Hello World! { item }
        </mj-text>
        <mj-raw>{ end if }</mj-raw>
        <mj-text css-class="text">
          Hello World! { item + 1 }
        </mj-text>
        <mj-image css-class="image" src="https://placehold.co/150x30"/>
      </mj-column>
    </mj-section>
    <mj-raw>{ end if }</mj-raw>
  </mj-body>
</mjml>
`

		html, err := mjml.MJML2Html(input, &core.MJMLOptions{
			KeepComments:    true,
			IgnoreIncludes:  true,
			ValidationLevel: core.ValidationLevelSoft,
		})

		assert.Nil(t, err)
		snaps.MatchSnapshot(t, html)

		nodes := parseDOM(html)

		matches, err := query.SelectAll(".text div", nodes, nil)
		assert.Nil(t, err)
		assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string { return n.Attribs.GetOrDefault("data-id", "") }), []string{"42", "42"})

		matches, err = query.SelectAll(".image td", nodes, nil)
		assert.Nil(t, err)
		assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string { return n.Attribs.GetOrDefault("data-name", "") }), []string{"43"})

		// should not alter templating syntax, or move the content that is outside any tag (mj-raws)
		expected := []string{
			"{ if item < 5 }",
			"class=\"section\"",
			"{ if item > 10 }",
			"class=\"text\"",
			"{ item }",
			"{ end if }",
			"{ item + 1 }",
		}

		indexes := utils.MapFunc(expected, func(s string) int { return strings.Index(html, s) })

		assert.NotContains(t, indexes, -1)

		sorted := slices.Clone(indexes)
		slices.Sort(sorted)

		assert.Equal(t, sorted, indexes)
	})

	t.Run("should render correct font-weight in CSS style values on accordion-title", func(t *testing.T) {
		input := `
    <mjml>
      <mj-head>
        <mj-attributes>
          <mj-accordion border="none" padding="1px" />
          <mj-accordion-element icon-wrapped-url="https://i.imgur.com/Xvw0vjq.png" icon-unwrapped-url="https://i.imgur.com/KKHenWa.png" icon-height="24px" icon-width="24px" />
          <mj-accordion-title font-family="Roboto, Open Sans, Helvetica, Arial, sans-serif" background-color="#fff" color="#031017" padding="15px" font-size="18px" />
          <mj-accordion-text font-family="Open Sans, Helvetica, Arial, sans-serif" background-color="#fafafa" padding="15px" color="#505050" font-size="14px" />
        </mj-attributes>
      </mj-head>
      <mj-body>
        <mj-section padding="20px" background-color="#ffffff">
          <mj-column background-color="#dededd">
            <mj-accordion>
              <mj-accordion-element>
                <mj-accordion-title font-weight="bold" css-class="accordion-title">Why use an accordion?</mj-accordion-title>
                <mj-accordion-text font-weight="bold">
                  <span style="line-height:20px">
                    Because emails with a lot of content are most of the time a very bad experience on mobile, mj-accordion comes handy when you want to deliver a lot of information in a concise way.
                  </span>
                </mj-accordion-text>
              </mj-accordion-element>
              <mj-accordion-element>
                <mj-accordion-title font-weight="700" css-class="accordion-title">How it works</mj-accordion-title>
                <mj-accordion-text font-weight="700">
                  <span style="line-height:20px">
                    Content is stacked into tabs and users can expand them at will. If responsive styles are not supported (mostly on desktop clients), tabs are then expanded and your content is readable at once.
                  </span>
                </mj-accordion-text>
              </mj-accordion-element>
            </mj-accordion>
          </mj-column>
        </mj-section>
      </mj-body>
    </mjml>
`

		html, err := mjml.MJML2Html(input, &core.MJMLOptions{
			KeepComments:    true,
			IgnoreIncludes:  true,
			ValidationLevel: core.ValidationLevelSoft,
		})

		assert.Nil(t, err)

		nodes := parseDOM(html)

		matches, err := query.SelectAll(".accordion-title", nodes, nil)
		assert.Nil(t, err)
		assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string {
			style := n.Attribs.GetOrDefault("style", "")
			start := strings.Index(style, "font-weight:") + 12
			end := strings.Index(style[start:], ";")
			return style[start : start+end]
		}), []string{"bold", "700"})
	})

	t.Run("should render correct font-family in CSS style values on accordion-title and accordion-text", func(t *testing.T) {
		input := `
    <mjml>
      <mj-body>
        <mj-section>
          <mj-column>
            <mj-accordion css-class="my-accordion-1" font-family="serif">
              <mj-accordion-element>
                <mj-accordion-title>Why use an accordion?</mj-accordion-title>
                <mj-accordion-text>
                    Because emails with a lot of content are most of the time a very bad experience on mobile, mj-accordion comes handy when you want to deliver a lot of information in a concise way.
                </mj-accordion-text>
              </mj-accordion-element>
            </mj-accordion>
          </mj-column>
        </mj-section>
        <mj-section>
          <mj-column>
            <mj-accordion css-class="my-accordion-2" font-family="serif">
              <mj-accordion-element font-family="sans-serif">
                <mj-accordion-title font-family="monospace">Why use an accordion?</mj-accordion-title>
                <mj-accordion-text font-family="monospace">
                    Because emails with a lot of content are most of the time a very bad experience on mobile, mj-accordion comes handy when you want to deliver a lot of information in a concise way.
                </mj-accordion-text>
              </mj-accordion-element>
            </mj-accordion>
          </mj-column>
        </mj-section>
      </mj-body>
    </mjml>
`

		html, err := mjml.MJML2Html(input, &core.MJMLOptions{
			KeepComments:    true,
			IgnoreIncludes:  true,
			ValidationLevel: core.ValidationLevelStrict,
		})

		assert.Nil(t, err)

		nodes := parseDOM(html)

		matches, err := query.SelectAll(".my-accordion-1 .mj-accordion-title td:first-child, .my-accordion-1 .mj-accordion-content td:first-child, .my-accordion-2 .mj-accordion-title td:first-child, .my-accordion-2 .mj-accordion-content td:first-child", nodes, nil)
		assert.Nil(t, err)
		assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string {
			style := n.Attribs.GetOrDefault("style", "")
			start := strings.Index(style, "font-family:") + 12
			end := strings.Index(style[start:], ";")
			return style[start : start+end]
		}), []string{"serif", "serif", "monospace", "monospace"})
	})

	t.Run("should render correct padding in CSS style values on accordion-title and accordion-text", func(t *testing.T) {
		input := `
    <mjml>
      <mj-body>
        <mj-section>
          <mj-column>
            <mj-accordion>
              <mj-accordion-element>
                <mj-accordion-title padding="20px" padding-bottom="40px" padding-left="40px" padding-right="40px" padding-top="40px">Why use an accordion?</mj-accordion-title>
                <mj-accordion-text padding="20px" padding-bottom="40px" padding-left="40px" padding-right="40px" padding-top="40px">
                    Because emails with a lot of content are most of the time a very bad experience on mobile, mj-accordion comes handy when you want to deliver a lot of information in a concise way.
                </mj-accordion-text>
              </mj-accordion-element>
            </mj-accordion>
          </mj-column>
        </mj-section>
      </mj-body>
    </mjml>
`

		html, err := mjml.MJML2Html(input, &core.MJMLOptions{
			KeepComments:    true,
			IgnoreIncludes:  true,
			ValidationLevel: core.ValidationLevelStrict,
		})

		assert.Nil(t, err)

		nodes := parseDOM(html)

		paddings := []string{
			"padding-left",
			"padding-right",
			"padding-top",
			"padding-bottom",
		}

		for _, padding := range paddings {
			matches, err := query.SelectAll(".mj-accordion-title td:first-child, .mj-accordion-content td:first-child", nodes, nil)
			assert.Nil(t, err)
			assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string {
				style := n.Attribs.GetOrDefault("style", "")
				start := strings.Index(style, padding+":") + len(padding) + 1
				end := strings.Index(style[start:], ";")
				return style[start : start+end]
			}), []string{"40px", "40px"})
		}
	})

	t.Run("should render correct display in CSS style values on mj-carousel-thumbnail", func(t *testing.T) {
		input := `
    <mjml>
      <mj-body>
        <mj-section>
          <mj-column>
            <mj-carousel thumbnails="supported">
              <mj-carousel-image src="https://placehold.co/450x300/333/ccc/png" />
              <mj-carousel-image src="https://placehold.co/450x300/ccc/000/png" />
              <mj-carousel-image src="https://placehold.co/450x300/f45e43/fff/png" />
            </mj-carousel>
          </mj-column>
        </mj-section>
      </mj-body>
    </mjml>
`

		html, err := mjml.MJML2Html(input, &core.MJMLOptions{
			KeepComments:    true,
			IgnoreIncludes:  true,
			ValidationLevel: core.ValidationLevelStrict,
		})

		assert.Nil(t, err)

		nodes := parseDOM(html)

		matches, err := query.SelectAll(".mj-carousel-thumbnail", nodes, nil)
		assert.Nil(t, err)
		assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string {
			style := n.Attribs.GetOrDefault("style", "")
			start := strings.Index(style, "display:") + 8
			end := strings.Index(style[start:], ";")
			return style[start : start+end]
		}), []string{"none", "none", "none"})
	})

	t.Run("should render correct border-radius / inner-border-radius (and border-collapse) in CSS style values on mj-column", func(t *testing.T) {
		input := `
    <mjml>
      <mj-body>
        <mj-section>
          <mj-column border-radius="50px" inner-border-radius="40px" padding="50px" border="5px solid #000" inner-border="5px solid #666">
            <mj-text>Hello World</mj-text>
          </mj-column>
        </mj-section>
      </mj-body>
    </mjml>
`

		html, err := mjml.MJML2Html(input, &core.MJMLOptions{
			KeepComments:    true,
			IgnoreIncludes:  true,
			ValidationLevel: core.ValidationLevelStrict,
		})

		assert.Nil(t, err)

		nodes := parseDOM(html)

		matches, err := query.SelectAll(".mj-column-per-100 > table > tbody > tr > td, .mj-column-per-100 > table > tbody > tr > td > table", nodes, nil)
		assert.Nil(t, err)
		assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string {
			style := n.Attribs.GetOrDefault("style", "")
			start := strings.Index(style, "border-radius:") + 14
			end := strings.Index(style[start:], ";")
			return style[start : start+end]
		}), []string{"50px", "40px"})

		matches, err = query.SelectAll(".mj-column-per-100 > table > tbody > tr > td, .mj-column-per-100 > table > tbody > tr > td > table", nodes, nil)
		assert.Nil(t, err)
		assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string {
			style := n.Attribs.GetOrDefault("style", "")
			start := strings.Index(style, "border-collapse:") + 16
			end := strings.Index(style[start:], ";")
			return style[start : start+end]
		}), []string{"separate", "separate"})
	})

	t.Run("should not alter the whitespace between the opening/closing comment tags and the comment content", func(t *testing.T) {
		input := `
    <mjml>
      <mj-body>
        <mj-section>
          <mj-column>
            <mj-text>
            <p>View source to see comments below</p>
            <!-- comment with standard spaces -->
            <br>
            <!--comment without spaces-->
            <br>
            <!--     comment with 5 spaces     -->
            </mj-text>
          </mj-column>
        </mj-section>
      </mj-body>
    </mjml>
`

		html, err := mjml.MJML2Html(input, &core.MJMLOptions{
			KeepComments:    true,
			IgnoreIncludes:  true,
			ValidationLevel: core.ValidationLevelStrict,
		})

		assert.Nil(t, err)

		// should not alter templating syntax, or move the content that is outside any tag (mj-raws)
		expected := []string{
			"<!-- comment with standard spaces -->",
			"<!--comment without spaces-->",
			"<!--     comment with 5 spaces     -->",
		}

		indexes := utils.MapFunc(expected, func(s string) int { return strings.Index(html, s) })

		assert.NotContains(t, indexes, -1)
	})

	t.Run("should render correct padding in CSS style values on navbar hamburger icon", func(t *testing.T) {
		input := `
    <mjml>
      <mj-body>
        <mj-section>
          <mj-column>
            <mj-navbar hamburger="hamburger" ico-padding="20px" ico-padding-bottom="20px" ico-padding-left="30px" ico-padding-right="40px"  ico-padding-top="50px" >
                <mj-navbar-link href="/gettings-started-onboard" color="#ffffff">Getting started</mj-navbar-link>
                <mj-navbar-link href="/try-it-live" color="#ffffff">Try it live</mj-navbar-link>
            </mj-navbar>
          </mj-column>
        </mj-section>
      </mj-body>
    </mjml>
`

		html, err := mjml.MJML2Html(input, &core.MJMLOptions{
			KeepComments:    true,
			IgnoreIncludes:  true,
			ValidationLevel: core.ValidationLevelStrict,
		})

		assert.Nil(t, err)

		nodes := parseDOM(html)

		paddings := []string{
			"padding-bottom",
			"padding-left",
			"padding-right",
			"padding-top",
		}

		// Padding should be ['20px', '30px', '40px', '50px']
		expected := [][]string{
			{"20px"},
			{"30px"},
			{"40px"},
			{"50px"},
		}

		for idx, padding := range paddings {
			matches, err := query.SelectAll(".mj-menu-label", nodes, nil)
			assert.Nil(t, err)
			assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string {
				style := n.Attribs.GetOrDefault("style", "")
				start := strings.Index(style, padding+":") + len(padding) + 1
				end := strings.Index(style[start:], ";")
				return style[start : start+end]
			}), expected[idx])
		}
	})

	t.Run("should render correct align in CSS style values on mj-social-element", func(t *testing.T) {
		input := `
    <mjml>
      <mj-body>
        <mj-section>
          <mj-column>
            <mj-social mode="vertical">
              <mj-social-element name="facebook" href="https://mjml.io/" icon-position="right" align="right"css-class="my-social-element">
                Facebook
              </mj-social-element>
            </mj-social>
          </mj-column>
        </mj-section>
      </mj-body>
    </mjml>
`

		html, err := mjml.MJML2Html(input, &core.MJMLOptions{
			KeepComments:    true,
			IgnoreIncludes:  true,
			ValidationLevel: core.ValidationLevelStrict,
		})

		assert.Nil(t, err)

		nodes := parseDOM(html)

		matches, err := query.SelectAll(".my-social-element > td:first-child", nodes, nil)
		assert.Nil(t, err)
		assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string {
			style := n.Attribs.GetOrDefault("style", "")
			start := strings.Index(style, "text-align:") + 11
			end := strings.Index(style[start:], ";")
			return style[start : start+end]
		}), []string{"right"})
	})

	t.Run("should render correct icon-height align in CSS style values on mj-social", func(t *testing.T) {
		input := `
    <mjml>
      <mj-body>
        <mj-section>
          <mj-column css-class="my-social-element">
            <mj-social icon-height="40px">
              <mj-social-element name="facebook" href="https://mjml.io/" css-class="my-social-element">
                Facebook
              </mj-social-element>
            </mj-social>
          </mj-column>
        </mj-section>
      </mj-body>
    </mjml>
`

		html, err := mjml.MJML2Html(input, &core.MJMLOptions{
			KeepComments:    true,
			IgnoreIncludes:  true,
			ValidationLevel: core.ValidationLevelStrict,
		})

		assert.Nil(t, err)

		nodes := parseDOM(html)

		matches, err := query.SelectAll(".my-social-element > td > table > tbody > tr > td", nodes, nil)
		assert.Nil(t, err)
		assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string {
			style := n.Attribs.GetOrDefault("style", "")
			start := strings.Index(style, "height:") + 7
			end := strings.Index(style[start:], ";")
			return style[start : start+end]
		}), []string{"40px"})

		matches, err = query.SelectAll(".my-social-element > td > table > tbody > tr > td img", nodes, nil)
		assert.Nil(t, err)
		assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string {
			return n.Attribs.GetOrDefault("height", "")
		}), []string{""})
	})

	t.Run("should render correct cellspacing (and border-collapse) in HTML tag / CSS style values on mj-table", func(t *testing.T) {
		input := `
    <mjml>
      <mj-body>
        <mj-section>
          <mj-column>
            <mj-table border="1px solid #000" width="auto" cellpadding="20" cellspacing="10" css-class="my-table">
              <tr style="border-bottom:1px solid #000;text-align:left;">
                <th style="background:#ddd;">Year</th>
                <th style="background:#ddd;">Language</th>
                <th style="background:#ddd;">Inspired from</th>
              </tr>
              <tr>
                <td style="background:#ddd;">1995</td>
                <td style="background:#ddd;">PHP</td>
                <td style="background:#ddd;">C, Shell Unix</td>
              </tr>
            </mj-table>
          </mj-column>
        </mj-section>
      </mj-body>
    </mjml>
`

		html, err := mjml.MJML2Html(input, &core.MJMLOptions{
			KeepComments:    true,
			IgnoreIncludes:  true,
			ValidationLevel: core.ValidationLevelStrict,
		})

		assert.Nil(t, err)

		nodes := parseDOM(html)

		matches, err := query.SelectAll(".my-table > table", nodes, nil)
		assert.Nil(t, err)
		assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string {
			return n.Attribs.GetOrDefault("cellspacing", "")
		}), []string{"10"})

		matches, err = query.SelectAll(".my-table > table", nodes, nil)
		assert.Nil(t, err)
		assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string {
			style := n.Attribs.GetOrDefault("style", "")
			start := strings.Index(style, "border-collapse:") + 16
			end := strings.Index(style[start:], ";")
			return style[start : start+end]
		}), []string{"separate"})
	})

	t.Run("should render correct width in CSS style values on mj-table", func(t *testing.T) {
		input := `
    <mjml>
        <mj-body>
            <mj-wrapper>
                <mj-section>
                    <mj-column>
                        <mj-table css-class="table">
                            <tr>
                                <th style="border: 1px solid black;text-align: left;">
                                    Default Width
                                </th>
                                <td style="border: 1px solid black;">
                                    100%
                                </td>
                            </tr>
                        </mj-table>
                    </mj-column>
                </mj-section>
                <mj-section>
                    <mj-column>
                        <mj-table width="500px" css-class="table">
                            <tr>
                                <th style="border: 1px solid black;text-align: left;">
                                    Pixel Width
                                </th>
                                <td style="border: 1px solid black;">
                                    500px
                                </td>
                            </tr>
                        </mj-table>
                    </mj-column>
                </mj-section>
                <mj-section>
                    <mj-column>
                        <mj-table width="80%" css-class="table">
                            <tr>
                                <th style="border: 1px solid black;text-align: left;">
                                    Percentage Width
                                </th>
                                <td style="border: 1px solid black;">
                                    80%
                                </td>
                            </tr>
                        </mj-table>
                    </mj-column>
                </mj-section>
                <mj-section css-class="section">
                    <mj-column>
                        <mj-table width="auto" css-class="table">
                            <tr>
                                <th style="border: 1px solid black;text-align: left;">
                                    Auto Width
                                </th>
                                <td style="border: 1px solid black;">
                                    Auto
                                </td>
                            </tr>
                        </mj-table>
                    </mj-column>
                </mj-section>
            </mj-wrapper>
        </mj-body>
    </mjml>
`

		html, err := mjml.MJML2Html(input, &core.MJMLOptions{
			KeepComments:    true,
			IgnoreIncludes:  true,
			ValidationLevel: core.ValidationLevelStrict,
		})

		assert.Nil(t, err)

		nodes := parseDOM(html)

		matches, err := query.SelectAll(".table table", nodes, nil)
		assert.Nil(t, err)
		assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string {
			return n.Attribs.GetOrDefault("width", "")
		}), []string{"100%", "500", "80%", "auto"})

		matches, err = query.SelectAll(".table table", nodes, nil)
		assert.Nil(t, err)
		assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string {
			style := n.Attribs.GetOrDefault("style", "")
			start := strings.Index(style, "width:") + 6
			end := strings.Index(style[start:], ";")
			return style[start : start+end]
		}), []string{"100%", "500px", "80%", "auto"})
	})

	t.Run("should render correct border-radius (and border-collapse) in CSS style values on mj-wrapper and mj-section", func(t *testing.T) {
		input := `
    <mjml>
      <mj-body>
        <mj-wrapper border="1px solid red" border-radius="10px">
          <mj-section>
            <mj-column>
              <mj-text font-size="20px" color="#F45E43" font-family="helvetica">Hello World</mj-text>
            </mj-column>
          </mj-section>
        </mj-wrapper>
      </mj-body>
    </mjml>
`

		html, err := mjml.MJML2Html(input, &core.MJMLOptions{
			KeepComments:    true,
			IgnoreIncludes:  true,
			ValidationLevel: core.ValidationLevelStrict,
		})

		assert.Nil(t, err)

		nodes := parseDOM(html)

		matches, err := query.SelectAll("body > div > div > table:first-child > tbody > tr > td, body > div > div", nodes, nil)
		assert.Nil(t, err)
		assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string {
			style := n.Attribs.GetOrDefault("style", "")
			start := strings.Index(style, "border-radius:") + 14
			end := strings.Index(style[start:], ";")
			return style[start : start+end]
		}), []string{"10px", "10px"})

		matches, err = query.SelectAll("body > div > div", nodes, nil)
		assert.Nil(t, err)
		assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string {
			style := n.Attribs.GetOrDefault("style", "")
			start := strings.Index(style, "overflow:") + 9
			end := strings.Index(style[start:], ";")
			return style[start : start+end]
		}), []string{"hidden"})

		matches, err = query.SelectAll("body > div > div > table:first-child", nodes, nil)
		assert.Nil(t, err)
		assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string {
			style := n.Attribs.GetOrDefault("style", "")
			start := strings.Index(style, "border-collapse:") + 16
			end := strings.Index(style[start:], ";")
			return style[start : start+end]
		}), []string{"separate"})
	})

	t.Run("should render correct gap values in CSS style values on children mj-section", func(t *testing.T) {
		input := `
    <mjml>
      <mj-body>
        <mj-wrapper gap="20px" css-class="my-wrapper" background-color="#000"> 
          <mj-section css-class="my-section" background-color="#f45e43" padding="10px">
            <mj-column>
              <mj-text>Section 1</mj-text>
            </mj-column>
          </mj-section>
          <mj-section css-class="my-section" background-color="#ccc" padding="10px">
            <mj-column>
              <mj-text>Section 2</mj-text>
            </mj-column>
          </mj-section>
          <mj-section css-class="my-section" background-color="#333" padding="10px">
            <mj-column>
              <mj-text color="#fff">Section 3</mj-text>
            </mj-column>
          </mj-section>
        </mj-wrapper>
      </mj-body>
    </mjml>
`

		html, err := mjml.MJML2Html(input, &core.MJMLOptions{
			KeepComments:    true,
			IgnoreIncludes:  true,
			ValidationLevel: core.ValidationLevelStrict,
		})

		assert.Nil(t, err)

		nodes := parseDOM(html)

		matches, err := query.SelectAll(".my-section", nodes, nil)
		assert.Nil(t, err)
		assert.Equal(t, utils.MapFunc(matches, func(n *dom.Node) string {
			style := n.Attribs.GetOrDefault("style", "")
			substr := "margin-top:"

			if strings.Contains(style, substr) {
				start := strings.Index(style, substr) + 11
				end := strings.Index(style[start:], ";")
				return style[start : start+end]
			} else {
				return ""
			}
		}), []string{"", "20px", "20px"})
	})
}
