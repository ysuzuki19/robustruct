package strchain_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ysuzuki19/robustruct/cmd/generators/testdocgen/internal/strchain"
)

func TestFrom(t *testing.T) {
	require := require.New(t)
	// testdoc begin From
	s := strchain.From("testing")
	require.Equal("testing", s.String())
	// testdoc end
}

func TestTrimSpace(t *testing.T) {
	require := require.New(t)
	// testdoc begin single.TrimSpace
	s := strchain.From("  testing  ")
	s = s.TrimSpace()
	require.Equal("testing", s.String())
	// testdoc end
}

func TestReplace(t *testing.T) {
	require := require.New(t)
	// testdoc begin single.Replace
	s := strchain.From("testing")
	s = s.Replace("t", "T", 1)
	require.Equal("Testing", s.String())
	// testdoc end
}

func TestMatch(t *testing.T) {
	require := require.New(t)
	// testdoc begin single.Match
	s := strchain.From("testing")
	re := regexp.MustCompile("^test")
	require.True(s.Match(re))
	// testdoc end
}

func TestMainchAndStrip(t *testing.T) {
	require := require.New(t)
	// testdoc begin single.MatchAndStrip
	s := strchain.From("testing")
	re := regexp.MustCompile("^test")
	m, ok := s.MatchAndStrip(re)
	require.True(ok)
	require.Equal("ing", m.String())
	// testdoc end
}

func TestSplit(t *testing.T) {
	require := require.New(t)
	// testdoc begin single.Split
	s := strchain.From("a,b,c")
	m := s.Split(",")
	require.Equal([]string{"a", "b", "c"}, m.Collect())
	// testdoc end
}
