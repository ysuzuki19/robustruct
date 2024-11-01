package fields_align

import (
	"flag"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/ysuzuki19/robustruct/internal/field_init"
	"github.com/ysuzuki19/robustruct/internal/struct_init"
)

var Analyzer = &analysis.Analyzer{
	Name: "fields_align",
	Doc:  "checks that all fields of a struct are sorted by defined order",
	URL:  "",
	Flags: flag.FlagSet{
		Usage: func() {},
	},
	Run:              run,
	RunDespiteErrors: false,
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	ResultType:       nil,
	FactTypes:        []analysis.Fact{},
}

func run(pass *analysis.Pass) (interface{}, error) {
	structInits := struct_init.List(*pass)
	for _, si := range structInits {
		// Fast path:
		// - if ignore comment exists
		// - if struct has no fields
		// - if struct definition is unnamed
		// - if struct has no visible fields
		if si.IsIgnored("ignore:fields_align") ||
			si.TypeStruct.NumFields() == 0 ||
			si.IsUnnamed() ||
			len(si.ListVisibleFields()) == 0 {
			{
				continue
			}
		}

		isAligned := true
		for i, elt := range si.CompLit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok || kv == nil {
				continue
			}
			field := si.TypeStruct.Field(i)
			if field == nil {
				continue
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok || key == nil {
				continue
			}
			if key.Name != field.Name() {
				isAligned = false
			}
		}
		if isAligned {
			continue
		}

		kves := make(map[string]*ast.KeyValueExpr)
		for _, elt := range si.CompLit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok || kv == nil {
				continue
			}
			ident, ok := kv.Key.(*ast.Ident)
			if !ok || ident == nil {
				continue
			}
			kves[ident.Name] = kv
		}

		alignedFields := field_init.NewFieldInits(pass, si.TypeStruct.NumFields())
		for i := 0; i < si.TypeStruct.NumFields(); i++ {
			field := si.TypeStruct.Field(i)
			if kve, ok := kves[field.Name()]; ok {
				alignedFields.PushKeyValueExpr(kve)
			}
		}

		newText, err := alignedFields.ToBytes()
		if err != nil {
			return nil, err
		}
		newText = append([]byte{'\n'}, newText...)

		pass.Report(analysis.Diagnostic{
			Pos:      si.CompLit.Pos(),
			End:      0,
			Category: "",
			Message:  "all fields of the struct must be sorted by defined order",
			URL:      "",
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
			Related: []analysis.RelatedInformation{},
		})
	}
	return nil, nil
}
