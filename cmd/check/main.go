package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/ysuzuki19/robustruct/pkg/robustruct"
)

func main() {
	multichecker.Main(robustruct.FeatureAnalyzers...)
}
