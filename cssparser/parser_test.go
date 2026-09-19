package cssparser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectors(t *testing.T) {
	selectors := []string{
		"myNameSpace|a {}",
		"*|p {}",
		".my\\(class {}",
		".intro {}",
		"#Lastname {}",
		".intro, #Lastname {}",
		"h1 {}",
		"h1, p {}",
		"div p {}",
		"div > p {}",
		"ul + p {}",
		"ul ~ table {}",
		"* {}",
		"p.myquote {}",
		"[id] {}",
		"[id=my-Address] {}",
		"[id$=ess] {}",
		"[id|=my] {}",
		"[id^=L] {}",
		"[title~=beautiful] {}",
		"[id*=s] {}",
		":checked {}",
		":disabled {}",
		":enabled {}",
		":empty {}",
		":focus {}",
		"p:first-child {}",
		"p::first-letter {}",
		"p::first-line {}",
		"p:first-of-type {}",
		"h1:hover {}",
		"input:in-range {}",
		"input:out-of-range {}",
		"input:invalid {}",
		"input:valid {}",
		"p:lang(it) {}",
		"p:last-child {}",
		"p:last-of-type {}",
		"tr:nth-child(even) {}",
		"tr:nth-child(odd) {}",
		"li:nth-child(1) {}",
		"li:nth-last-child(1) {}",
		"li:nth-of-type(2) {}",
		"li:nth-last-of-type(2) {}",
		"b:only-child {}",
		"h3:only-of-type {}",
		":root {}",
	}

	for _, selector := range selectors {
		t.Run(fmt.Sprintf("selector should be valid: %s", selector), func(t *testing.T) {
			rules, err := ParseCss(selector)
			assert.Nil(t, err)
			assert.Equal(t, 1, len(rules))
		})
	}
}

func TestCssParser(t *testing.T) {
	t.Run("Should parse attribute selectors", func(t *testing.T) {
		rules, err := ParseCss(`input[type=range] {
			background: lightblue url("img_tree.gif") no-repeat fixed center;
			accent-color: rgb(0, 0, 255);
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(rules))
	})

	t.Run("Allow hash tokens inside property value", func(t *testing.T) {
		rules, err := ParseCss(`.icon {
			background: #ffffff;
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(rules))
	})

	t.Run("Allow hash tokens inside property value", func(t *testing.T) {
		rules, err := ParseCss(`.icon {
			background: #ffffff;
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(rules))
	})

	t.Run("Should parse !important inside property value", func(t *testing.T) {
		rules, err := ParseCss(`/* This is a comment */
			p > a, a.hyperlink {
				color: blue;
				text-decoration: underline !important;
				font-weight: normal   !IMPORTANT    ;
			}
		`)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(rules))
	})

	t.Run("Allow css variables", func(t *testing.T) {
		rules, err := ParseCss(`
		:root {
			--primary-bg-color: #1e90ff;
			--primary-color: #ffffff;
		}

		body {
			background-color: var(--primary-bg-color);
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(rules))

		rules, err = ParseCss(`
		.foo {
			--width-a: 100px;
			--width-b: calc(var(--width-a) / 2);
			--width-c: calc(var(--width-b) / 2);
			width: var(--width-c);
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(rules))
	})

	t.Run("Allow use complex calc funcions", func(t *testing.T) {
		rules, err := ParseCss(`h1 {
			font-size: calc(1.5rem + 3vw);
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(rules))

		rules, err = ParseCss(`
		.foo {
			color: lch(from aquamarine l c calc(h + 180))
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(rules))
	})
}

func TestCssNesting(t *testing.T) {
	t.Run("Nested pseudo-class with ampersand", func(t *testing.T) {
		rules, err := ParseCss(`.card {
			color: red;
			&:hover {
				color: green;
			}
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(rules))
		assert.Equal(t, ".card", rules[0].Selector)
		assert.Equal(t, "color", rules[0].Declarations[0].Property)
		assert.Equal(t, "red", rules[0].Declarations[0].Value)

		assert.Equal(t, ".card:hover", rules[1].Selector)
		assert.Equal(t, "color", rules[1].Declarations[0].Property)
		assert.Equal(t, "green", rules[1].Declarations[0].Value)
	})

	t.Run("Nested compound class with ampersand", func(t *testing.T) {
		rules, err := ParseCss(`.btn {
			padding: 10px;
			&.primary {
				background: blue;
			}
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(rules))
		assert.Equal(t, ".btn", rules[0].Selector)
		assert.Equal(t, ".btn.primary", rules[1].Selector)
		assert.Equal(t, "background", rules[1].Declarations[0].Property)
		assert.Equal(t, "blue", rules[1].Declarations[0].Value)
	})

	t.Run("Nested element selector without ampersand", func(t *testing.T) {
		rules, err := ParseCss(`.card {
			padding: 20px;
			p {
				margin: 0;
			}
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(rules))
		assert.Equal(t, ".card", rules[0].Selector)
		assert.Equal(t, ".card p", rules[1].Selector)
		assert.Equal(t, "margin", rules[1].Declarations[0].Property)
		assert.Equal(t, "0", rules[1].Declarations[0].Value)
	})

	t.Run("Comma-separated parents and children", func(t *testing.T) {
		rules, err := ParseCss(`.a, .b {
			& .c, & .d {
				margin: 0;
			}
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(rules))
		assert.Equal(t, ".a,.b", rules[0].Selector)
		assert.Equal(t, ".a .c, .b .c, .a .d, .b .d", rules[1].Selector)
	})

	t.Run("Deep nesting (3 levels)", func(t *testing.T) {
		rules, err := ParseCss(`.nav {
			ul {
				li {
					color: black;
				}
			}
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 3, len(rules))
		assert.Equal(t, ".nav", rules[0].Selector)
		assert.Equal(t, ".nav ul", rules[1].Selector)
		assert.Equal(t, ".nav ul li", rules[2].Selector)
		assert.Equal(t, "color", rules[2].Declarations[0].Property)
		assert.Equal(t, "black", rules[2].Declarations[0].Value)
	})

	t.Run("Declarations before and after nested rule", func(t *testing.T) {
		rules, err := ParseCss(`.card {
			color: red;
			&:hover {
				color: green;
			}
			background: white;
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(rules))
		assert.Equal(t, ".card", rules[0].Selector)
		assert.Equal(t, 2, len(rules[0].Declarations))
		assert.Equal(t, "color", rules[0].Declarations[0].Property)
		assert.Equal(t, "red", rules[0].Declarations[0].Value)
		assert.Equal(t, "background", rules[0].Declarations[1].Property)
		assert.Equal(t, "white", rules[0].Declarations[1].Value)

		assert.Equal(t, ".card:hover", rules[1].Selector)
		assert.Equal(t, 1, len(rules[1].Declarations))
		assert.Equal(t, "color", rules[1].Declarations[0].Property)
		assert.Equal(t, "green", rules[1].Declarations[0].Value)
	})
}

func TestAtRules(t *testing.T) {
	t.Run("Should return error for At-rules", func(t *testing.T) {
		_, err := ParseCss(`@media (max-width: 600px) { body { color: red; } }`)
		assert.NotNil(t, err)
		assert.Equal(t, "parser does not support At-rules", err.Error())
	})
}
