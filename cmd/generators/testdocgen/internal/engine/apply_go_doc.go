package engine

import (
	"sort"

	"github.com/ysuzuki19/robustruct/cmd/generators/testdocgen/internal/strchain"
)

func ApplyGoDoc(source string, plans []Plan) string {
	sort.Slice(plans, func(i, j int) bool {
		return plans[i].InsertIndex > plans[j].InsertIndex
	})

	lines := strchain.From(source).Split("\n")
	for _, plan := range plans {
		lines = lines.Splice(plan.InsertIndex, plan.ReplaceCount, plan.Lines)
	}

	return lines.Join("\n").String()
}
