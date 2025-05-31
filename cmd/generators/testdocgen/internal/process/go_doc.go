package process

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func PlanGoDoc(source string, tds []TestDoc) ([]Plan, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	plans := []Plan{}

	for _, td := range tds {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			if structName, ok := td.StructName.Get(); ok {
				recvTypeName, ok := recvTypeName(fn).Get()
				if !ok {
					continue
				}

				fnName := fn.Name.Name
				if recvTypeName == structName && fnName == td.FuncName {
					begin, end, err := FindExampleRange(fset, fn.Doc.List)
					if err != nil {
						return nil, fmt.Errorf("failed to find example range: %w", err)
					}
					lines := strings.Split(td.Content, "\n")
					for i := range lines {
						lines[i] = "// " + lines[i]
					}
					plans = append(plans, Plan{
						Begin: begin,
						End:   end,
						Lines: lines,
					})
				}
			} else {
				if fn.Name.Name == td.FuncName {
					begin, end, err := FindExampleRange(fset, fn.Doc.List)
					if err != nil {
						return nil, fmt.Errorf("failed to find example range: %w", err)
					}
					lines := strings.Split(td.Content, "\n")
					for i := range lines {
						lines[i] = "// " + lines[i]
					}
					plans = append(plans, Plan{
						Begin: begin,
						End:   end,
						Lines: lines,
					})
				}
			}
		}
	}
	return plans, nil
}
