package field_inits

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"

	"github.com/ysuzuki19/robustruct/internal/chain_writer"
	"github.com/ysuzuki19/robustruct/internal/field_init"
)

type FieldInits struct {
	pass *analysis.Pass
	list []*field_init.FieldInit
}

func New(pass *analysis.Pass, cap int) FieldInits {
	return FieldInits{
		pass: pass,
		list: make([]*field_init.FieldInit, 0, cap),
	}
}

func (fis *FieldInits) List() []*field_init.FieldInit {
	return fis.list
}

func (fis *FieldInits) Len() int {
	return len(fis.list)
}

func (fis *FieldInits) ToBytes() ([]byte, error) {
	cw := chain_writer.New(fis.pass)
	for _, fi := range fis.list {
		cw.Push(fi.KeyValueExpr()).Push(",\n")
	}
	return cw.Bytes()
}

func (fis *FieldInits) Push(i *field_init.FieldInit) {
	fis.list = append(fis.list, i)
}

func (fis *FieldInits) PushKeyValueExpr(kve *ast.KeyValueExpr) {
	fis.Push(field_init.FromKeyValueExpr(kve))
}

func (fis *FieldInits) PushVarDefault(field *types.Var, isSamePackage bool) {
	fis.Push(field_init.FromVarDefault(field, isSamePackage))
}
