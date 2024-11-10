package field_init

import (
	"go/ast"
	"go/token"
	"go/types"
)

type FieldInit struct {
	kve *ast.KeyValueExpr
}

func (fi *FieldInit) KeyValueExpr() *ast.KeyValueExpr {
	return fi.kve
}

func (fi *FieldInit) Key() ast.Expr {
	return fi.kve.Key
}

func FromKeyValueExpr(kve *ast.KeyValueExpr) *FieldInit {
	return &FieldInit{kve}
}

func FromVarDefault(field *types.Var, isSamePackage bool) *FieldInit {
	return &FieldInit{
		&ast.KeyValueExpr{
			Key:   ast.NewIdent(field.Name()),
			Colon: 0,
			Value: generateDefaultExpr(field.Type(), isSamePackage),
		},
	}
}

func generateDefaultExpr(typ types.Type, isSamePackage bool) ast.Expr {
	switch t := typ.Underlying().(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
			types.Float32, types.Float64, types.Complex64, types.Complex128:
			return ast.NewIdent("0")
		case types.String:
			return &ast.BasicLit{
				ValuePos: 0,
				Kind:     token.STRING,
				Value:    `""`,
			}
		case types.Bool:
			return ast.NewIdent("false")
		}
	}
	if isSamePackage {
		if named, ok := typ.(*types.Named); ok {
			// return without package
			return ast.NewIdent(named.Obj().Name() + "{}")
		}
	} else {
		// return with package
		return ast.NewIdent(typ.String() + "{}")
	}
	return ast.NewIdent("nil")
}
