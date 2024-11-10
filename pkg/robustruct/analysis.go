package robustruct

import (
	"flag"

	"golang.org/x/tools/go/analysis"

	"github.com/ysuzuki19/robustruct/pkg/fields_align"
	"github.com/ysuzuki19/robustruct/pkg/fields_require"
)

var analyzers = []*analysis.Analyzer{
	fields_require.Analyzer,
	fields_align.Analyzer,
}

var Analyzer = &analysis.Analyzer{
	Name:             "robustruct",
	Doc:              "robustruct is a suite of analyzers for struct literals",
	URL:              "",
	Flags:            flag.FlagSet{Usage: func() {}},
	Run:              run,
	RunDespiteErrors: false,
	Requires:         analyzers,
	ResultType:       nil,
	FactTypes:        []analysis.Fact{},
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, analyzer := range analyzers {
		if _, err := analyzer.Run(pass); err != nil {
			return nil, err
		}
	}
	return nil, nil
}
