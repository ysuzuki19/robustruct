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

var Analyzer = &analysis.Analyzer{
	Name:             settings.FeatureConstGroupSwitchCover.String(),
	Doc:              "checks that all cases in a switch statement with a constant group expression are full arms",
	URL:              "",
	Flags:            flag.FlagSet{Usage: func() {}},
	Run:              run,
	RunDespiteErrors: false,
	Requires:         []*analysis.Analyzer{inspect.Analyzer},
	ResultType:       nil,
	FactTypes:        []analysis.Fact{},
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file == nil {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			if ss, ok := n.(*ast.SwitchStmt); ok {
				id, ok := ss.Tag.(*ast.Ident)
				if !ok || id == nil {
					return true
				}

				obj := pass.TypesInfo.Uses[id]
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

				cases := []types.Type{}
				for _, stmt := range ss.Body.List {
					caseStmt, ok := stmt.(*ast.CaseClause)
					if !ok || len(caseStmt.List) == 0 {
						continue
					}

					for _, expr := range caseStmt.List {
						if isHardCoded(expr) {
							logger.Debug("Hard-coded expression found:", expr)
							pass.Report(analysis.Diagnostic{
								Pos:            ss.Pos(),
								End:            0,
								Category:       "",
								Message:        "robustruct/linters/switch_case_cover: case value requires type related const value",
								URL:            "",
								SuggestedFixes: []analysis.SuggestedFix{},
								Related:        []analysis.RelatedInformation{},
							})
							return true
						}

						caseType := pass.TypesInfo.Types[expr].Type
						if caseType == nil {
							logger.Debug("No type found for expression:", expr)
							continue
						}

						logger.Debug("Case expression type:", caseType)
						cases = append(cases, caseType)
					}
				}
				logger.Debug("Cases found:", len(cases))

				isInternal := namedType.Obj().Pkg().Path() == pass.Pkg.Path()
				var consts []*types.Const
				if isInternal {
					consts = runInternal(pass, namedType)
				} else {
					consts = runExternal(pass, namedType)
				}

				if len(consts) != len(cases) {
					pass.Report(analysis.Diagnostic{
						Pos:            ss.Pos(),
						End:            0,
						Category:       "",
						Message:        "robustruct/linters/switch_case_cover: case body uncovered grouped const value",
						URL:            "",
						SuggestedFixes: []analysis.SuggestedFix{},
						Related:        []analysis.RelatedInformation{},
					})
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
