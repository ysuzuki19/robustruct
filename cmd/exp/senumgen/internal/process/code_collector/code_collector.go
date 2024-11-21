package tmpl

import (
	"bytes"
	"fmt"
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

func (tc *CodeCollector) LF() *CodeCollector {
	if tc.err != nil {
		return tc
	}
	tc.buf.WriteString("\n")
	return tc
}

func (tc *CodeCollector) Str(content string) *CodeCollector {
	if tc.err != nil {
		return tc
	}
	tc.buf.WriteString(content)
	return tc
}

func (tc *CodeCollector) Format(format string, a ...any) *CodeCollector {
	content := fmt.Sprintf(format, a...)
	tc.buf.WriteString(content)
	return tc
}

func (tc *CodeCollector) Func(f func(tc *CodeCollector) *CodeCollector) *CodeCollector {
	if tc.err != nil {
		return tc
	}
	return f(tc)
}

func (tc *CodeCollector) Tmpl(tmpl string, args interface{}) *CodeCollector {
	if tc.err != nil {
		return tc
	}

	t, err := template.New("tmpl").Funcs(
		template.FuncMap{
			"capitalize": Capitalize,
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
