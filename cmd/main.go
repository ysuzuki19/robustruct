package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/ysuzuki19/robustruct"
)

func main() {
	multichecker.Main(robustruct.Analyzer)
}
