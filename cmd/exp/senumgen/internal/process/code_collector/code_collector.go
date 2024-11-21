package tmpl

import (
	"bytes"
	"fmt"
	"text/template"
)

type CodeCollector struct {
	buf     bytes.Buffer
	globals map[string]interface{}
	err     error
}

func New() *CodeCollector {
	var buf bytes.Buffer
	return &CodeCollector{
		buf:     buf,
		globals: make(map[string]interface{}),
		err:     nil,
	}
}

func (tc *CodeCollector) Globals(globals map[string]interface{}) *CodeCollector {
	if tc.err != nil {
		return tc
	}
	tc.globals = globals
	return tc
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

func (tc *CodeCollector) Tmpl(tmpl string, args map[string]interface{}) *CodeCollector {
	if tc.err != nil {
		return tc
	}

	locals := make(map[string]interface{})
	for k, v := range tc.globals {
		locals[k] = v
	}
	for k, v := range args {
		locals[k] = v
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

	if err := t.Execute(&tc.buf, locals); err != nil {
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
