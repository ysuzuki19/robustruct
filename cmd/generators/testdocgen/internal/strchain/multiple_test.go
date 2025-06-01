package strchain_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ysuzuki19/robustruct/cmd/generators/testdocgen/internal/strchain"
)

func TestFromSlice(t *testing.T) {
	require := require.New(t)
	// testdoc begin FromSlice
	m := strchain.FromSlice([]string{"a", "b", "c"})
	require.Equal([]string{"a", "b", "c"}, m.Collect())
	// testdoc end
}

func TestSlice(t *testing.T) {
	require := require.New(t)
	// testdoc begin multiple.Slice
	input := []string{"a", "b", "c", "d", "e"}
	require.Equal([]string{"a", "b"}, strchain.FromSlice(input).Slice(0, 2).Collect())
	require.Equal([]string{"d", "e"}, strchain.FromSlice(input).Slice(3, -1).Collect())
	// testdoc end
}

func TestSplice(t *testing.T) {
	require := require.New(t)
	// testdoc begin multiple.Splice
	m := strchain.FromSlice([]string{"a", "b", "c", "d"})
	m = m.Splice(1, 2, []string{"x", "y"})
	require.Equal([]string{"a", "x", "y", "d"}, m.Collect())
	// testdoc end
}

func Map(t *testing.T) {
	require := require.New(t)
	// testdoc begin multiple.Map
	data := strchain.
		FromSlice([]string{"a", "b", "c"}).
		Map(func(s string) string {
			return s + "!"
		}).
		Collect()
	require.Equal([]string{"a!", "b!", "c!"}, data)
	// testdoc end
}

func TestAppend(t *testing.T) {
	require := require.New(t)
	// testdoc begin multiple.Append
	m := strchain.
		FromSlice([]string{"a", "b"}).
		Append("c", "d")
	require.Equal([]string{"a", "b", "c", "d"}, m.Collect())
	// testdoc end
}

func TestExtend(t *testing.T) {
	require := require.New(t)
	// testdoc begin multiple.Extend
	m1 := strchain.FromSlice([]string{"a", "b"})
	m2 := strchain.FromSlice([]string{"c", "d"})
	m := m1.Extend(m2)
	require.Equal([]string{"a", "b", "c", "d"}, m.Collect())
	// testdoc end
}

func TestJoin(t *testing.T) {
	require := require.New(t)
	// testdoc begin multiple.Join
	m := strchain.FromSlice([]string{"a", "b", "c"})
	s := m.Join(", ")
	require.Equal("a, b, c", s.String())
	// testdoc end
}

func TestCollect(t *testing.T) {
	require := require.New(t)
	// testdoc begin multiple.Collect
	m := strchain.FromSlice([]string{"a", "b", "c"})
	data := m.Collect()
	require.Equal([]string{"a", "b", "c"}, data)
	// testdoc end
}

func TestEntries(t *testing.T) {
	require := require.New(t)
	// testdoc begin multiple.Entries
	entries := strchain.FromSlice([]string{"a", "b", "c"}).Entries()
	require.Equal(strchain.From("a"), entries[0])
	require.Equal(strchain.From("b"), entries[1])
	require.Equal(strchain.From("c"), entries[2])
	// testdoc end
}
