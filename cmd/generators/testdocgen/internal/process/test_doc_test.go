package process

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
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
