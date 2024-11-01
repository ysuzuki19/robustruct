package field_init

import (
	"bytes"
	"go/ast"
	"go/format"

	"golang.org/x/tools/go/analysis"
)

type FieldInit struct {
	kve *ast.KeyValueExpr
}

func NewFieldInit(kve *ast.KeyValueExpr) *FieldInit {
	return &FieldInit{
		kve: kve,
	}
}

func (fi *FieldInit) Write(pass *analysis.Pass, buf *bytes.Buffer) error {
	return format.Node(buf, pass.Fset, fi.kve)
}

type FieldInits struct {
	pass *analysis.Pass
	list []*FieldInit
}

func NewFieldInits(pass *analysis.Pass, cap int) FieldInits {
	return FieldInits{
		pass: pass,
		list: make([]*FieldInit, 0, cap),
	}
}

func (fis *FieldInits) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("\n")
	for _, fi := range fis.list {
		err := fi.Write(fis.pass, &buf)
		if err != nil {
			continue
		}
		buf.WriteString(",\n")
	}
	return buf.Bytes(), nil
}

func (fis *FieldInits) Push(i *FieldInit) {
	fis.list = append(fis.list, i)
}
