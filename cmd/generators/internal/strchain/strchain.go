package strchain

import "strings"

type String string
type Strings []string

func (s String) String() string {
	return string(s)
}

func (s String) Split(sep string) Strings {
	parts := strings.Split(string(s), sep)
	return Strings(parts)
}

func (s String) AsStrings() Strings {
	return Strings{string(s)}
}

func (ss Strings) Join(sep string) String {
	return String(strings.Join(ss, sep))
}

func (ss Strings) Map(fn func(string) string) Strings {
	mapped := make(Strings, len(ss))
	for i, s := range ss {
		mapped[i] = fn(s)
	}
	return mapped
}
