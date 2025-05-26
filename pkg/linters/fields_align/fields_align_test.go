package fields_align_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/ysuzuki19/robustruct/pkg/linters/fields_align"
	"github.com/ysuzuki19/robustruct/pkg/linters/robustruct/settings"
)

func TestFix(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, fields_align.Analyzer, "fix")
}

func TestIgnore(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, fields_align.Factory(settings.Flags{settings.FlagDisableTest}), "ignore")
}
