package strchain

import (
	"regexp"
	"strings"
)

type single struct {
	s string
}

// From returns a string chaining object
//
// Example:
//
//	s := strchain.From("testing")
//	require.Equal("testing", s.String())
func From(s string) single {
	return single{s}
}

// String returns the standard string
func (s single) String() string {
	return s.s
}

// TrimSpace returns a strings.TrimSpace applied string.
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

// Replace returns a strings.Replace applied string.
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

// Match returns true if the string matches the regular expression.
//
// Example:
//
//	s := strchain.From("testing")
//	re := regexp.MustCompile("^test")
//	require.True(s.Match(re))
func (s single) Match(re *regexp.Regexp) bool {
	return re.MatchString(s.s)
}

// MatchAndStrip returns a new single with the matched substring removed.
//
// Example:
//
//	s := strchain.From("testing")
//	re := regexp.MustCompile("^test")
//	m, ok := s.MatchAndStrip(re)
//	require.True(ok)
//	require.Equal("ing", m.String())
func (s single) MatchAndStrip(re *regexp.Regexp) (single, bool) {
	loc := re.FindStringIndex(s.s)
	if loc == nil {
		return s, false
	}
	s.s = s.s[:loc[0]] + s.s[loc[1]:]
	return s, true
}

// Split returns a multiple containing the substrings of the string
//
// Example:
//
//	s := strchain.From("a,b,c")
//	m := s.Split(",")
//	require.Equal([]string{"a", "b", "c"}, m.Collect())
func (s single) Split(sep string) multiple {
	return FromSlice(strings.Split(s.String(), sep))
}

//go:generate go run github.com/ysuzuki19/robustruct/cmd/generators/testdocgen -file=$GOFILE
