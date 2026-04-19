package parser

import (
	"testing"

	"github.com/elliotchance/orderedmap/v3"
	"github.com/stretchr/testify/assert"
)

type mjmlComponent struct {
	tagName string
}

func (c *mjmlComponent) GetTagName() string {
	return c.tagName
}
func (c *mjmlComponent) IsEndingTag() bool {
	return true
}

var endingComponents = []MJMLComponent{
	{TagName: "mj-breakpoint", EndingTag: true},
	{TagName: "mj-preview", EndingTag: true},
	{TagName: "mj-style", EndingTag: true},
	{TagName: "mj-title", EndingTag: true},
	{TagName: "mj-button", EndingTag: true},
	{TagName: "mj-raw", EndingTag: true},
	{TagName: "mj-text", EndingTag: true},
	{TagName: "mj-table", EndingTag: true},
	{TagName: "mj-social-element", EndingTag: true},
	{TagName: "mj-navbar-link", EndingTag: true},
	{TagName: "mj-accordion-text", EndingTag: true},
	{TagName: "mj-accordion-title", EndingTag: true},
	{TagName: "mj-carousel-image", EndingTag: true},
}

func prepareNode(node *MJMLNode) *MJMLNode {
	node.AbsoluteFilePath = ""
	node.File = "."

	for _, inc := range node.IncludedIn {
		inc.File = "."
	}

	for _, child := range node.Children {
		prepareNode(child)
	}

	return node
}

func parse(mjml string) (*MJMLNode, error) {
	return MJMLParser([]byte(mjml), MJMLParserOptions{
		KeepComments: true,
		Components:   endingComponents,
		FilePath:     ".",
	}, nil)
}

type parserTest struct {
	name   string
	input  string
	output *MJMLNode
}

