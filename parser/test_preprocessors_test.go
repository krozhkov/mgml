package parser

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/gkampitakis/go-snaps/snaps"
)

type TemplateButton struct {
	Title string
}

type TemplateData struct {
	Buttons []*TemplateButton
}

var xml string = `<mjml>
<mj-body>
  <mj-section mj-class="content">
    {{- range .Buttons }}
    <mj-text>{{ .Title }}</mj-text>
    {{- end }}
  </mj-section>
</mj-body>
</mjml>`

func TestPreprocessors(t *testing.T) {
	testPreprocessor := func(content []byte) []byte {
		data := TemplateData{
			Buttons: []*TemplateButton{
				{Title: "Title 1"},
				{Title: "Title 2"},
			},
		}

		tmpl, err := template.New("test").Parse(string(content))
		if err != nil {
			t.Error(err)
		}

		output := bytes.NewBufferString("")
		err = tmpl.Execute(output, data)
		if err != nil {
			t.Error(err)
		}

		return output.Bytes()
	}

	result, err := MJMLParser([]byte(xml), MJMLParserOptions{
		KeepComments:  true,
		Components:    endingComponents,
		Preprocessors: []MJMLPreprocessor{testPreprocessor},
	}, nil)

	if err != nil {
		t.Error(err)
	}

	snaps.MatchSnapshot(t, prepareNode(result))
}
