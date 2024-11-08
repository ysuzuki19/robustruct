package chain_writer

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"

	"golang.org/x/tools/go/analysis"
)

type ChainWriter struct {
	pass *analysis.Pass
	buf  bytes.Buffer
	err  error
}

func New(pass *analysis.Pass) *ChainWriter {
	var buf bytes.Buffer
	var err error
	return &ChainWriter{pass, buf, err}
}

func (cw *ChainWriter) Push(data any) *ChainWriter {
	if cw.err != nil {
		return cw
	}
	var err error
	switch v := data.(type) {
	case string:
		_, err = cw.buf.WriteString(v)
	case []byte:
		_, err = cw.buf.Write(v)
	case *ast.Ident, *ast.KeyValueExpr:
		err = format.Node(&cw.buf, cw.pass.Fset, v)
	default:
		err = fmt.Errorf("unsupported type %T", v)
	}
	if err != nil {
		cw.err = err
	}
	return cw
}

func (cw *ChainWriter) Bytes() ([]byte, error) {
	return cw.buf.Bytes(), cw.err
}

func (cw *ChainWriter) String() (string, error) {
	return cw.buf.String(), cw.err
}
