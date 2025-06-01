package engine

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ysuzuki19/robustruct/pkg/option"
)

func TestRegexps(t *testing.T) {
	require := require.New(t)

	cases := []struct {
		line string
		re   *regexp.Regexp
	}{
		{
			line: "// testdoc begin StructName.FuncName",
			re:   tdRegex,
		},
		{
			line: "begin StructName.FuncName",
			re:   tdBeginRegex,
		},
		{
			line: "end",
			re:   tdEndRegex,
		},
	}
	for _, c := range cases {
		require.True(c.re.MatchString(c.line), "Must to be matched: `%s`", c.line)
		require.True(c.re.MatchString(" "+c.line), "Must to be matched: `%s`", c.line)
		require.True(c.re.MatchString("\t"+c.line), "Must to be matched: `%s`", c.line)
		require.False(c.re.MatchString("x"+c.line), "Must to be unmatched: `%s`", c.line)
		require.False(c.re.MatchString("// "+c.line), "Must to be unmatched: `%s`", c.line)
	}
}

func TestParseTestDocsNormal(t *testing.T) {
	require := require.New(t)

	test := `
package sample_test

func TestUtility(t *testing.T) {
	require := require.New(t)
	// testdoc begin Utility
	require.Equal(1, Utility(1))
	// testdoc end
}

func TestSampleMethod(t *testing.T) {
	require := require.New(t)
	// testdoc begin Sample.Method
	s := Sample{}
	require.Equal(1, s.Method())
	// testdoc end
}
	`

	tds, err := ParseTestDocs(test)
	require.NoError(err)
	require.Equal([]TestDoc{
		{
			StructName: option.None[string](),
			FuncName:   "Utility",
			Content:    "\trequire.Equal(1, Utility(1))",
		},
		{
			StructName: option.NewSome("Sample"),
			FuncName:   "Method",
			Content:    "\ts := Sample{}\n\trequire.Equal(1, s.Method())",
		},
	}, tds)
}

func TestParseTestDocsLines(t *testing.T) {
	cases := []struct {
		source string
		isErr  bool
	}{
		{
			source: `
			`,
			isErr: false,
		},
		{
			source: `
				// testdoc begin Utility
				// testdoc end
			`,
			isErr: false,
		},
		{
			source: `
				// testdoc begin Sample.Method
				// testdoc end
			`,
			isErr: false,
		},
		{
			source: `
				// testdoc begin Sample.
				// testdoc end
			`,
			isErr: true,
		},
		{
			source: `
				// testdoc begin .Method
				// testdoc end
			`,
			isErr: true,
		},
		{
			source: `
				// testdoc begin
				// testdoc end
			`,
			isErr: true,
		},
		{
			source: `
				// testdoc begin Utility
				// testdoc begin Utility
				// testdoc end
			`,
			isErr: true,
		},
		{
			source: `
				// testdoc end
			`,
			isErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.source, func(t *testing.T) {
			require := require.New(t)
			_, err := ParseTestDocs(c.source)
			if c.isErr {
				require.Error(err)
			} else {
				require.NoError(err)
			}
		})
	}
}
