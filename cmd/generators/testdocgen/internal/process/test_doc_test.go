package process

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegexps(t *testing.T) {
	require := require.New(t)

	line := "// testdoc begin StructName.FuncName"
	require.True(tdRegex.MatchString(line), "Must to be matched: `%s`", line)

	line = "begin StructName.FuncName"
	require.True(tdBeginRegex.MatchString(line), "Must to be matched: `%s`", line)
	require.True(tdBeginRegex.MatchString(" "+line), "Must to be matched: ` %s`", line)

	line = "end"
	require.True(tdEndRegex.MatchString(line), "Must to be matched: `%s`", line)
	require.True(tdEndRegex.MatchString(" "+line), "Must to be matched: ` %s`", line)
}
