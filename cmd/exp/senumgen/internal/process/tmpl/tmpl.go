package tmpl

import (
	"bytes"
	"text/template"
)

type TemplateCollector struct {
	buf bytes.Buffer
	err error
}

func NewTemplateCollector() *TemplateCollector {
	var buf bytes.Buffer
	return &TemplateCollector{
		buf: buf,
		err: nil,
	}
}

func (tc *TemplateCollector) Merge(tmpl string, args interface{}) *TemplateCollector {
	if tc.err != nil {
		return tc
	}

	t, err := template.New("tmpl").Funcs(
		template.FuncMap{
			"capitalize": capitalize,
			"bracket":    Bracket,
			"csvConnect": csvConnect,
		},
	).Parse(tmpl)
	if err != nil {
		tc.err = err
	}

	if err := t.Execute(&tc.buf, args); err != nil {
		tc.err = err
	}

	return tc
}

func (tc *TemplateCollector) Export() ([]byte, error) {
	if tc.err != nil {
		return nil, tc.err
	}
	return tc.buf.Bytes(), nil
}
