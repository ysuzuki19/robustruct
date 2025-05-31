package strchain

import "strings"

type single struct {
	s string
}

// From
//
// Example:
//
//	s := strchain.From("testing")
//	require.Equal("testing", s.String())
func From(s string) single {
	return single{s}
}

// String
func (s single) String() string {
	return s.s
}

// TrimSpace
//
// Example:
//
//	s := strchain.From("  testing  ")
//	s = s.TrimSpace()
//	require.Equal("testing", s.String())
func (s single) TrimSpace() single {
	s.s = strings.TrimSpace(s.s)
	return s
}

// Replace
//
// Example:
//
//	s := strchain.From("testing")
//	s = s.Replace("t", "T", 1)
//	require.Equal("Testing", s.String())
func (s single) Replace(old, new string, n int) single {
	s.s = strings.Replace(s.s, old, new, n)
	return s
}

// Split
//
// Example:
//
//	s := strchain.From("a,b,c")
//	m := s.Split(",")
//	require.Equal([]string{"a", "b", "c"}, m.Collect())
func (s single) Split(sep string) multiple {
	return FromSlice(strings.Split(s.String(), sep))
}

//go:generate go run ../../main.go -- -file=$GOFILE
