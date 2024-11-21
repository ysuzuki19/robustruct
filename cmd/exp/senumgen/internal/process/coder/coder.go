package coder

import (
	"bytes"
	"fmt"
	"text/template"
)

type Coder struct {
	buf     bytes.Buffer
	globals map[string]interface{}
	err     error
}

func New() *Coder {
	var buf bytes.Buffer
	return &Coder{
		buf:     buf,
		globals: make(map[string]interface{}),
		err:     nil,
	}
}

func (tc *Coder) Globals(globals map[string]interface{}) *Coder {
	if tc.err != nil {
		return tc
	}
	tc.globals = globals
	return tc
}

func (tc *Coder) LF() *Coder {
	if tc.err != nil {
		return tc
	}
	tc.buf.WriteString("\n")
	return tc
}

func (tc *Coder) Str(content string) *Coder {
	if tc.err != nil {
		return tc
	}
	tc.buf.WriteString(content)
	return tc
}

func (tc *Coder) Format(format string, a ...any) *Coder {
	content := fmt.Sprintf(format, a...)
	tc.buf.WriteString(content)
	return tc
}

func (tc *Coder) Func(f func(tc *Coder) *Coder) *Coder {
	if tc.err != nil {
		return tc
	}
	return f(tc)
}

func (tc *Coder) Block(start, end string, f func(tc *Coder) *Coder) *Coder {
	if tc.err != nil {
		return tc
	}
	tc.buf.WriteString(start)
	f(tc)
	tc.buf.WriteString(end)
	return tc
}

func (tc *Coder) Tmpl(tmpl string, args map[string]interface{}) *Coder {
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

func (tc *Coder) Export() ([]byte, error) {
	if tc.err != nil {
		return nil, tc.err
	}
	return tc.buf.Bytes(), nil
}
