package strchain

import "strings"

type multiple struct {
	ss []string
}

// FromSlice returns a string chaining object from a `[]string`.
//
// Example:
//
//	m := strchain.FromSlice([]string{"a", "b", "c"})
//	require.Equal([]string{"a", "b", "c"}, m.Collect())
func FromSlice(ss []string) multiple {
	return multiple{ss}
}

// Collect returns the `[]string` from the chaining object.
//
// Example:
//
//	m := strchain.FromSlice([]string{"a", "b", "c"})
//	data := m.Collect()
//	require.Equal([]string{"a", "b", "c"}, data)
func (m multiple) Collect() []string {
	return m.ss
}

// Slice returns a new string chaining object with a slice of the original strings.
//
// Example:
//
//	input := []string{"a", "b", "c", "d", "e"}
//	require.Equal([]string{"a", "b"}, strchain.FromSlice(input).Slice(0, 2).Collect())
//	require.Equal([]string{"d", "e"}, strchain.FromSlice(input).Slice(3, -1).Collect())
func (m multiple) Slice(start, end int) multiple {
	if end < 0 {
		end = len(m.ss)
	}
	if start < 0 || end > len(m.ss) || start > end {
		return FromSlice([]string{})
	}
	return FromSlice(m.ss[start:end])
}

// Splice modifies the string chaining object by removing a slice of strings
//
// Example:
//
//	m := strchain.FromSlice([]string{"a", "b", "c", "d"})
//	m = m.Splice(1, 2, []string{"x", "y"})
//	require.Equal([]string{"a", "x", "y", "d"}, m.Collect())
func (m multiple) Splice(start, count int, ss []string) multiple {
	if start < 0 || start > len(m.ss) {
		return m
	}
	if count < 0 || start+count > len(m.ss) {
		count = len(m.ss) - start
	}
	m.ss = append(m.ss[:start], append(ss, m.ss[start+count:]...)...)
	return m
}

// Map applies a function to each string in the chaining object and returns a new chaining object.
//
// Example:
//
//	data := strchain.
//		FromSlice([]string{"a", "b", "c"}).
//		Map(func(s string) string {
//			return s + "!"
//		}).
//		Collect()
//	require.Equal([]string{"a!", "b!", "c!"}, data)
func (m multiple) Map(fn func(string) string) multiple {
	for i := range m.ss {
		m.ss[i] = fn(m.ss[i])
	}
	return m
}

// Append adds one or more strings to the end of the chaining object.
//
// Example:
//
//	m := strchain.
//		FromSlice([]string{"a", "b"}).
//		Append("c", "d")
//	require.Equal([]string{"a", "b", "c", "d"}, m.Collect())
func (m multiple) Append(s ...string) multiple {
	m.ss = append(m.ss, s...)
	return m
}

// Extend appends another chaining object to the current object.
//
// Example:
//
//	m1 := strchain.FromSlice([]string{"a", "b"})
//	m2 := strchain.FromSlice([]string{"c", "d"})
//	m := m1.Extend(m2)
//	require.Equal([]string{"a", "b", "c", "d"}, m.Collect())
func (m multiple) Extend(other multiple) multiple {
	m.ss = append(m.ss, other.ss...)
	return m
}

// Join concatenates the strings in the chaining object with a separator.
//
// Example:
//
//	m := strchain.FromSlice([]string{"a", "b", "c"})
//	s := m.Join(", ")
//	require.Equal("a, b, c", s.String())
func (m multiple) Join(sep string) single {
	return From(strings.Join(m.ss, sep))
}

// Entries returns a slice of single objects representing each string in the chaining object.
//
// Example:
//
//	entries := strchain.FromSlice([]string{"a", "b", "c"}).Entries()
//	require.Equal(strchain.From("a"), entries[0])
//	require.Equal(strchain.From("b"), entries[1])
//	require.Equal(strchain.From("c"), entries[2])
func (m multiple) Entries() []single {
	entries := make([]single, len(m.ss))
	for i, s := range m.ss {
		entries[i] = From(s)
	}
	return entries
}

//go:generate go run ../../main.go -- -file=$GOFILE
