package postgenerate_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ysuzuki19/robustruct/cmd/gen/internal/postgenerate"
)

func TestPostGenerate(t *testing.T) {
	require := require.New(t)
	// testdoc begin PostGenerate
	buf := []byte("package main\nfunc main() {\nprintln(\"testing\")}")
	formattedCode, err := postgenerate.PostGenerate(
		postgenerate.PostGenerateArgs{
			Buf: buf,
		},
	)
	// testdoc end
	require.NoError(err)
	require.NotEmpty(formattedCode)
	require.Equal("package main\n\nfunc main() {\n\tprintln(\"testing\")\n}\n", string(formattedCode))
}
