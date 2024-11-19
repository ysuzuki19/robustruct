package golangci_lint

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/plugin-module-register/register"

	"github.com/ysuzuki19/robustruct/pkg/fields_align"
	"github.com/ysuzuki19/robustruct/pkg/fields_require"
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
	for _, feature := range pr.settings.Features {
		switch feature {
		case settings.FeatureFieldsRequire:
			analyzers = append(analyzers, fields_require.Factory(pr.settings.Flags))
		case settings.FeatureFieldsAlign:
			analyzers = append(analyzers, fields_align.Factory(pr.settings.Flags))
		}
	}
	return analyzers, nil
}

func (pr *PluginRobustruct) GetLoadMode() string {
	return register.LoadModeSyntax
}