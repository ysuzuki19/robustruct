package fields_require

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/ysuzuki19/robustruct/internal/chain_writer"
	"github.com/ysuzuki19/robustruct/internal/field_init"
	"github.com/ysuzuki19/robustruct/internal/struct_init"
)

var Analyzer = &analysis.Analyzer{
	Name:             "fields_require",
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
	handler := handlerFactory(pass)
	if err := struct_init.List(*pass).ForEach(handler); err != nil {
		return nil, err
	}
	return nil, nil
}

func generateDefaultExpr(typ types.Type, isSamePackage bool) ast.Expr {
	switch t := typ.Underlying().(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
			types.Float32, types.Float64, types.Complex64, types.Complex128:
			return ast.NewIdent("0")
		case types.String:
			return &ast.BasicLit{
				ValuePos: 0,
				Kind:     token.STRING,
				Value:    `""`,
			}
		case types.Bool:
			return ast.NewIdent("false")
		}
	}
	if isSamePackage {
		if named, ok := typ.(*types.Named); ok {
			// return without package
			return ast.NewIdent(named.Obj().Name() + "{}")
		}
	} else {
		// return with package
		return ast.NewIdent(typ.String() + "{}")
	}
	return ast.NewIdent("nil")
}

func handlerFactory(pass *analysis.Pass) func(si struct_init.StructInit) error {
	return func(si struct_init.StructInit) error {
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

		missingFields := field_init.NewFieldInits(pass, si.TypeStruct.NumFields()-len(si.CompLit.Elts))
		for _, field := range si.ListVisibleFields() {
			if initializedFields[field.Name()] {
				continue
			}
			missingFields.PushExpr(field.Name(), generateDefaultExpr(field.Type(), si.IsSamePackage()))
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
}
