package astutil_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ysuzuki19/robustruct/cmd/generators/testdocgen/internal/process/astutil"
	"github.com/ysuzuki19/robustruct/cmd/generators/testdocgen/internal/strchain"
	"github.com/ysuzuki19/robustruct/pkg/option"
)

func createSource(parts ...string) string {
	return strchain.
		FromSlice([]string{
			"package main",
		}).
		Append(parts...).
		Join("\n").
		String()
}

func TestListFnDecls(t *testing.T) {
	require := require.New(t)

	source := `
package sample
func Utility() {}
func (s Sample) Method() {}
func (s *Sample) RefMethod() {}
func (s Sample[T]) GenericMethod() {}
func (s *Sample[T]) GenericRefMethod() {}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	require.NoError(err)

	funcs := astutil.ListFnDecls(fset, file)
	require.Len(funcs, 5)

	require.Equal("Utility", funcs[0].Name)
	require.Equal(option.None[string](), funcs[0].Recv)

	require.Equal("Method", funcs[1].Name)
	require.Equal(option.NewSome("Sample"), funcs[1].Recv)

	require.Equal("Method", funcs[1].Name)
	require.Equal(option.NewSome("Sample"), funcs[1].Recv)

	require.Equal("RefMethod", funcs[2].Name)
	require.Equal(option.NewSome("Sample"), funcs[2].Recv)

	require.Equal("GenericMethod", funcs[3].Name)
	require.Equal(option.NewSome("Sample"), funcs[3].Recv)

	require.Equal("GenericRefMethod", funcs[4].Name)
	require.Equal(option.NewSome("Sample"), funcs[4].Recv)
}

func TestRecvTypeName(t *testing.T) {
	require := require.New(t)

	cases := []struct {
		decl     string
		receiver option.Option[string]
	}{
		{
			decl:     "func Utility() {}",
			receiver: option.None[string](),
		},
		{
			decl:     "func (s Sample) Method() {}",
			receiver: option.NewSome("Sample"),
		},
		{
			decl:     "func (s *Sample) RefMethod() {}",
			receiver: option.NewSome("Sample"),
		},
		{
			decl:     "func (s Sample[T]) GenericMethod() {}",
			receiver: option.NewSome("Sample"),
		},
		{
			decl:     "func (s *Sample[T]) GenericRefMethod() {}",
			receiver: option.NewSome("Sample"),
		},
	}

	for _, c := range cases {
		source := createSource(c.decl)

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "", source, parser.ParseComments)
		require.NoError(err)

		if len(file.Decls) != 1 {
			require.Fail("no declarations found")
			continue
		}

		fn, ok := file.Decls[0].(*ast.FuncDecl)
		require.True(ok)
		name := astutil.RecvTypeName(fn)
		require.Equal(c.receiver.UnwrapOrDefault(), name.UnwrapOrDefault())
	}
}
