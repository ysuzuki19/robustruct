package plugin

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/plugin-module-register/register"

	"github.com/ysuzuki19/robustruct"
)

func init() {
	register.Plugin("robustruct", New)
}

type PluginRobustruct struct {
	analyzer *analysis.Analyzer
}

var _ register.LinterPlugin = &PluginRobustruct{}

func New(_ any) (register.LinterPlugin, error) {
	return &PluginRobustruct{
		analyzer: robustruct.Analyzer,
	}, nil
}

func (pr *PluginRobustruct) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{pr.analyzer}, nil
}

func (pr *PluginRobustruct) GetLoadMode() string {
	return register.LoadModeSyntax
}
