package strchain

import "strings"

type multiple struct {
	ss []string
}

func FromSlice(ss []string) multiple {
	return multiple{ss}
}

func (m multiple) Slice(start, end int) multiple {
	if end < 0 {
		end = len(m.ss)
	}
	if start < 0 || end > len(m.ss) || start > end {
		return FromSlice([]string{})
	}
	return FromSlice(m.ss[start:end])
}

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

func (m multiple) Map(fn func(string) string) multiple {
	for i := range m.ss {
		m.ss[i] = fn(m.ss[i])
	}
	return m
}

func (m multiple) Append(s ...string) multiple {
	m.ss = append(m.ss, s...)
	return m
}

func (m multiple) Extend(other multiple) multiple {
	m.ss = append(m.ss, other.ss...)
	return m
}

func (m multiple) Join(sep string) single {
	return From(strings.Join(m.ss, sep))
}

func (m multiple) Collect() []string {
	return m.ss
}
