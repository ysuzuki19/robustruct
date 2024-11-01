package fields_require

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

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
	structInits := struct_init.List(*pass)
	for _, si := range structInits {
		// Fast path: all fields are initialized
		if si.TypeStruct.NumFields() == len(si.CompLit.Elts) ||
			si.IsIgnored("ignore:fields_require") {
			continue
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

		var missingFields []*ast.KeyValueExpr
		for _, field := range si.ListVisibleFields() {
			if initializedFields[field.Name()] {
				continue
			}
			missingFields = append(missingFields, &ast.KeyValueExpr{
				Key:   ast.NewIdent(field.Name()),
				Colon: 0,
				Value: generateDefaultExpr(field.Type(), si.IsSamePackage()),
			})
		}

		if len(missingFields) > 0 {
			var fieldsCSV bytes.Buffer
			_ = format.Node(&fieldsCSV, pass.Fset, missingFields[0].Key)
			for _, field := range missingFields[1:] {
				fieldsCSV.WriteString(", ")
				_ = format.Node(&fieldsCSV, pass.Fset, field.Key)
			}

			// if all fields are missing, add a newline before the first field
			var buf bytes.Buffer
			if len(missingFields) == si.TypeStruct.NumFields() {
				buf.WriteString("\n")
			}
			for _, field := range missingFields {
				if err := format.Node(&buf, pass.Fset, field); err != nil {
					continue
				}
				buf.WriteString(",\n")
			}
			newText := buf.Bytes()

			pass.Report(analysis.Diagnostic{
				Pos:      si.CompLit.Pos(),
				End:      0,
				Category: "",
				Message:  fmt.Sprintf("fields '%s' are not initialized", fieldsCSV.String()),
				URL:      "",
				SuggestedFixes: []analysis.SuggestedFix{{
					Message:   "Add a missing fields",
					TextEdits: []analysis.TextEdit{{Pos: si.CompLit.Rbrace, End: si.CompLit.Rbrace, NewText: newText}},
				}},
				Related: []analysis.RelatedInformation{},
			})
		}
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
