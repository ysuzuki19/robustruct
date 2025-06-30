package const_group_switch_cover

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func runInternal(pass *analysis.Pass, ss *ast.SwitchStmt, namedType *types.Named) {
	info := pass.TypesInfo
	name := namedType.Obj().Name()
	var tagType types.Type
	for ident, obj := range info.Defs {
		if typName, ok := obj.(*types.TypeName); ok && ident.Name == name {
			tagType = typName.Type()
			break
		}
	}
	if tagType == nil {
		fmt.Println("No matching type found for:", name)
		return
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
		return
	}
	fmt.Println("Constants found:", len(consts))

	cases := []string{}
	for _, stmt := range ss.Body.List {
		// Check if the case statement is a full arm
		caseStmt, ok := stmt.(*ast.CaseClause)
		if !ok || len(caseStmt.List) == 0 {
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
				return
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
