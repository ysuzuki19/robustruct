// Code generated by senum; DO NOT EDIT.
package result

type tag int

const (
	tagOk  tag = iota
	tagErr tag = iota
)

type Result[T any] struct {
	result[T]
	tag tag
}

func NewOk[T any](v *T) Result[T] {
	return Result[T]{
		result: result[T]{
			ok: v,
		},
		tag: tagOk,
	}
}

func NewErr[T any](v error) Result[T] {
	return Result[T]{
		result: result[T]{
			err: v,
		},
		tag: tagErr,
	}
}

func (e *Result[T]) IsOk() bool {
	return e.tag == tagOk
}

func (e *Result[T]) IsErr() bool {
	return e.tag == tagErr
}

func (e *Result[T]) AsOk() (*T, bool) {
	if e.IsOk() {
		return e.result.ok, true
	}
	return nil, false
}

func (e *Result[T]) AsErr() (error, bool) {
	if e.IsErr() {
		return e.result.err, true
	}
	return nil, false
}

type Switcher[T any] struct {
	Ok  func(v *T)
	Err func(v error)
}

func (e *Result[T]) Switch(s Switcher[T]) {
	switch e.tag {
	case tagOk:
		s.Ok(e.result.ok)
	case tagErr:
		s.Err(e.result.err)
	}
}

type Matcher[MatchResult any, T any] struct {
	Ok  func(v *T) MatchResult
	Err func(v error) MatchResult
}

func Match[MatchResult any, T any](e *Result[T], m Matcher[MatchResult, T]) MatchResult {
	switch e.tag {
	case tagOk:
		return m.Ok(e.result.ok)
	case tagErr:
		return m.Err(e.result.err)
	}
	panic("unreachable: invalid tag")
}
