package e2e_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/ysuzuki19/robustruct/pkg/linters/robustruct"
)

func TestE2E(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, robustruct.Analyzer, "e2e")
}

func TestSample(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, robustruct.Analyzer, "sample")
}
