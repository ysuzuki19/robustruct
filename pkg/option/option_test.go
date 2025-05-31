package option_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ysuzuki19/robustruct/pkg/option"
)

func Test_Some(t *testing.T) {
	require := require.New(t)
	// testdoc begin Some
	v := 1
	o := option.Some(&v)
	require.True(o.IsSome())
	require.Equal(&v, o.Ptr())
	v = 2
	require.Equal(2, *o.Ptr())
	// testdoc end
}

func Test_NewSome(t *testing.T) {
	require := require.New(t)
	// testdoc begin NewSome
	o := option.NewSome(1)
	require.True(o.IsSome())
	require.Equal(1, *o.Ptr())
	// testdoc end
}

func Test_None(t *testing.T) {
	require := require.New(t)
	// testdoc begin None
	o := option.None[int]()
	require.True(o.IsNone())
	require.Nil(o.Ptr())
	// testdoc end
}

func Test_Get(t *testing.T) {
	require := require.New(t)
	// testdoc begin Option.Get
	v := 1
	ptr, isSome := option.Some(&v).Get()
	require.True(isSome)
	require.Equal(&v, ptr)

	ptr, isSome = option.None[int]().Get()
	require.False(isSome)
	require.Nil(ptr)
	// testdoc end
}

func Test_IsSomeAnd(t *testing.T) {
	require := require.New(t)
	// testdoc begin Option.IsSomeAnd
	require.True(option.NewSome(2).IsSomeAnd(func(x int) bool { return x > 1 }))
	require.False(option.NewSome(0).IsSomeAnd(func(x int) bool { return x > 1 }))
	require.False(option.None[int]().IsSomeAnd(func(x int) bool { return x > 1 }))
	// testdoc end
}

func Test_UnwrapOr(t *testing.T) {
	require := require.New(t)
	// testdoc begin Option.UnwrapOr
	require.Equal("car", option.NewSome("car").UnwrapOr("bike"))
	require.Equal("bike", option.None[string]().UnwrapOr("bike"))
	// testdoc end
}

func Test_UnwrapOrElse(t *testing.T) {
	require := require.New(t)
	// testdoc begin Option.UnwrapOrElse
	require.Equal(4, option.NewSome(4).UnwrapOrElse(func() int { return 20 }))
	require.Equal(20, option.None[int]().UnwrapOrElse(func() int { return 20 }))
	// testdoc end
}

func Test_UnwrapOrDefault(t *testing.T) {
	require := require.New(t)
	// testdoc begin Option.UnwrapOrDefault
	require.Equal(0, option.None[int]().UnwrapOrDefault())
	require.Equal(12, option.NewSome(12).UnwrapOrDefault())
	// testdoc end
}

func Test_Filter(t *testing.T) {
	require := require.New(t)
	// testdoc begin Option.Filter
	isEven := func(x *int) bool { return *x%2 == 0 }

	require.True(option.None[int]().Filter(isEven).IsNone())
	require.True(option.NewSome(3).Filter(isEven).IsNone())
	require.Equal(4, *option.NewSome(4).Filter(isEven).Ptr())
	// testdoc end
}

func Test_Or(t *testing.T) {
	require := require.New(t)

	// testdoc begin Option.Or
	x := option.NewSome(2)
	y := option.None[int]()
	require.Equal(x.Or(y), x)

	x = option.None[int]()
	y = option.NewSome(3)
	require.Equal(x.Or(y), y)
	// testdoc end

	x = option.NewSome(2)
	y = option.NewSome(3)
	require.Equal(x.Or(y), x)

	x = option.None[int]()
	y = option.None[int]()
	require.Equal(x.Or(y), y)
}

func Test_OrElse(t *testing.T) {
	require := require.New(t)

	// testdoc begin Option.OrElse
	v := 3
	noneFactory := func() option.Option[int] { return option.None[int]() }
	someFactory := func() option.Option[int] { return option.NewSome(v) }

	require.Equal(option.NewSome(2).OrElse(someFactory), option.NewSome(2))
	require.Equal(option.None[int]().OrElse(someFactory), option.NewSome(v))
	require.Equal(option.None[int]().OrElse(noneFactory), option.None[int]())
	// testdoc end
}

func Test_Xor(t *testing.T) {
	require := require.New(t)
	{
		// testdoc begin Option.Xor
		x := option.NewSome(2)
		y := option.None[int]()
		require.Equal(x.Xor(y), x)
		// testdoc end
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
		// testdoc begin Option.Insert
		v := o.Insert(1)
		require.Equal(o.Ptr(), v)
		require.Equal(1, *v)
		require.Equal(1, *o.Ptr())
		// testdoc end
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
		// testdoc begin Option.GetOrInsert
		v := o.GetOrInsert(1)
		require.Equal(1, *v)
		*v = 7
		// testdoc end
	}
	require.Equal(option.NewSome(7), o)
}

func Test_GetOrInsertDefault(t *testing.T) {
	require := require.New(t)
	o := option.None[int]()
	{
		// testdoc begin Option.GetOrInsertDefault
		v := o.GetOrInsertDefault()
		require.Equal(0, *v)
		*v = 7
		// testdoc end
	}
	require.Equal(option.NewSome(7), o)
}

func Test_GetOrInsertWith(t *testing.T) {
	require := require.New(t)
	// testdoc begin Option.GetOrInsertWith
	o := option.None[int]()
	{
		v := o.GetOrInsertWith(func() int { return 1 })
		require.Equal(1, *v)
		*v = 7
	}
	require.Equal(option.NewSome(7), o)
	// testdoc end
}

func Test_Take(t *testing.T) {
	require := require.New(t)
	{
		// testdoc begin Option.Take
		o := option.NewSome(2)
		o2 := o.Take()
		require.Equal(option.None[int](), o)
		require.Equal(option.NewSome(2), o2)
		// testdoc end
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
		// testdoc begin Option.TakeIf
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
		// testdoc end
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
		// testdoc begin Option.Replace
		o := option.NewSome(2)
		prev := o.Replace(3)
		require.Equal(option.NewSome(2), prev)
		require.Equal(option.NewSome(3), o)
		// testdoc end
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
		// testdoc begin Option.Clone
		o := option.NewSome(2)
		o2 := o.Clone()
		require.Equal(option.NewSome(2), o)
		require.Equal(option.NewSome(2), o2)
		require.Equal(o, o2)
		*o.Ptr() = 3 // o2 should not be affected
		require.NotEqual(o, o2)
		// testdoc end
	}
}
