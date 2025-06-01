package process_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ysuzuki19/robustruct/cmd/generators/testdocgen/internal/process"
)

func TestApplyGoDoc(t *testing.T) {
	require := require.New(t)
	source := strings.Join([]string{"0", "1", "2", "3"}, "\n")

	plans := []process.Plan{
		{
			InsertIndex:  0,
			ReplaceCount: 0,
			Lines:        []string{"a", "b"},
		},
	}
	updated := process.ApplyGoDoc(source, plans)
	require.Equal("a\nb\n0\n1\n2\n3", updated)

	plans = []process.Plan{
		{
			InsertIndex:  0,
			ReplaceCount: 0,
			Lines:        []string{"a", "b"},
		},
		{
			InsertIndex:  1,
			ReplaceCount: 0,
			Lines:        []string{"a", "b"},
		},
	}
	updated = process.ApplyGoDoc(source, plans)
	require.Equal("a\nb\n0\na\nb\n1\n2\n3", updated)

	plans = []process.Plan{
		{
			InsertIndex:  0,
			ReplaceCount: 2,
			Lines:        []string{"a", "b"},
		},
	}
	updated = process.ApplyGoDoc(source, plans)
	require.Equal("a\nb\n2\n3", updated)

}
