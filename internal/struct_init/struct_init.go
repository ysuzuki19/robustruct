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

func (si StructInit) IsIgnored(pattern string) bool {
	if si.pass.Fset == nil {
		return false
	}
	structPos := si.pass.Fset.Position(si.CompLit.Pos())
	for _, commentGroup := range si.AstFile.Comments {
		commentPos := si.pass.Fset.Position(commentGroup.End())
		if commentPos.Line+1 == structPos.Line {
			for _, comment := range commentGroup.List {
				if strings.Contains(comment.Text, "ignore:robustruct") ||
					strings.Contains(comment.Text, pattern) {
					return true
				}
			}
		}
	}
	return false
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
				CompLit:    *compLit,
				TypeStruct: *typeStruct,
				AstFile:    *file,
			})
			return true
		})
	}
	return
}
