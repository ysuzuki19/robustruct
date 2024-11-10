package fields_align_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/ysuzuki19/robustruct/pkg/fields_align"
)

func TestFix(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, fields_align.Analyzer, "fix")
}

func TestIgnore(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, fields_align.Analyzer, "ignore")
}
