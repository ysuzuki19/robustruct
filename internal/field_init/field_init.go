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

func (fi *FieldInit) Key() ast.Expr {
	return fi.kve.Key
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

func (fis *FieldInits) List() []*FieldInit {
	return fis.list
}

func (fis *FieldInits) Len() int {
	return len(fis.list)
}

func (fis *FieldInits) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
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

func (fis *FieldInits) PushKeyValueExpr(kve *ast.KeyValueExpr) {
	fis.Push(&FieldInit{kve})
}

func (fis *FieldInits) PushExpr(key string, value ast.Expr) {
	fis.Push(&FieldInit{
		&ast.KeyValueExpr{
			Key:   ast.NewIdent(key),
			Colon: 0,
			Value: value,
		},
	})
}
