package fields_align

import (
	"bytes"
	"go/ast"
	"go/format"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/ysuzuki19/robustruct/internal/struct_init"
)

var Analyzer = &analysis.Analyzer{
	Name: "fields_align",
	Doc:  "checks that all fields of a struct are sorted by defined order",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	structInits := struct_init.List(pass)
	for _, si := range structInits {
		// Fast path:
		// - if the struct is ignored by comment
		// - if the number of fields of the struct is not equal to the number of initialized fields
		// - if the number of fields of the struct is 0
		if si.IsIgnored("ignore:fields_align") ||
			// len(si.CompLit.Elts) != si.TypeStruct.NumFields() ||
			si.TypeStruct.NumFields() == 0 {
			continue
		}

		isUnnamed := false
		isAligned := true
		for i, elt := range si.CompLit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				isUnnamed = true
			}
			field := si.TypeStruct.Field(i)
			if kv.Key.(*ast.Ident).Name != field.Name() {
				isAligned = false
			}
		}
		if isUnnamed {
			continue
		}
		if isAligned {
			continue
		}

		fieldInits := make(map[string]*ast.KeyValueExpr)
		for _, elt := range si.CompLit.Elts {
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				if ident, ok := kv.Key.(*ast.Ident); ok {
					fieldInits[ident.Name] = kv
				} else {
					continue
				}
			}
		}

		var buf bytes.Buffer
		buf.WriteString("\n")
		for i := 0; i < si.TypeStruct.NumFields(); i++ {
			field := si.TypeStruct.Field(i)
			if kv, ok := fieldInits[field.Name()]; ok {
				buf.WriteString(field.Name() + ": ")
				if err := format.Node(&buf, pass.Fset, kv.Value); err != nil {
					panic(err)
				}
				buf.WriteString(",\n")
			}
		}
		newText := buf.Bytes()

		pass.Report(analysis.Diagnostic{
			Pos:     si.CompLit.Pos(),
			Message: "all fields of the struct must be sorted by defined order",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "Align fields by defined order",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     si.CompLit.Lbrace + 1,
							End:     si.CompLit.Rbrace,
							NewText: newText,
						},
					},
				},
			},
		})
	}
	return nil, nil
}
