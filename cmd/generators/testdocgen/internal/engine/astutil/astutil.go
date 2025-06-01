package astutil

import (
	"go/ast"
	"go/token"
	"regexp"

	"github.com/ysuzuki19/robustruct/pkg/option"
)

type Function struct {
	Name string
	Recv option.Option[string]
	decl *ast.FuncDecl
	fset *token.FileSet
}

var docExampleRegex = regexp.MustCompile(`^\s*//\s*Example:?.*`)

func (f Function) ExamplePosition() (begin int, count int, err error) {
	if f.decl.Doc == nil || len(f.decl.Doc.List) == 0 {
		return f.fset.Position(f.decl.Pos()).Line - 1, 0, nil
	}
	found := option.None[int]()
	var searched int
	for _, comment := range f.decl.Doc.List {
		searched = f.fset.Position(comment.Pos()).Line
		if found.IsSome() {
			// if exampleAnnotation.IsSome() {
			// if comment.Text == "//" {
			// 	if begin, ok := begin.Get(); ok {
			// 		return *begin, searched, nil
			// 	}
			// } else {
			// 	if begin.IsNone() {
			// 		pos := fset.Position(comment.Pos())
			// 		begin = option.NewSome(pos.Line)
			// 	}
			// }
		} else if docExampleRegex.MatchString(comment.Text) {
			found = option.NewSome(searched - 1)
		}
	}
	if exm, ok := found.Get(); ok {
		return *exm, searched - *exm, nil
	}
	// if `Example:` is not found, we return the last searched line as the end
	return searched, 0, nil
}

func ListFnDecls(fset *token.FileSet, file *ast.File) []Function {
	funcs := []Function{}
	// decls := []*ast.FuncDecl{}
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			funcs = append(funcs, Function{
				Name: fn.Name.Name,
				Recv: RecvTypeName(fn),
				decl: fn,
				fset: fset,
			})
		}
	}
	return funcs
}

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
