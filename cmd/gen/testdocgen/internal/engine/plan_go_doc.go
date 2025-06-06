package engine

import (
	"fmt"
	"go/parser"
	"go/token"

	"github.com/ysuzuki19/robustruct/cmd/gen/testdocgen/internal/engine/astutil"
	"github.com/ysuzuki19/robustruct/cmd/gen/testdocgen/internal/strchain"
)

type Plan struct {
	InsertIndex  int
	ReplaceCount int
	Lines        []string
}

func PlanGoDoc(source string, tds []TestDoc) ([]Plan, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	plans := []Plan{}

	for _, td := range tds {
		planed := false
		for _, fn := range astutil.ListFnDecls(fset, file) {
			if fn.Recv.UnwrapOrDefault() != td.StructName.UnwrapOrDefault() {
				// None and Some => unmatched
				// Some and None => unmatched
				// Some and Some => check equality
				continue
			}

			if fn.Name == td.FuncName {
				insertLine, replaceCount, err := fn.ExamplePosition()
				if err != nil {
					return nil, fmt.Errorf("failed to find example range: %w", err)
				}
				// fmt.Printf("fn(%s) example position: %d %d\n", fn.Name.Name, insertLine, replaceCount)
				lines := strchain.FromSlice([]string{"", "Example:"}).
					Extend(strchain.From(td.Content).Split("\n")).
					Map(func(line string) string {
						return "//" + line
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
