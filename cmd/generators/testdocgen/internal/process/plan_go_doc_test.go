package process_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ysuzuki19/robustruct/cmd/generators/testdocgen/internal/process"
	"github.com/ysuzuki19/robustruct/pkg/option"
)

func TestPlanGoDo(t *testing.T) {
	require := require.New(t)

	tds := []process.TestDoc{
		{
			StructName: option.None[string](),
			FuncName:   "Utility",
			Content:    "example for Utility",
		},
		{
			StructName: option.NewSome("Sample"),
			FuncName:   "Method",
			Content:    "example for Sample.Method",
		},
	}
	source := strings.Join([]string{
		"package sample",
		"func Utility() {}",
		"func Dummy() {}",
		"func (s Sample) Method() {}",
	}, "\n")

	plans, err := process.PlanGoDoc(source, tds)
	require.NoError(err)
	require.Len(plans, 2)

	require.Equal(process.Plan{
		InsertIndex:  1,
		ReplaceCount: 0,
		Lines: []string{
			"//",
			"//Example:",
			"//example for Utility",
		},
	}, plans[0])

	require.Equal(process.Plan{
		InsertIndex:  3,
		ReplaceCount: 0,
		Lines: []string{
			"//",
			"//Example:",
			"//example for Sample.Method",
		},
	}, plans[1])
}
