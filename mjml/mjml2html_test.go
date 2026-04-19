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
}
