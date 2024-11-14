package option

type Option[T any] struct {
	ptr *T
}

type Check[T any] func(T) bool
type ValueFactory[T any] func() T
type OptionFactory[T any] func() Option[T]

func Some[T any](v *T) Option[T] {
	return Option[T]{ptr: v}
}

func NewSome[T any](v T) Option[T] {
	return Some(&v)
}

func None[T any]() Option[T] {
	return Option[T]{ptr: nil}
}

func (e Option[T]) IsSome() bool {
	return e.ptr != nil
}

func (e Option[T]) IsNone() bool {
	return e.ptr == nil
}

func (o Option[T]) Ptr() *T {
	return o.ptr
}

func (o Option[T]) Get() (*T, bool) {
	return o.ptr, o.IsSome()
}

func (o Option[T]) IsSomeAnd(f Check[T]) bool {
	if o.IsSome() {
		return f(*o.ptr)
	}
	return false
}

func (o Option[T]) IsNoneOr(f Check[T]) bool {
	if o.IsNone() {
		return true
	}
	return f(*o.ptr)
}

func (o Option[T]) UnwrapOr(v T) T {
	if o.IsNone() {
		return v
	}
	return *o.ptr
}

func (o Option[T]) UnwrapOrElse(f ValueFactory[T]) T {
	if o.IsNone() {
		return f()
	}
	return *o.ptr
}

func (o Option[T]) UnwrapOrDefault() T {
	if o.IsNone() {
		var v T // set zero value
		return v
	}
	return *o.ptr
}

func (o Option[T]) Filter(f Check[*T]) Option[T] {
	if o.IsSome() && f(o.ptr) {
		return o
	}
	return None[T]()
}

func (o Option[T]) Or(optb Option[T]) Option[T] {
	if o.IsSome() {
		return o
	}
	return optb
}

func (o Option[T]) OrElse(f OptionFactory[T]) Option[T] {
	if o.IsSome() {
		return o
	}
	return f()
}

func (o Option[T]) Xor(optb Option[T]) Option[T] {
	if o.IsSome() && optb.IsNone() {
		return o
	}
	if o.IsNone() && optb.IsSome() {
		return optb
	}
	return None[T]()
}

func (o *Option[T]) Insert(v T) *T {
	o.ptr = &v
	return o.ptr
}

func (o *Option[T]) GetOrInsert(v T) *T {
	if o.IsNone() {
		o.ptr = &v
	}

	return o.ptr
}

func (o *Option[T]) GetOrInsertDefault() *T {
	var v T // set zero value
	return o.GetOrInsert(v)
}

func (o *Option[T]) GetOrInsertWith(f ValueFactory[T]) *T {
	if o.IsNone() {
		v := f()
		o.ptr = &v
	}

	return o.ptr
}

func (o *Option[T]) Take() Option[T] {
	if o.IsNone() {
		return None[T]()
	}
	v := o.ptr
	o.ptr = nil // change to None

	return Some[T](v)
}

func (o *Option[T]) TakeIf(predicate Check[*T]) Option[T] {
	if o.IsNone() {
		return None[T]()
	}
	v := o.ptr
	if !predicate(o.ptr) {
		return None[T]()
	}
	o.ptr = nil // change to None
	return Some[T](v)
}

func (o *Option[T]) Replace(v T) Option[T] {
	prev := *o // copy current value
	o.ptr = &v
	return prev
}

func (o Option[T]) Clone() Option[T] {
	if o.IsNone() {
		return None[T]()
	}
	copied := *o.ptr // copy current value
	return Some(&copied)
}
