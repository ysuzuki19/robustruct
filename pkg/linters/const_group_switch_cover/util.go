package const_group_switch_cover

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func findImport(pass *analysis.Pass, name string) *types.Package {
	for _, imported := range pass.Pkg.Imports() {
		importedName := imported.Name()
		if importedName == name {
			return imported
		}
	}
	return nil
}

func typeEqual(a, b types.Type) bool {
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
