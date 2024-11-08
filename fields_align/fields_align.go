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
	if err := struct_init.Inspect(pass, handler); err != nil {
		return nil, err
	}
	return nil, nil
}

func handler(pass *analysis.Pass, si struct_init.StructInit) error {
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
			return nil
		}
	}

	definedOrder := map[string]int{}
	for i, field := range si.ListVisibleFields() {
		definedOrder[field.Name()] = i
	}

	{
		isAligned := true
		cursor := -1
		for _, elt := range si.CompLit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok || kv == nil {
				return nil
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok || key == nil {
				return nil
			}
			keyCursor, ok := definedOrder[key.Name]
			if !ok {
				return nil
			}
			if cursor < keyCursor {
				cursor = keyCursor
			} else {
				isAligned = false
				break
			}
		}
		if isAligned {
			return nil
		}
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
		return err
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
	return nil
}
