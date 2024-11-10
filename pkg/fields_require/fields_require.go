package fields_require

import (
	"flag"
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/ysuzuki19/robustruct/internal/chain_writer"
	"github.com/ysuzuki19/robustruct/internal/field_inits"
	"github.com/ysuzuki19/robustruct/internal/struct_init"
	"github.com/ysuzuki19/robustruct/pkg/robustruct/settings"
)

var Analyzer = &analysis.Analyzer{
	Name:             settings.FeatureFieldsRequire.String(),
	Doc:              "checks that all fields of a struct are initialized in a composite literal",
	URL:              "",
	Flags:            flag.FlagSet{Usage: func() {}},
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
	// Fast path: all fields are initialized
	if si.TypeStruct.NumFields() == len(si.CompLit.Elts) ||
		si.IsIgnored("ignore:fields_require") {
		return nil
	}

	initializedFields := make(map[string]bool)
	for _, elt := range si.CompLit.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if ident, ok := kv.Key.(*ast.Ident); ok {
				initializedFields[ident.Name] = true
			} else {
				continue
			}
		}
	}

	missingFields := field_inits.New(pass, si.TypeStruct.NumFields()-len(si.CompLit.Elts))
	for _, field := range si.VisibleFields() {
		if initializedFields[field.Name()] {
			continue
		}
		missingFields.PushVarDefault(field, si.IsSamePackage())
	}

	if missingFields.Len() == 0 {
		return nil
	}

	cw := chain_writer.New(pass)
	for idx, field := range missingFields.List() {
		if idx != 0 {
			cw.Push(", ")
		}
		cw.Push(field.Key())
	}
	fieldsCSV, _ := cw.String()

	newText, err := missingFields.ToBytes()
	if err != nil {
		return err
	}
	if len(initializedFields) == 0 {
		newText = append([]byte{'\n'}, newText...)
	}

	pass.Report(analysis.Diagnostic{
		Pos:      si.CompLit.Pos(),
		End:      0,
		Category: "",
		Message:  fmt.Sprintf("fields '%s' are not initialized", fieldsCSV),
		URL:      "",
		SuggestedFixes: []analysis.SuggestedFix{{
			Message: "Add a missing fields",
			TextEdits: []analysis.TextEdit{{
				Pos:     si.CompLit.Rbrace,
				End:     si.CompLit.Rbrace,
				NewText: newText,
			}},
		}},
		Related: []analysis.RelatedInformation{},
	})
	return nil
}
