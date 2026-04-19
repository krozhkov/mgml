package main

import (
	"os"
	"unsafe"

	"github.com/krozhkov/mgml/core"
	"github.com/krozhkov/mgml/mjml"
)

func main() {
	input := os.Args[1]
	output := os.Args[2]

	content, err := os.ReadFile(input)
	if err != nil {
		panic(err)
	}

	str := unsafe.String(unsafe.SliceData(content), len(content))
	html, err := mjml.MJML2Html(str, &core.MJMLOptions{
		KeepComments:    true,
		ValidationLevel: core.ValidationLevelSoft,
	})

	bytes := unsafe.Slice(unsafe.StringData(html), len(html))
	err = os.WriteFile(output, bytes, 0o644)
	if err != nil {
		panic(err)
	}
}
