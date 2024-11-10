package robustruct

import (
	"flag"

	"golang.org/x/tools/go/analysis"

	"github.com/ysuzuki19/robustruct/pkg/fields_align"
	"github.com/ysuzuki19/robustruct/pkg/fields_require"
	"github.com/ysuzuki19/robustruct/pkg/robustruct/settings"
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

func Factory(settings settings.Settings) *analysis.Analyzer {
	enabled := []*analysis.Analyzer{}
	for _, a := range analyzers {
		if settings.Features.Contains(a.Name) {
			enabled = append(enabled, a)
		}
	}
	run := func(pass *analysis.Pass) (interface{}, error) {
		for _, analyzer := range enabled {
			if _, err := analyzer.Run(pass); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}
	analyzer := &analysis.Analyzer{
		Name:             Analyzer.Name,
		Doc:              Analyzer.Doc,
		URL:              Analyzer.URL,
		Flags:            flag.FlagSet{Usage: func() {}},
		Run:              run,
		RunDespiteErrors: false,
		Requires:         enabled,
		ResultType:       nil,
		FactTypes:        []analysis.Fact{},
	}
	return analyzer
}
