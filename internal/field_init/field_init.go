package field_init

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/ysuzuki19/robustruct/internal/chain_writer"
)

type FieldInit struct {
	kve *ast.KeyValueExpr
}

func (fi *FieldInit) Key() ast.Expr {
	return fi.kve.Key
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
	cw := chain_writer.New(fis.pass)
	for _, fi := range fis.list {
		cw.Push(fi.kve).Push(",\n")
	}
	return cw.Bytes()
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
