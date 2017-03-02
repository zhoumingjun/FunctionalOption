package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"testing"
)

func TestTemplate(t *testing.T) {
	tmpl, err := template.New("model").Parse(tmplOptions)
	if err != nil {
		fmt.Println(err)
		return
	}

	options := []Option{
		Option{
			FieldName:  "Name",
			Type:       "string",
			OptionName: "Name",
		},
		Option{
			FieldName:  "Age",
			Type:       "string",
			OptionName: "Old",
		},
	}
	var b bytes.Buffer

	err = tmpl.Execute(&b, map[string]interface{}{
		"TypeName": "example",
		"Options":  options,
	})

	out, err := format.Source(b.Bytes())
	if err != nil {
		fmt.Println(" format", err)
		return
	}

	fmt.Println(string(out))
}
