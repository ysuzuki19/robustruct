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
	structPos := si.pass.Fset.Position(si.CompLit.Pos()).Line
	structEnd := si.pass.Fset.Position(si.CompLit.End()).Line
	for _, commentGroup := range si.AstFile.Comments {
		commentLine := si.pass.Fset.Position(commentGroup.End()).Line
		if commentLine+1 == structPos || commentLine == structEnd {
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

func (si StructInit) VisibleFields() (fields []*types.Var) {
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

type InspectInput struct {
	Pass        *analysis.Pass
	DisableTest bool
	Handler     func(passs *analysis.Pass, si StructInit) error
}

func Inspect(input InspectInput) error {
	for _, file := range input.Pass.Files {
		if file == nil {
			continue
		}
		if input.DisableTest && strings.HasSuffix(input.Pass.Fset.File(file.Pos()).Name(), "_test.go") {
			continue
		}
		var err error
		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil || err != nil {
				return false
			}
			compLit, ok := n.(*ast.CompositeLit)
			if !ok || compLit == nil {
				return true
			}

			if input.Pass.TypesInfo == nil {
				return true
			}
			typ := input.Pass.TypesInfo.TypeOf(compLit)
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

			err = input.Handler(input.Pass, StructInit{
				pass:       *input.Pass,
				AstFile:    *file,
				CompLit:    *compLit,
				TypeStruct: *typeStruct,
			})
			return true
		})
		if err != nil {
			return err
		}
	}
	return nil
}
