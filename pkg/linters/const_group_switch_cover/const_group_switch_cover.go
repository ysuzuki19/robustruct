package const_group_switch_cover

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"

	"github.com/ysuzuki19/robustruct/pkg/linters/robustruct/settings"
)

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

func loadInfo(pass *analysis.Pass, file *ast.File) (*types.Info, error) {
	conf := types.Config{
		Context:          &types.Context{},
		GoVersion:        "",
		IgnoreFuncBodies: false,
		FakeImportC:      false,
		Error: func(err error) {
			panic("TODO")
		},
		Importer:                 nil,
		Sizes:                    nil,
		DisableUnusedImportCheck: false,
	}
	info := &types.Info{
		Types:        make(map[ast.Expr]types.TypeAndValue),
		Instances:    map[*ast.Ident]types.Instance{},
		Defs:         make(map[*ast.Ident]types.Object),
		Uses:         make(map[*ast.Ident]types.Object),
		Implicits:    map[ast.Node]types.Object{},
		Selections:   map[*ast.SelectorExpr]*types.Selection{},
		Scopes:       map[ast.Node]*types.Scope{},
		InitOrder:    []*types.Initializer{},
		FileVersions: map[*ast.File]string{},
	}
	_, err := conf.Check(pass.Fset.File(file.Pos()).Name(), pass.Fset, []*ast.File{file}, info)
	if err != nil {
		return nil, fmt.Errorf("failed to type check: %w", err)
	}

	return info, nil
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file == nil {
			continue
		}

		info, err := loadInfo(pass, file)
		if err != nil {
			return nil, fmt.Errorf("failed to load type info: %w", err)
		}

		// var err error
		ast.Inspect(file, func(n ast.Node) bool {
			if ss, ok := n.(*ast.SwitchStmt); ok {
				id, ok := ss.Tag.(*ast.Ident)
				if !ok || id == nil {
					return true
				}

				fmt.Println("Ident found:", id.Name)

				obj := info.Uses[id]
				if obj == nil {
					return true
				}
				typ := obj.Type()
				fmt.Println("Object found:", obj.Name(), "Type:", typ)

				named, ok := typ.(*types.Named)
				if !ok {
					// ignore check
					return true
				}
				name := named.Obj().Name()
				fmt.Println("Named type found:", name)

				var tagType types.Type
				for ident, obj := range info.Defs {
					if typName, ok := obj.(*types.TypeName); ok && ident.Name == name {
						tagType = typName.Type()
						break
					}
				}
				if tagType == nil {
					fmt.Println("No matching type found for:", name)
					return true
				}

				var consts []*types.Const
				for _, obj := range info.Defs {
					if obj == nil {
						continue
					}

					// Check if the object is a constant and matches the detected type
					if c, ok := obj.(*types.Const); ok && types.Identical(c.Type(), tagType) {
						consts = append(consts, c)
					}
				}
				if len(consts) == 0 {
					fmt.Println("No constants found for type:", tagType)
					return true
				}
				fmt.Println("Constants found:", len(consts))

				cases := []string{}
				for _, stmt := range ss.Body.List {
					// Check if the case statement is a full arm
					caseStmt, ok := stmt.(*ast.CaseClause)
					if !ok {
						continue
					}

					if len(caseStmt.List) == 0 {
						// This is a default case, which is always a full arm
						fmt.Println("Default case found")
						continue
					}

					for _, expr := range caseStmt.List {
						// Check if the expression is an identifier that matches the constant group
						if ident, ok := expr.(*ast.Ident); ok {
							cases = append(cases, ident.Name)
						} else {
							// If it's not an identifier, we can skip it for now
							fmt.Println("Non-identifier case expression found:", expr)
							pass.Report(analysis.Diagnostic{
								Pos:            ss.Pos(),
								End:            0,
								Category:       "",
								Message:        "robustruct/linters/switch_case_cover: case body requires const value",
								URL:            "",
								SuggestedFixes: []analysis.SuggestedFix{},
								Related:        []analysis.RelatedInformation{},
							})
							return true
						}
					}
				}
				fmt.Println("Cases found:", len(cases))

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
