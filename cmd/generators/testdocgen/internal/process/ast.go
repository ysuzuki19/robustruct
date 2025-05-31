package process

import (
	"go/ast"

	"github.com/ysuzuki19/robustruct/pkg/option"
)

func RecvTypeName(fn *ast.FuncDecl) option.Option[string] {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return option.None[string]()
	}
	switch t := fn.Recv.List[0].Type.(type) {
	case *ast.Ident:
		return option.Some(&t.Name)
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return option.Some(&ident.Name)
		}
	case *ast.IndexExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return option.Some(&ident.Name)
		}
	case *ast.IndexListExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return option.Some(&ident.Name)
		}
	}
	return option.None[string]()
}
