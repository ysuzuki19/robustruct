package const_group_switch_cover

import (
	"flag"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/ysuzuki19/robustruct/internal/logger"
	"github.com/ysuzuki19/robustruct/pkg/linters/robustruct/settings"
)

func TypeEqual(a, b types.Type) bool {
	if a == nil || b == nil {
		return false
	}
	if a == b {
		return true
	}
	if namedA, ok := a.(*types.Named); ok {
		if namedB, ok := b.(*types.Named); ok {
			return namedA.Obj().Pkg() == namedB.Obj().Pkg() && namedA.Obj().Name() == namedB.Obj().Name()
		}
	}
	return false
}

var Analyzer = &analysis.Analyzer{
	Name: settings.FeatureConstGroupSwitchCover.String(),
	Doc:  "checks that all cases in a switch statement with a constant group expression are full arms",
	URL:  "",
	Flags: flag.FlagSet{
		Usage: func() {},
	},
	Run:              run,
	RunDespiteErrors: false,
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	ResultType:       nil,
	FactTypes:        []analysis.Fact{},
}

func findImport(pass *analysis.Pass, name string) *types.Package {
	for _, imported := range pass.Pkg.Imports() {
		importedName := imported.Name()
		if importedName == name {
			return imported
		}
	}
	return nil
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file == nil {
			continue
		}

		info := pass.TypesInfo

		ast.Inspect(file, func(n ast.Node) bool {
			if ss, ok := n.(*ast.SwitchStmt); ok {
				id, ok := ss.Tag.(*ast.Ident)
				if !ok || id == nil {
					return true
				}

				obj := info.Uses[id]
				if obj == nil {
					logger.Debug("No object found for ident:", id.Name)
					return true
				}
				typ := obj.Type()
				logger.Debug("Switch Tag Type:", typ)

				namedType, ok := typ.(*types.Named)
				if !ok {
					logger.Debug("[ignore-check] Switch Tag Type is not a named type")
					return true
				}
				logger.Debug("Named type found:", namedType.Obj().Pkg().Path(), ":", namedType.Obj().Name())

				isInternal := namedType.Obj().Pkg().Path() == pass.Pkg.Path()
				if isInternal {
					runInternal(pass, ss, namedType)
				} else {
					runExternal(pass, ss, namedType)
				}
			}
			return true
		})
	}
	return nil, nil
}

func Factory(flags settings.Flags) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:             Analyzer.Name,
		Doc:              Analyzer.Doc,
		URL:              Analyzer.URL,
		Flags:            flag.FlagSet{Usage: func() {}},
		Run:              run,
		RunDespiteErrors: false,
		Requires:         []*analysis.Analyzer{inspect.Analyzer},
		ResultType:       nil,
		FactTypes:        []analysis.Fact{},
	}
}
