package const_group_switch_cover_test

import (
	"testing"

	"github.com/ysuzuki19/robustruct/pkg/linters/const_group_switch_cover"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLogLevel(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, const_group_switch_cover.Analyzer, "log_level")
}

func TestRole(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, const_group_switch_cover.Analyzer, "role")
}

func TestNamedNoConsts(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, const_group_switch_cover.Analyzer, "named_no_consts")
}

func TestExternal(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, const_group_switch_cover.Analyzer, "external")
}

func TestIgnore(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, const_group_switch_cover.Analyzer, "ignore")
}
