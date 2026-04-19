package cssparser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func parse(input string) ([]*Rule, error) {
	parser := NewCssParser()
	return parser.Parse(input)
}

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
			rules, err := parse(selector)
			assert.Nil(t, err)
			assert.Equal(t, 1, len(rules))
		})
	}
}

func TestCssParser(t *testing.T) {
	t.Run("Should parse attribute selectors", func(t *testing.T) {
		rules, err := parse(`input[type=range] {
			background: lightblue url("img_tree.gif") no-repeat fixed center;
			accent-color: rgb(0, 0, 255);
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(rules))
	})

	t.Run("Allow hash tokens inside property value", func(t *testing.T) {
		rules, err := parse(`.icon {
			background: #ffffff;
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(rules))
	})

	t.Run("Allow hash tokens inside property value", func(t *testing.T) {
		rules, err := parse(`.icon {
			background: #ffffff;
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(rules))
	})

	t.Run("Should parse !important inside property value", func(t *testing.T) {
		rules, err := parse(`/* This is a comment */
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
		rules, err := parse(`
		:root {
			--primary-bg-color: #1e90ff;
			--primary-color: #ffffff;
		}

		body {
			background-color: var(--primary-bg-color);
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(rules))

		rules, err = parse(`
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
		rules, err := parse(`h1 {
			font-size: calc(1.5rem + 3vw);
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(rules))

		rules, err = parse(`
		.foo {
			color: lch(from aquamarine l c calc(h + 180))
		}`)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(rules))
	})
}
