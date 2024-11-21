package coder

import (
	"bytes"
	"fmt"
	"text/template"
)

type Vars map[string]interface{}

type Coder struct {
	buf     bytes.Buffer
	globals Vars
	err     error
}

func New() *Coder {
	var buf bytes.Buffer
	return &Coder{
		buf:     buf,
		globals: make(Vars),
		err:     nil,
	}
}

func (c *Coder) Globals(globals Vars) *Coder {
	if c.err != nil {
		return c
	}
	c.globals = globals
	return c
}

func (c *Coder) LF() *Coder {
	if c.err != nil {
		return c
	}
	c.buf.WriteString("\n")
	return c
}

func (c *Coder) Str(content string) *Coder {
	if c.err != nil {
		return c
	}
	c.buf.WriteString(content)
	return c
}

func (c *Coder) Space() *Coder {
	return c.Str(" ")
}

func (c *Coder) Capitalize(s string) *Coder {
	return c.Str(Capitalize(s))
}

func (c *Coder) Format(format string, a ...any) *Coder {
	content := fmt.Sprintf(format, a...)
	c.buf.WriteString(content)
	return c
}

func (c *Coder) Func(f func(c *Coder) *Coder) *Coder {
	if c.err != nil {
		return c
	}
	return f(c)
}

func (c *Coder) Block(start, end string, f func(c *Coder) *Coder) *Coder {
	if c.err != nil {
		return c
	}
	c.buf.WriteString(start)
	f(c)
	c.buf.WriteString(end)
	return c
}

func (c *Coder) Tmpl(tmpl string, args Vars) *Coder {
	if c.err != nil {
		return c
	}

	locals := make(Vars)
	for k, v := range c.globals {
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
		c.err = err
	}

	if err := t.Execute(&c.buf, locals); err != nil {
		c.err = err
	}

	return c
}

func (c *Coder) Export() ([]byte, error) {
	if c.err != nil {
		return nil, c.err
	}
	return c.buf.Bytes(), nil
}
