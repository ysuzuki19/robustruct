// Package option provides a generic Option type, representing an optional value.
//
// The Option type is inspired by functional programming languages and is useful for:
// - Representing optional values (e.g., nullable values)
// - Avoiding the use of nil pointers
// - Providing a safer alternative to error-prone nil checks
//
// Example usage:
//
//	package main
//
//	import (
//	    "fmt"
//
//	    "github.com/ysuzuki19/robustruct/pkg/option"
//	)
//
//	func main() {
//	    someValue := option.NewSome(42)
//	    noneValue := option.None[int]()
//
//	    if someValue.IsSome() {
//	        fmt.Println("Some value:", *someValue.Ptr())
//	    }
//
//	    if noneValue.IsNone() {
//	        fmt.Println("No value")
//	    }
//	}
package option

// Option represents an optional value: every Option is either Some and contains a value, or None and does not.
//
// Example:
//
//	opt := option.NewSome(42)
//	fmt.Println(opt.IsSome()) // Output: true
//	fmt.Println(opt.IsNone()) // Output: false
//
//	none := option.None[int]()
//	fmt.Println(none.IsNone()) // Output: true
//	fmt.Println(none.IsSome()) // Output: false
type Option[T any] struct {
	ptr *T
}

// Check is a function type that takes a value of type T and returns a boolean.
// It is used in methods like IsSomeAnd and IsNoneOr.
type Check[T any] func(T) bool

// ValueFactory is a function type that produces a value of type T.
// It is used in methods like UnwrapOrElse and GetOrInsertWith.
type ValueFactory[T any] func() T

// OptionFactory is a function type that produces an Option[T].
// It is used in methods like OrElse.
type OptionFactory[T any] func() Option[T]

// Some creates an Option that contains the given value.
//
// Example:
//
//	v := 1
//	o := option.Some(&v)
//	require.True(o.IsSome())
//	require.Equal(&v, o.Ptr())
//	v = 2
//	require.Equal(2, *o.Ptr())
func Some[T any](v *T) Option[T] {
	return Option[T]{ptr: v}
}

// NewSome creates an Option that contains the given value.
//
// Example:
//
//	opt := option.NewSome(42)
//	fmt.Println(opt.IsSome()) // Output: true
func NewSome[T any](v T) Option[T] {
	return Some(&v)
}

// None creates an Option that does not contain a value.
//
// Example:
//
//	opt := option.None[int]()
//	fmt.Println(opt.IsNone()) // Output: true
func None[T any]() Option[T] {
	return Option[T]{ptr: nil}
}

// IsSome returns true if the Option contains a value.
//
// Example:
//
//	opt := option.NewSome(42)
//	fmt.Println(opt.IsSome()) // Output: true
func (o Option[T]) IsSome() bool {
	return o.ptr != nil
}

// IsNone returns true if the Option does not contain a value.
//
// Example:
//
//	opt := option.None[int]()
//	fmt.Println(opt.IsNone()) // Output: true
func (o Option[T]) IsNone() bool {
	return o.ptr == nil
}

// Ptr returns a pointer to the contained value, or nil if the Option is None.
//
// Example:
//
//	opt := option.NewSome(42)
//	fmt.Println(*opt.Ptr()) // Output: 42
func (o Option[T]) Ptr() *T {
	return o.ptr
}

// Get returns the contained value and a boolean indicating whether the Option is Some.
//
// Example:
//
//	if v, ok := opt.Get(); ok {
//		// something with `v`
//	} else {
//		// fallback without `v`
//	}
func (o Option[T]) Get() (*T, bool) {
	return o.ptr, o.IsSome()
}

// IsSomeAnd returns true if the Option is Some and the contained value satisfies the given predicate.
//
// Example:
//
//	require.True(option.NewSome(2).IsSomeAnd(func(x int) bool { return x > 1 }))
//	require.False(option.NewSome(0).IsSomeAnd(func(x int) bool { return x > 1 }))
//	require.False(option.None[int]().IsSomeAnd(func(x int) bool { return x > 1 }))
func (o Option[T]) IsSomeAnd(f Check[T]) bool {
	if o.IsSome() {
		return f(*o.ptr)
	}
	return false
}

