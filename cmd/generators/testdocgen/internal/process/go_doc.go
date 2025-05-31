package process

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/ysuzuki19/robustruct/cmd/generators/testdocgen/internal/strchain"
)

func PlanGoDoc(source string, tds []TestDoc) ([]Plan, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	plans := []Plan{}

	for _, td := range tds {
		planed := false
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

				if *recvTypeName != *structName {
					continue
				}
			}

			if fn.Name.Name == td.FuncName {
				insertLine, replaceCount, err := FindExamplePosition(fset, fn.Doc.List)
				if err != nil {
					return nil, fmt.Errorf("failed to find example range: %w", err)
				}
				// fmt.Printf("fn(%s) example position: %d %d\n", fn.Name.Name, insertLine, replaceCount)
				lines := strchain.FromSlice([]string{"", "Example:"}).
					Extend(strchain.From(td.Content).Split("\n")).
					Map(func(line string) string {
						return "// " + line
					}).Collect()
				plans = append(plans, Plan{
					InsertIndex:  insertLine,
					ReplaceCount: replaceCount,
					Lines:        lines,
				})
				planed = true
				break
			}
		}
		if !planed {
			return plans, fmt.Errorf("no matching function found for %s.%s", td.StructName.UnwrapOrDefault(), td.FuncName)
		}
	}
	return plans, nil
}
