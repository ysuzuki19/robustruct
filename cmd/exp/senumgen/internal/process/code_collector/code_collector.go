package tmpl

import (
	"bytes"
	"text/template"
)

type CodeCollector struct {
	buf bytes.Buffer
	err error
}

func NewCodeCollector() *CodeCollector {
	var buf bytes.Buffer
	return &CodeCollector{
		buf: buf,
		err: nil,
	}
}

func (tc *CodeCollector) Merge(tmpl string, args interface{}) *CodeCollector {
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

func (tc *CodeCollector) Export() ([]byte, error) {
	if tc.err != nil {
		return nil, tc.err
	}
	return tc.buf.Bytes(), nil
}