// UnwrapOr returns the contained value if the Option is Some, or the given default value if the Option is None.
//
// Example:
//
//	require.Equal("car", option.NewSome("car").UnwrapOr("bike"))
//	require.Equal("bike", option.None[string]().UnwrapOr("bike"))
func (o Option[T]) UnwrapOr(v T) T {
	if o.IsNone() {
		return v
	}
	return *o.ptr
}

// UnwrapOrElse returns the contained value if the Option is Some, or invokes the given fallback function and returns its result if the Option is None.
//
// Example:
//
//	require.Equal(4, option.NewSome(4).UnwrapOrElse(func() int { return 20 }))
//	require.Equal(20, option.None[int]().UnwrapOrElse(func() int { return 20 }))
func (o Option[T]) UnwrapOrElse(f ValueFactory[T]) T {
	if o.IsNone() {
		return f()
	}
	return *o.ptr
}

// UnwrapOrDefault returns the contained value if the Option is Some, or the zero value of type T if the Option is None.
//
// Example:
//
//	require.Equal(0, option.None[int]().UnwrapOrDefault())
//	require.Equal(12, option.NewSome(12).UnwrapOrDefault())
func (o Option[T]) UnwrapOrDefault() T {
	if o.IsNone() {
		var v T // set zero value
		return v
	}
	return *o.ptr
}

// Filter returns the Option itself if it is Some and the contained value satisfies the given predicate, or None otherwise.
//
// Example:
//
//	isEven := func(x *int) bool { return *x%2 == 0 }
//
//	require.True(option.None[int]().Filter(isEven).IsNone())
//	require.True(option.NewSome(3).Filter(isEven).IsNone())
//	require.Equal(4, *option.NewSome(4).Filter(isEven).Ptr())
func (o Option[T]) Filter(f Check[*T]) Option[T] {
	if o.IsSome() && f(o.ptr) {
		return o
	}
	return None[T]()
}

// Or returns the Option itself if it is Some, or the given alternative Option if the Option is None.
//
// Example:
//
//	x := option.NewSome(2)
//	y := option.None[int]()
//	require.Equal(x.Or(y), x)
//
//	x = option.None[int]()
//	y = option.NewSome(3)
//	require.Equal(x.Or(y), y)
//
//	x = option.NewSome(2)
//	y = option.NewSome(3)
//	require.Equal(x.Or(y), x)
//
//	x = option.None[int]()
//	y = option.None[int]()
//	require.Equal(x.Or(y), y)
func (o Option[T]) Or(optb Option[T]) Option[T] {
	if o.IsSome() {
		return o
	}
	return optb
}

// OrElse returns the Option itself if it is Some, or invokes the given alternative factory function and returns its result if the Option is None.
//
// Example:
//
//	v := 3
//	noneFactory := func() option.Option[int] { return option.None[int]() }
//	someFactory := func() option.Option[int] { return option.NewSome(v) }
//
//	require.Equal(option.NewSome(2).OrElse(someFactory), option.NewSome(2))
//	require.Equal(option.None[int]().OrElse(someFactory), option.NewSome(v))
//	require.Equal(option.None[int]().OrElse(noneFactory), option.None[int]())
func (o Option[T]) OrElse(f OptionFactory[T]) Option[T] {
	if o.IsSome() {
		return o
	}
	return f()
}

// Xor returns the Option itself if it is Some and the given alternative Option is None, or the given alternative Option if it is Some and this Option is None. Returns None if both are Some or both are None.
//
// Example:
//
//	x := option.NewSome(2)
//	y := option.None[int]()
//	require.Equal(x.Xor(y), x)
func (o Option[T]) Xor(optb Option[T]) Option[T] {
	if o.IsSome() && optb.IsNone() {
		return o
	}
	if o.IsNone() && optb.IsSome() {
		return optb
	}
	return None[T]()
}

