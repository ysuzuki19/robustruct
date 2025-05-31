package strchain

import "strings"

type single struct {
	s string
}

// From
func From(s string) single {
	return single{s}
}

// String
func (s single) String() string {
	return s.s
}

// TrimSpace
func (s single) TrimSpace() single {
	s.s = strings.TrimSpace(s.s)
	return s
}

// Replace
func (s single) Replace(old, new string, n int) single {
	s.s = strings.Replace(s.s, old, new, n)
	return s
}

// Split
func (s single) Split(sep string) multiple {
	return FromSlice(strings.Split(s.String(), sep))
}

//go:generate go run ../../main.go -- -file=$GOFILE
