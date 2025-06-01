package process

import (
	"go/ast"

	"github.com/ysuzuki19/robustruct/pkg/option"
)

func ListFnDecls(file *ast.File) []*ast.FuncDecl {
	decls := []*ast.FuncDecl{}
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			decls = append(decls, fn)
		}
	}
	return decls
}

func recvTypeName(fn *ast.FuncDecl) option.Option[string] {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return option.None[string]()
	}
	switch t := fn.Recv.List[0].Type.(type) {
	case *ast.Ident:
		return option.Some(&t.Name)
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return option.Some(&ident.Name)
		} else {
			if indexExpr, ok := t.X.(*ast.IndexExpr); ok {
				if ident, ok := indexExpr.X.(*ast.Ident); ok {
					return option.Some(&ident.Name)
				}
			} else if indexListExpr, ok := t.X.(*ast.IndexListExpr); ok {
				if ident, ok := indexListExpr.X.(*ast.Ident); ok {
					return option.Some(&ident.Name)
				}
			}
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
