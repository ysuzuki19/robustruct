package fields_require_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/ysuzuki19/robustruct/pkg/linters/fields_require"
	"github.com/ysuzuki19/robustruct/pkg/linters/robustruct/settings"
)

func TestFix(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, fields_require.Analyzer, "fix")
}

func TestIgnore(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, fields_require.Factory(settings.Flags{settings.FlagDisableTest}), "ignore")
}

func TestExternal(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, fields_require.Analyzer, "external")
}