var testValues = []parserTest{
	{
		name: "Special characters",
		input: `
<mjml>
  <mj-body>
    <mj-section background-color="#CCCCCC" full-width="full-width">
      <mj-button href="<% dynamic %>" pouf="$2">
        &end<br />
        Blu & end $1
        &amp;
        lorem
      </mj-button>
      <mj-button href="https://mjml.io?encodedUrl=https%3A%2F%2Fmjml.io&coin=coi">
        Blu
      </mj-button>
      <mj-button href="&é(§&è!çà)">https%3A%2F%2Fmjml.io</mj-button>
      <mj-raw>
        <coin color="#CCCCCC">bla</coin>
      </mj-raw>
    </mj-section>
  </mj-body>
</mjml>
`,
		output: &MJMLNode{
			File:    ".",
			Line:    2,
			TagName: "mjml",
			Children: []*MJMLNode{
				{
					File:    ".",
					Line:    3,
					TagName: "mj-body",
					Children: []*MJMLNode{
						{
							File:    ".",
							Line:    4,
							TagName: "mj-section",
							Attributes: orderedmap.NewOrderedMapWithElements(
								&Attribute{Key: "background-color", Value: "#CCCCCC"},
								&Attribute{Key: "full-width", Value: "full-width"},
							),
							Children: []*MJMLNode{
								{
									File:    ".",
									Line:    5,
									TagName: "mj-button",
									Attributes: orderedmap.NewOrderedMapWithElements(
										&Attribute{Key: "href", Value: "<% dynamic %>"},
										&Attribute{Key: "pouf", Value: "$2"},
									),
									Content: "&end<br />\n        Blu & end $1\n        &amp;\n        lorem",
								},
								{
									File:    ".",
									Line:    11,
									TagName: "mj-button",
									Attributes: orderedmap.NewOrderedMapWithElements(
										&Attribute{Key: "href", Value: "https://mjml.io?encodedUrl=https%3A%2F%2Fmjml.io&coin=coi"},
									),
									Content: "Blu",
								},
								{
									File:    ".",
									Line:    14,
									TagName: "mj-button",
									Attributes: orderedmap.NewOrderedMapWithElements(
										&Attribute{Key: "href", Value: "&é(§&è!çà)"},
									),
									Content: "https%3A%2F%2Fmjml.io",
								},
								{
									File:    ".",
									Line:    15,
									TagName: "mj-raw",
									Content: "<coin color=\"#CCCCCC\">bla</coin>",
								},
							},
						},
					},
				},
			},
		},
	},
	{
		name: "Similar tags",
		input: `
<mjml>
  <mj-body>
    <mj-text-test-wrapper>
      <mj-text>MJML</mj-text>
      <mj-text attr="val">FTW</mj-text>
    </mj-text-test-wrapper>
    <mj-text-test-wrapper>
      <mj-text attr="val">FTW</mj-text>
      <mj-text>MJML</mj-text>
    </mj-text-test-wrapper>
  </mj-body>
</mjml>
`,
		output: &MJMLNode{
			File:    ".",
			Line:    2,
			TagName: "mjml",
			Children: []*MJMLNode{
				{
					File:    ".",
					Line:    3,
					TagName: "mj-body",
					Children: []*MJMLNode{
						{
							File:    ".",
							Line:    4,
							TagName: "mj-text-test-wrapper",
							Children: []*MJMLNode{
								{
									File:    ".",
									Line:    5,
									TagName: "mj-text",
									Content: "MJML",
								},
								{
									File:    ".",
									Line:    6,
									TagName: "mj-text",
									Attributes: orderedmap.NewOrderedMapWithElements(
										&Attribute{Key: "attr", Value: "val"},
									),
									Content: "FTW",
								},
							},
						},
						{
							File:    ".",
							Line:    8,
							TagName: "mj-text-test-wrapper",
							Children: []*MJMLNode{
								{
									File:    ".",
									Line:    9,
									TagName: "mj-text",
									Attributes: orderedmap.NewOrderedMapWithElements(
										&Attribute{Key: "attr", Value: "val"},
									),
									Content: "FTW",
								},
								{
									File:    ".",
									Line:    10,
									TagName: "mj-text",
									Content: "MJML",
								},
							},
						},
					},
				},
			},
		},
	},
	{
		name: "Self closing tags",
		input: `
<mjml>
  <mj-head>
    <mj-attributes>
      <mj-text color="blue" />
      <mj-text font-size="40px" />
    </mj-attributes>
  </mj-head>
  <mj-body>
    <mj-section>
      <mj-column>
        <mj-text>
          Hello !
        </mj-text>
      </mj-column>
    </mj-section>
  </mj-body>
</mjml>
`,
		output: &MJMLNode{
			File:    ".",
			Line:    2,
			TagName: "mjml",
			Children: []*MJMLNode{
				{
					File:    ".",
					Line:    3,
					TagName: "mj-head",
					Children: []*MJMLNode{
						{
							File:    ".",
							Line:    4,
							TagName: "mj-attributes",
							Children: []*MJMLNode{
								{
									File:    ".",
									Line:    5,
									TagName: "mj-text",
									Attributes: orderedmap.NewOrderedMapWithElements(
										&Attribute{Key: "color", Value: "blue"},
									),
								},
								{
									File:    ".",
									Line:    6,
									TagName: "mj-text",
									Attributes: orderedmap.NewOrderedMapWithElements(
										&Attribute{Key: "font-size", Value: "40px"},
									),
								},
							},
						},
					},
				},
				{
					File:    ".",
					Line:    9,
					TagName: "mj-body",
					Children: []*MJMLNode{
						{
							File:    ".",
							Line:    10,
							TagName: "mj-section",
							Children: []*MJMLNode{
								{
									File:    ".",
									Line:    11,
									TagName: "mj-column",
									Children: []*MJMLNode{
										{
											File:    ".",
											Line:    12,
											TagName: "mj-text",
											Content: "Hello !",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
	// Input that matches most of the CDATAs regex but not all, potentially resulting in regex timeout
	{
		name: "Regex timeout",
		input: `
<mj-section>
  <mj-text font-family="Arial" />
  <mj-column background-color="#ffffff" css-class="column1"></mj-column>
</mj-section>
`,
		output: &MJMLNode{
			File:    ".",
			Line:    2,
			TagName: "mj-section",
			Children: []*MJMLNode{
				{
					File:    ".",
					Line:    3,
					TagName: "mj-text",
					Attributes: orderedmap.NewOrderedMapWithElements(
						&Attribute{Key: "font-family", Value: "Arial"},
					),
				},
				{
					File:    ".",
					Line:    4,
					TagName: "mj-column",
					Attributes: orderedmap.NewOrderedMapWithElements(
						&Attribute{Key: "background-color", Value: "#ffffff"},
						&Attribute{Key: "css-class", Value: "column1"},
					),
				},
			},
		},
	},
	{
		name: "Multiline attributes",
		input: `
<mj-text
    padding-left="16px"

    padding-right="16px">
    <a href="https://www.test.com" style="color: #60788c">View blog ]]post</a>
</mj-text>
`,
		output: &MJMLNode{
			File:    ".",
			Line:    2,
			TagName: "mj-text",
			Attributes: orderedmap.NewOrderedMapWithElements(
				&Attribute{Key: "padding-left", Value: "16px"},
				&Attribute{Key: "padding-right", Value: "16px"},
			),
			Content: "<a href=\"https://www.test.com\" style=\"color: #60788c\">View blog ]]post</a>",
		},
	},
	{
		name: "Self closing Ending Tags",
		input: `
      <mjml>
        <mj-head>
          <mj-title></mj-title>
          <mj-attributes>
            <mj-text font-size="27px" />
          </mj-attributes>
        </mj-head>
        <mj-body>
          <mj-section>
            <mj-column width="65%">
              <mj-text mj-class="small" align="left" font-family="Helvetica" color="#000000" padding-top="20px">
                coin
                <a href="https://test" style="text-decoration:underline;color:#336666;font-weight:bold" class="mobile-small-letters">Majors and Minors</a>
                bla
                <a href="https://test" style="text-decoration:underline;color:#336666;font-weight:bold" class="mobile-small-letters">Majors and Minors</a>
                <mj-raw>
                  coin
                </mj-raw>
              </mj-text>
            </mj-column>
          </mj-section>
        </mj-body>
      </mjml>
`,
		output: &MJMLNode{
			File:    ".",
			Line:    2,
			TagName: "mjml",
			Children: []*MJMLNode{
				{
					File:    ".",
					Line:    3,
					TagName: "mj-head",
					Children: []*MJMLNode{
						{
							File:    ".",
							Line:    4,
							TagName: "mj-title",
						},
						{
							File:    ".",
							Line:    5,
							TagName: "mj-attributes",
							Children: []*MJMLNode{
								{
									File:    ".",
									Line:    6,
									TagName: "mj-text",
									Attributes: orderedmap.NewOrderedMapWithElements(
										&Attribute{Key: "font-size", Value: "27px"},
									),
								},
							},
						},
					},
				},
				{
					File:    ".",
					Line:    9,
					TagName: "mj-body",
					Children: []*MJMLNode{
						{
							File:    ".",
							Line:    10,
							TagName: "mj-section",
							Children: []*MJMLNode{
								{
									File:    ".",
									Line:    11,
									TagName: "mj-column",
									Attributes: orderedmap.NewOrderedMapWithElements(
										&Attribute{Key: "width", Value: "65%"},
									),
									Children: []*MJMLNode{
										{
											File:    ".",
											Line:    12,
											TagName: "mj-text",
											Attributes: orderedmap.NewOrderedMapWithElements(
												&Attribute{Key: "mj-class", Value: "small"},
												&Attribute{Key: "align", Value: "left"},
												&Attribute{Key: "font-family", Value: "Helvetica"},
												&Attribute{Key: "color", Value: "#000000"},
												&Attribute{Key: "padding-top", Value: "20px"},
											),
											Content: "coin\n                <a href=\"https://test\" style=\"text-decoration:underline;color:#336666;font-weight:bold\" class=\"mobile-small-letters\">Majors and Minors</a>\n                bla\n                <a href=\"https://test\" style=\"text-decoration:underline;color:#336666;font-weight:bold\" class=\"mobile-small-letters\">Majors and Minors</a>\n                <mj-raw>\n                  coin\n                </mj-raw>",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		name: "Single opening tag in endingTag, single and multi-line",
		input: `
      <mjml>
        <mj-body>
          <mj-section>
            <mj-column>
              <mj-raw test="test"><?php endif ?></mj-raw>
              <mj-raw>
                <?php endif ?>
              </mj-raw>
            </mj-column>
          </mj-section>
        </mj-body>
      </mjml>
`,
		output: &MJMLNode{
			File:    ".",
			Line:    2,
			TagName: "mjml",
			Children: []*MJMLNode{
				{
					File:    ".",
					Line:    3,
					TagName: "mj-body",
					Children: []*MJMLNode{
						{
							File:    ".",
							Line:    4,
							TagName: "mj-section",
							Children: []*MJMLNode{
								{
									File:    ".",
									Line:    5,
									TagName: "mj-column",
									Children: []*MJMLNode{
										{
											File:    ".",
											Line:    6,
											TagName: "mj-raw",
											Attributes: orderedmap.NewOrderedMapWithElements(
												&Attribute{Key: "test", Value: "test"},
											),
											Content: "<?php endif ?>",
										},
										{
											File:    ".",
											Line:    7,
											TagName: "mj-raw",
											Content: "<?php endif ?>",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		name: "Include",
		input: `
      <mjml>
        <mj-body>
          <mj-section>
            <mj-include path="./test/incl.mjml" />
          </mj-section>
        </mj-body>
      </mjml>
`,
		output: &MJMLNode{
			File:    ".",
			Line:    2,
			TagName: "mjml",
			Children: []*MJMLNode{
				{
					File:    ".",
					Line:    3,
					TagName: "mj-body",
					Children: []*MJMLNode{
						{
							File:    ".",
							Line:    4,
							TagName: "mj-section",
							Children: []*MJMLNode{
								{
									File:       ".",
									Line:       1,
									TagName:    "mj-column",
									IncludedIn: []*MJMLIncludedIn{{File: ".", Line: 5}},
									Children: []*MJMLNode{
										{
											File:       ".",
											Line:       2,
											TagName:    "mj-text",
											IncludedIn: []*MJMLIncludedIn{{File: ".", Line: 5}},
											Attributes: orderedmap.NewOrderedMapWithElements(
												&Attribute{Key: "font-size", Value: "22px"},
											),
											Content: "COIN\n    <a src=\"test\">aze</a>",
										},
										{
											File:       ".",
											Line:       6,
											TagName:    "mj-text",
											IncludedIn: []*MJMLIncludedIn{{File: ".", Line: 5}},
											Attributes: orderedmap.NewOrderedMapWithElements(
												&Attribute{Key: "font-size", Value: "22px"},
											),
											Content: "COIN2\n    <a src=\"test\">aze2</a>",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

func TestValues(t *testing.T) {
	for _, tt := range testValues {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parse(tt.input)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, prepareNode(result), tt.output)
		})
	}
}
