package struct_init

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type StructInit struct {
	pass       *analysis.Pass
	AstFile    *ast.File
	CompLit    *ast.CompositeLit
	TypeStruct types.Struct
}

func (si StructInit) IsIgnored(pattern string) bool {
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

func List(pass *analysis.Pass) (found []StructInit) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil {
				return false
			}
			compLit, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}

			typ := pass.TypesInfo.TypeOf(compLit)
			if typ == nil {
				return true
			}

			typeStruct, ok := pass.TypesInfo.TypeOf(compLit).Underlying().(*types.Struct)
			if !ok {
				return true
			}

			found = append(found, StructInit{
				pass:       pass,
				CompLit:    compLit,
				TypeStruct: *typeStruct,
				AstFile:    file,
			})
			return true
		})
	}
	return
}
