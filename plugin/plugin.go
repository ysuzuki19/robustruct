package plugin

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/plugin-module-register/register"

	"github.com/ysuzuki19/robustruct/pkg/robustruct"
	"github.com/ysuzuki19/robustruct/pkg/robustruct/settings"
)

func init() {
	register.Plugin("robustruct", New)
}

type PluginRobustruct struct {
	settings settings.Settings
}

// ignore:fields_require
var _ register.LinterPlugin = &PluginRobustruct{}

func New(input any) (register.LinterPlugin, error) {
	settings, err := register.DecodeSettings[settings.Settings](input)
	if err != nil {
		return nil, err
	}
	return &PluginRobustruct{
		settings: settings,
	}, nil
}

func (pr *PluginRobustruct) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	analyzers := []*analysis.Analyzer{}
	for _, analyzer := range robustruct.FeatureAnalyzers {
		if pr.settings.Features.Contains(analyzer.Name) {
			analyzers = append(analyzers, analyzer)
		}
	}
	return analyzers, nil
}

func (pr *PluginRobustruct) GetLoadMode() string {
	return register.LoadModeSyntax
}
