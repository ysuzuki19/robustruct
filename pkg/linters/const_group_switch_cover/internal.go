package const_group_switch_cover

import (
	"go/types"

	"golang.org/x/tools/go/analysis"

	"github.com/ysuzuki19/robustruct/internal/logger"
)

func findConstsInternaly(pass *analysis.Pass, namedType *types.Named) (consts []*types.Const) {
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
		logger.Debug("No matching type found for:", name)
		return nil
	}

	for _, obj := range info.Defs {
		if obj == nil {
			continue
		}

		if c, ok := obj.(*types.Const); ok && types.Identical(c.Type(), tagType) {
			consts = append(consts, c)
		}
	}
	if len(consts) == 0 {
		logger.Debug("No constants found for type:", tagType)
		return nil
	}
	logger.Debug("Constants found:", len(consts))

	return consts
}
