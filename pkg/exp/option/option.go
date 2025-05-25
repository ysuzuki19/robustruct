// torin DELETE BEGIN date=2025-10-01
package option

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead
type Option[T any] struct {
	ptr *T
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
type Check[T any] func(T) bool

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
type ValueFactory[T any] func() T

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
type OptionFactory[T any] func() Option[T]

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func Some[T any](v *T) Option[T] {
	return Option[T]{ptr: v}
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func NewSome[T any](v T) Option[T] {
	return Some(&v)
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func None[T any]() Option[T] {
	return Option[T]{ptr: nil}
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o Option[T]) IsSome() bool {
	return o.ptr != nil
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o Option[T]) IsNone() bool {
	return o.ptr == nil
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o Option[T]) Ptr() *T {
	return o.ptr
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o Option[T]) Get() (*T, bool) {
	return o.ptr, o.IsSome()
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o Option[T]) IsSomeAnd(f Check[T]) bool {
	if o.IsSome() {
		return f(*o.ptr)
	}
	return false
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o Option[T]) IsNoneOr(f Check[T]) bool {
	if o.IsNone() {
		return true
	}
	return f(*o.ptr)
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o Option[T]) UnwrapOr(v T) T {
	if o.IsNone() {
		return v
	}
	return *o.ptr
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o Option[T]) UnwrapOrElse(f ValueFactory[T]) T {
	if o.IsNone() {
		return f()
	}
	return *o.ptr
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o Option[T]) UnwrapOrDefault() T {
	if o.IsNone() {
		var v T // set zero value
		return v
	}
	return *o.ptr
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o Option[T]) Filter(f Check[*T]) Option[T] {
	if o.IsSome() && f(o.ptr) {
		return o
	}
	return None[T]()
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o Option[T]) Or(optb Option[T]) Option[T] {
	if o.IsSome() {
		return o
	}
	return optb
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o Option[T]) OrElse(f OptionFactory[T]) Option[T] {
	if o.IsSome() {
		return o
	}
	return f()
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o Option[T]) Xor(optb Option[T]) Option[T] {
	if o.IsSome() && optb.IsNone() {
		return o
	}
	if o.IsNone() && optb.IsSome() {
		return optb
	}
	return None[T]()
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o *Option[T]) Insert(v T) *T {
	o.ptr = &v
	return o.ptr
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o *Option[T]) GetOrInsert(v T) *T {
	if o.IsNone() {
		o.ptr = &v
	}

	return o.ptr
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o *Option[T]) GetOrInsertDefault() *T {
	var v T // set zero value
	return o.GetOrInsert(v)
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o *Option[T]) GetOrInsertWith(f ValueFactory[T]) *T {
	if o.IsNone() {
		v := f()
		o.ptr = &v
	}

	return o.ptr
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o *Option[T]) Take() Option[T] {
	if o.IsNone() {
		return None[T]()
	}
	v := o.ptr
	o.ptr = nil // change to None

	return Some(v)
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o *Option[T]) TakeIf(predicate Check[*T]) Option[T] {
	if o.IsNone() {
		return None[T]()
	}
	v := o.ptr
	if !predicate(o.ptr) {
		return None[T]()
	}
	o.ptr = nil // change to None
	return Some(v)
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o *Option[T]) Replace(v T) Option[T] {
	prev := *o // copy current value
	o.ptr = &v
	return prev
}

// Deprecated: use github.com/ysuzuki19/robustruct/pkg/option instead of github.com/ysuzuki19/robustruct/pkg/exp/option
func (o Option[T]) Clone() Option[T] {
	if o.IsNone() {
		return None[T]()
	}
	copied := *o.ptr // copy current value
	return Some(&copied)
}

// torin DELETE END
