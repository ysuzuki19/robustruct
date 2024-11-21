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

func (c *Coder) Str(content string) *Coder {
	if c.err != nil {
		return c
	}
	c.buf.WriteString(content)
	return c
}

func (c *Coder) LF() *Coder {
	return c.Str("\n")
}

func (c *Coder) Space() *Coder {
	return c.Str(" ")
}

func (c *Coder) Capitalize(s string) *Coder {
	return c.Str(Capitalize(s))
}

func (c *Coder) Format(format string, a ...any) *Coder {
	content := fmt.Sprintf(format, a...)
	return c.Str(content)
}

func (c *Coder) Func(f func(*Coder)) *Coder {
	if c.err != nil {
		return c
	}
	f(c)
	return c
}

func (c *Coder) Wrap(start, end string, f func(*Coder)) *Coder {
	if c.err != nil {
		return c
	}
	return c.Str(start).Func(f).Str(end)
}

func (c *Coder) Parens(f func(*Coder)) *Coder {
	return c.Wrap("(", ")", f)
}

func (c *Coder) Braces(f func(*Coder)) *Coder {
	return c.Wrap("{", "}", f)
}

func (c *Coder) Block(f func(*Coder)) *Coder {
	return c.Braces(func(c *Coder) {
		c.LF()
		f(c)
		c.LF()
	})
}

func (c *Coder) Brackets(f func(*Coder)) *Coder {
	return c.Wrap("[", "]", f)
}

func (c *Coder) Tmpl(tmpl string, varss ...Vars) *Coder {
	if c.err != nil {
		return c
	}

	locals := make(Vars)
	for k, v := range c.globals {
		locals[k] = v
	}
	for _, vars := range varss {
		for k, v := range vars {
			locals[k] = v
		}
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
		return c.buf.Bytes(), c.err
	}
	return c.buf.Bytes(), nil
}
