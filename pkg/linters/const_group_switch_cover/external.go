package const_group_switch_cover

import (
	"go/types"

	"golang.org/x/tools/go/analysis"

	"github.com/ysuzuki19/robustruct/internal/logger"
)

func runExternal(pass *analysis.Pass, namedType *types.Named) (consts []*types.Const) {
	typeName := namedType.Obj().Name()
	pkgName := namedType.Obj().Pkg().Name()

	tagTypePkg := findImport(pass, pkgName)
	if tagTypePkg == nil {
		logger.Debug("No import found for package:", pkgName)
		return nil
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
			if typeEqual(c.Type(), tagType) {
				logger.Debug("Matching constant found:", c.Type().String())
				consts = append(consts, c)
			}
		}
	}

	if len(consts) == 0 {
		logger.Debug("No constants found for type:", tagType)
		return nil
	}
	logger.Debug("Constants found:", len(consts))

	return consts
}
