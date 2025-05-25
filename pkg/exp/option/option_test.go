// torin DELETE BEGIN date=2025-10-01
package option_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ysuzuki19/robustruct/pkg/exp/option"
)

func Test_Some(t *testing.T) {
	require := require.New(t)
	v := 1
	o := option.Some(&v)
	require.True(o.IsSome())
	require.Equal(&v, o.Ptr())
	v = 2
	require.Equal(2, *o.Ptr())
}

func Test_NewSome(t *testing.T) {
	require := require.New(t)
	o := option.NewSome(1)
	require.True(o.IsSome())
	require.Equal(1, *o.Ptr())
}

func Test_None(t *testing.T) {
	require := require.New(t)
	o := option.None[int]()
	require.True(o.IsNone())
	require.Nil(o.Ptr())
}

func Test_Get(t *testing.T) {
	require := require.New(t)
	v := 1
	ptr, isSome := option.Some(&v).Get()
	require.True(isSome)
	require.Equal(&v, ptr)

	ptr, isSome = option.None[int]().Get()
	require.False(isSome)
	require.Nil(ptr)
}

func Test_IsSomeAnd(t *testing.T) {
	require := require.New(t)
	require.True(option.NewSome(2).IsSomeAnd(func(x int) bool { return x > 1 }))
	require.False(option.NewSome(0).IsSomeAnd(func(x int) bool { return x > 1 }))
	require.False(option.None[int]().IsSomeAnd(func(x int) bool { return x > 1 }))
}

func Test_IsNoneOr(t *testing.T) {
	require := require.New(t)
	require.True(option.NewSome(2).IsNoneOr(func(x int) bool { return x > 1 }))
	require.False(option.NewSome(0).IsNoneOr(func(x int) bool { return x > 1 }))
	require.True(option.None[int]().IsNoneOr(func(x int) bool { return x > 1 }))
}

func Test_UnwrapOr(t *testing.T) {
	require := require.New(t)
	require.Equal("car", option.NewSome("car").UnwrapOr("bike"))
	require.Equal("bike", option.None[string]().UnwrapOr("bike"))
}

func Test_UnwrapOrElse(t *testing.T) {
	require := require.New(t)
	k := 10
	require.Equal(4, option.NewSome(4).UnwrapOrElse(func() int { return k * 2 }))
	require.Equal(20, option.None[int]().UnwrapOrElse(func() int { return k * 2 }))
}

func Test_UnwrapOrDefault(t *testing.T) {
	require := require.New(t)
	require.Equal(0, option.None[int]().UnwrapOrDefault())
	require.Equal(12, option.NewSome(12).UnwrapOrDefault())
}

func Test_Filter(t *testing.T) {
	require := require.New(t)
	isEven := func(x *int) bool { return *x%2 == 0 }

	require.True(option.None[int]().Filter(isEven).IsNone())
	require.True(option.NewSome(3).Filter(isEven).IsNone())
	require.Equal(4, *option.NewSome(4).Filter(isEven).Ptr())
}

func Test_Or(t *testing.T) {
	require := require.New(t)

	x := option.NewSome(2)
	y := option.None[int]()
	require.Equal(x.Or(y), x)

	x = option.None[int]()
	y = option.NewSome(3)
	require.Equal(x.Or(y), y)

	x = option.NewSome(2)
	y = option.NewSome(3)
	require.Equal(x.Or(y), x)

	x = option.None[int]()
	y = option.None[int]()
	require.Equal(x.Or(y), y)
}

func Test_OrElse(t *testing.T) {
	require := require.New(t)

	v := 3
	noneFactory := func() option.Option[int] { return option.None[int]() }
	someFactory := func() option.Option[int] { return option.NewSome(v) }

	require.Equal(option.NewSome(2).OrElse(someFactory), option.NewSome(2))
	require.Equal(option.None[int]().OrElse(someFactory), option.NewSome(v))
	require.Equal(option.None[int]().OrElse(noneFactory), option.None[int]())
}

func Test_Xor(t *testing.T) {
	require := require.New(t)
	{
		x := option.NewSome(2)
		y := option.None[int]()
		require.Equal(x.Xor(y), x)
	}
	{
		x := option.None[int]()
		y := option.NewSome(3)
		require.Equal(x.Xor(y), y)
	}
	{
		x := option.NewSome(2)
		y := option.NewSome(3)
		require.Equal(x.Xor(y), option.None[int]())
	}
	{
		x := option.None[int]()
		y := option.None[int]()
		require.Equal(x.Xor(y), option.None[int]())
	}
}

func Test_Insert(t *testing.T) {
	require := require.New(t)
	o := option.None[int]()
	{
		v := o.Insert(1)
		require.Equal(o.Ptr(), v)
		require.Equal(1, *v)
		require.Equal(1, *o.Ptr())
	}
	{
		v := o.Insert(2)
		require.Equal(2, *v)
		*v = 3
		require.Equal(3, *o.Ptr())
	}
}

func Test_GetOrInsert(t *testing.T) {
	require := require.New(t)
	o := option.None[int]()
	{
		v := o.GetOrInsert(1)
		require.Equal(1, *v)
		*v = 7
	}
	require.Equal(option.NewSome(7), o)
}

func Test_GetOrInsertDefault(t *testing.T) {
	require := require.New(t)
	o := option.None[int]()
	{
		v := o.GetOrInsertDefault()
		require.Equal(0, *v)
		*v = 7
	}
	require.Equal(option.NewSome(7), o)
}

func Test_GetOrInsertWith(t *testing.T) {
	require := require.New(t)
	o := option.None[int]()
	{
		v := o.GetOrInsertWith(func() int { return 1 })
		require.Equal(1, *v)
		*v = 7
	}
	require.Equal(option.NewSome(7), o)
}

func Test_Take(t *testing.T) {
	require := require.New(t)
	{
		o := option.NewSome(2)
		o2 := o.Take()
		require.Equal(option.None[int](), o)
		require.Equal(option.NewSome(2), o2)
	}

	{
		o := option.None[int]()
		o2 := o.Take()
		require.Equal(option.None[int](), o)
		require.Equal(option.None[int](), o2)
	}
}

func Test_TakeIf(t *testing.T) {
	require := require.New(t)
	o := option.NewSome(2)
	{
		prev := o.TakeIf(func(x *int) bool {
			if *x == 2 {
				*x += 1
				return false
			} else {
				return false
			}
		})
		require.Equal(option.NewSome(3), o)
		require.Equal(option.None[int](), prev)
	}
	{
		prev := o.TakeIf(func(x *int) bool {
			require.Equal(3, *x)
			return *x == 3
		})
		require.Equal(option.None[int](), o)
		require.Equal(option.NewSome(3), prev)
	}
}

func Test_Replace(t *testing.T) {
	require := require.New(t)
	{
		o := option.NewSome(2)
		prev := o.Replace(3)
		require.Equal(option.NewSome(2), prev)
		require.Equal(option.NewSome(3), o)
	}
	{
		o := option.None[int]()
		prev := o.Replace(3)
		require.Equal(option.None[int](), prev)
		require.Equal(option.NewSome(3), o)
	}
}

func Test_Clone(t *testing.T) {
	require := require.New(t)
	{
		o := option.NewSome(2)
		o2 := o.Clone()
		require.Equal(option.NewSome(2), o)
		require.Equal(option.NewSome(2), o2)
		require.Equal(o, o2)
		*o.Ptr() = 3 // o2 should not be affected
		require.NotEqual(o, o2)
	}
}

// torin DELETE END
