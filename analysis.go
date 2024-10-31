package robustruct

import (
	"golang.org/x/tools/go/analysis"

	"github.com/ysuzuki19/robustruct/fields_align"
	"github.com/ysuzuki19/robustruct/fields_require"
)

var Analyzer = &analysis.Analyzer{
	Name: "robustruct",
	Doc:  "robustruct is a suite of analyzers for struct literals",
	Run:  run,
	Requires: []*analysis.Analyzer{
		fields_require.Analyzer,
		fields_align.Analyzer,
	},
}

var analyzers = []*analysis.Analyzer{
	fields_require.Analyzer,
	fields_align.Analyzer,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, analyzer := range analyzers {
		if _, err := analyzer.Run(pass); err != nil {
			return nil, err
		}
	}
	return nil, nil
}
