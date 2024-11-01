package struct_init

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type StructInit struct {
	pass       analysis.Pass
	AstFile    ast.File
	CompLit    ast.CompositeLit
	TypeStruct types.Struct
}

const disableAllPattern = "ignore:robustruct"

func (si StructInit) IsIgnored(pattern string) bool {
	if si.pass.Fset == nil {
		return false
	}
	structPos := si.pass.Fset.Position(si.CompLit.Pos())
	for _, commentGroup := range si.AstFile.Comments {
		commentPos := si.pass.Fset.Position(commentGroup.End())
		if commentPos.Line+1 == structPos.Line {
			for _, comment := range commentGroup.List {
				if strings.Contains(comment.Text, disableAllPattern) ||
					strings.Contains(comment.Text, pattern) {
					return true
				}
			}
		}
	}
	return false
}

func (si StructInit) IsUnnamed() bool {
	for _, elt := range si.CompLit.Elts {
		if _, ok := elt.(*ast.KeyValueExpr); !ok {
			return true
		}
	}
	return false
}

func (si StructInit) ListVisibleFields() (fields []*types.Var) {
	for i := 0; i < si.TypeStruct.NumFields(); i++ {
		field := si.TypeStruct.Field(i)
		if si.IsSamePackage() || field.Exported() {
			fields = append(fields, field)
		}
	}
	return
}

func (si StructInit) IsSamePackage() bool {
	if si.pass.Pkg == nil {
		return false
	}
	return si.pass.Pkg.Path() == si.TypeStruct.Field(0).Pkg().Path()
}

func List(pass analysis.Pass) (found []StructInit) {
	for _, file := range pass.Files {
		if file == nil {
			continue
		}
		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil {
				return false
			}
			compLit, ok := n.(*ast.CompositeLit)
			if !ok || compLit == nil {
				return true
			}

			if pass.TypesInfo == nil {
				return true
			}
			typ := pass.TypesInfo.TypeOf(compLit)
			if typ == nil {
				return true
			}

			ul := typ.Underlying()
			if ul == nil {
				return true
			}

			typeStruct, ok := ul.(*types.Struct)
			if !ok || typeStruct == nil {
				return true
			}

			found = append(found, StructInit{
				pass:       pass,
				AstFile:    *file,
				CompLit:    *compLit,
				TypeStruct: *typeStruct,
			})
			return true
		})
	}
	return
}