// Insert sets the Option to Some, containing the given value, and returns a pointer to the value.
//
// Example:
//
//	{
//		v := o.Insert(1)
//		require.Equal(o.Ptr(), v)
//		require.Equal(1, *v)
//		require.Equal(1, *o.Ptr())
//	}
//	{
//		v := o.Insert(2)
//		require.Equal(2, *v)
//		*v = 3
//		require.Equal(3, *o.Ptr())
//	}
func (o *Option[T]) Insert(v T) *T {
	o.ptr = &v
	return o.ptr
}

// GetOrInsert sets the Option to Some, containing the given value, if it is None, and returns a pointer to the value.
//
// Example:
//
//	o := option.None[int]()
//	{
//		v := o.GetOrInsert(1)
//		require.Equal(1, *v)
//		*v = 7
//	}
//	require.Equal(option.NewSome(7), o)
func (o *Option[T]) GetOrInsert(v T) *T {
	if o.IsNone() {
		o.ptr = &v
	}

	return o.ptr
}

// GetOrInsertDefault sets the Option to Some, containing the zero value of type T, if it is None, and returns a pointer to the value.
//
// Example:
//
//	o := option.None[int]()
//	{
//		v := o.GetOrInsertDefault()
//		require.Equal(0, *v)
//		*v = 7
//	}
//	require.Equal(option.NewSome(7), o)
func (o *Option[T]) GetOrInsertDefault() *T {
	var v T // set zero value
	return o.GetOrInsert(v)
}

// GetOrInsertWith sets the Option to Some, containing the value produced by the given factory function, if it is None, and returns a pointer to the value.
//
// Example:
//
//	o := option.None[int]()
//	{
//		v := o.GetOrInsertWith(func() int { return 1 })
//		require.Equal(1, *v)
//		*v = 7
//	}
//	require.Equal(option.NewSome(7), o)
func (o *Option[T]) GetOrInsertWith(f ValueFactory[T]) *T {
	if o.IsNone() {
		v := f()
		o.ptr = &v
	}

	return o.ptr
}

// Take takes the ownership of the contained value if the Option is Some, sets the Option to None, and returns an Option containing the previous value.
//
// Example:
//
//	o := option.NewSome(2)
//	o2 := o.Take()
//	require.Equal(option.None[int](), o)
//	require.Equal(option.NewSome(2), o2)
func (o *Option[T]) Take() Option[T] {
	if o.IsNone() {
		return None[T]()
	}
	v := o.ptr
	o.ptr = nil // change to None

	return Some(v)
}

// TakeIf takes the ownership of the contained value if the Option is Some and the value satisfies the given predicate, sets the Option to None, and returns an Option containing the previous value. Returns None if the Option is None or the predicate is not satisfied.
//
// Example:
//
//	o := option.NewSome(2)
//	{
//		prev := o.TakeIf(func(x *int) bool {
//			if *x == 2 {
//				*x += 1
//				return false
//			} else {
//				return false
//			}
//		})
//		require.Equal(option.NewSome(3), o)
//		require.Equal(option.None[int](), prev)
//	}
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

// Replace replaces the contained value with the given value and returns the previous Option value.
//
// Example:
//
//	o := option.NewSome(2)
//	prev := o.Replace(3)
//	require.Equal(option.NewSome(2), prev)
//	require.Equal(option.NewSome(3), o)
func (o *Option[T]) Replace(v T) Option[T] {
	prev := *o // copy current value
	o.ptr = &v
	return prev
}

// Clone creates a new Option containing a copy of the contained value if the Option is Some, or None if the Option is None.
//
// Example:
//
//	o := option.NewSome(2)
//	o2 := o.Clone()
//	require.Equal(option.NewSome(2), o)
//	require.Equal(option.NewSome(2), o2)
//	require.Equal(o, o2)
//	*o.Ptr() = 3 // o2 should not be affected
//	require.NotEqual(o, o2)
func (o Option[T]) Clone() Option[T] {
	if o.IsNone() {
		return None[T]()
	}
	copied := *o.ptr // copy current value
	return Some(&copied)
}

//go:generate go run ../../cmd/generators/testdocgen/main.go -- -file=$GOFILE
