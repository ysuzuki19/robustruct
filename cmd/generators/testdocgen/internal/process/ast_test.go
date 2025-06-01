package process

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
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

	decls := ListFnDecls(file)
	require.Len(decls, 5)
	require.Equal("Utility", decls[0].Name.Name)
	require.Equal("Method", decls[1].Name.Name)
	require.Equal("RefMethod", decls[2].Name.Name)
	require.Equal("GenericMethod", decls[3].Name.Name)
	require.Equal("GenericRefMethod", decls[4].Name.Name)
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
		name := recvTypeName(fn)
		require.Equal(c.receiver.UnwrapOrDefault(), name.UnwrapOrDefault())
	}
}
