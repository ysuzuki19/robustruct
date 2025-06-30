package const_group_switch_cover

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"

	"github.com/ysuzuki19/robustruct/internal/logger"
)

func runExternal(pass *analysis.Pass, ss *ast.SwitchStmt, namedType *types.Named) {
	info := pass.TypesInfo

	var consts []*types.Const
	typeName := namedType.Obj().Name()
	logger.Debug("Named type found:", typeName)
	pkgName := namedType.Obj().Pkg().Name()
	logger.Debug("Package name:", pkgName)

	tagTypePkg := findImport(pass, pkgName)
	if tagTypePkg == nil {
		logger.Debug("No import found for package:", pkgName)
		return
	}
	tagTypeScope := tagTypePkg.Scope()

	tagType := tagTypeScope.Lookup(typeName).Type()
	logger.Debug("Import object found:", tagType)

	for _, name := range tagTypeScope.Names() {
		obj := tagTypeScope.Lookup(name)
		if obj == nil {
			continue
		}

		if c, ok := obj.(*types.Const); ok {
			typeString := c.Type().String()
			logger.Debug("Constant type string:", typeString)
			if TypeEqual(c.Type(), tagType) {
				logger.Debug("Matching constant found:", c.Type().String())
				consts = append(consts, c)
			}
		}
	}

	if len(consts) == 0 {
		logger.Debug("No constants found for type:", tagType)
		return
	}
	logger.Debug("Constants found:", len(consts))

	cases := []types.Type{}
	for _, stmt := range ss.Body.List {
		caseStmt, ok := stmt.(*ast.CaseClause)
		if !ok {
			continue
		}

		if len(caseStmt.List) == 0 {
			continue
		}

		for _, expr := range caseStmt.List {
			caseType := info.Types[expr].Type
			if caseType == nil {
				fmt.Println("No type found for expression:", expr)
				return
			}

			if !TypeEqual(caseType, tagType) {
				pass.Report(analysis.Diagnostic{
					Pos:            ss.Pos(),
					End:            0,
					Category:       "",
					Message:        "robustruct/linters/switch_case_cover: case value requires type related const value",
					URL:            "",
					SuggestedFixes: []analysis.SuggestedFix{},
					Related:        []analysis.RelatedInformation{},
				})
				return
			}

			logger.Debug("Case expression type:", caseType.String())
			cases = append(cases, caseType)
		}
	}
	logger.Debug("Cases found:", len(cases))

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
