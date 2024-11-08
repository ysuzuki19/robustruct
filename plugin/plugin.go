package plugin

import (
	"slices"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/plugin-module-register/register"

	"github.com/ysuzuki19/robustruct"
)

func init() {
	register.Plugin("robustruct", New)
}

type Settings struct {
	Disables []string `json:"disables"`
}

type PluginRobustruct struct {
	settings  Settings
	analyzers []*analysis.Analyzer
}

// ignore:fields_require
var _ register.LinterPlugin = &PluginRobustruct{}

func New(input any) (register.LinterPlugin, error) {
	settings, err := register.DecodeSettings[Settings](input)
	if err != nil {
		return nil, err
	}
	return &PluginRobustruct{
		settings:  settings,
		analyzers: robustruct.Analyzer.Requires,
	}, nil
}

func (pr *PluginRobustruct) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	analyzers := []*analysis.Analyzer{}
	for _, analyzer := range pr.analyzers {
		if !slices.Contains(pr.settings.Disables, analyzer.Name) {
			analyzers = append(analyzers, analyzer)
		}
	}
	return analyzers, nil
}

func (pr *PluginRobustruct) GetLoadMode() string {
	return register.LoadModeSyntax
}
