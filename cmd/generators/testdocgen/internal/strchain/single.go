package strchain

import "strings"

type single struct {
	s string
}

func From(s string) single {
	return single{s}
}

func (s single) String() string {
	return s.s
}

func (s single) TrimSpace() single {
	s.s = strings.TrimSpace(s.s)
	return s
}

func (s single) Replace(old, new string, n int) single {
	s.s = strings.Replace(s.s, old, new, n)
	return s
}

func (s single) Split(sep string) multiple {
	return FromSlice(strings.Split(s.String(), sep))
}
