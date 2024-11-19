package process

import (
	"fmt"
	"go/ast"
)

type Variant struct {
	Name      string
	FieldName string
	TypeName  string
	HasData   bool
}

type TypeParam struct {
	Name     string
	TypeName string
}

type AnalyzeResult struct {
	TypeParams []TypeParam
	Variants   []Variant
}

func Analyze(targetDefinition TargetDefinition) (AnalyzeResult, error) {
	analyzeResult := AnalyzeResult{} // ignore:fields_require

	name := targetDefinition.Name
	typeSpec := targetDefinition.TypeSpec

	var typeParams []TypeParam
	var variants []Variant

	if typeSpec.TypeParams != nil {
		for _, typeParam := range typeSpec.TypeParams.List {
			tp := TypeParam{
				Name:     typeParam.Names[0].Name,
				TypeName: fmt.Sprintf("%s", typeParam.Type),
			}
			typeParams = append(typeParams, tp)
		}
	}

	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return analyzeResult, fmt.Errorf("%s is not a struct type", name)
	}

	for _, field := range structType.Fields.List {
		fieldName := field.Names[0].Obj.Name
		fieldType := field.Type
		hasData := true
		var typeName string
		switch typ := fieldType.(type) {
		case *ast.SelectorExpr:
			if xName, ok := typ.X.(*ast.Ident); ok {
				if xName.Name == "senum" && typ.Sel.Name == "NonVar" {
					hasData = false
				}
				typeName = fmt.Sprintf("%s.%s", xName.Name, typ.Sel.Name)
			} else {
				typeName = typ.Sel.Name
			}
		case *ast.StarExpr:
			switch typ := typ.X.(type) {
			case *ast.Ident:
				typeName = fmt.Sprintf("*%s", typ.Name)
			case *ast.SelectorExpr:
				if xName, ok := typ.X.(*ast.Ident); ok {
					typeName = fmt.Sprintf("*%s.%s", xName.Name, typ.Sel.Name)
				} else {
					return analyzeResult, fmt.Errorf("not implemented")
				}
			}
		case *ast.Ident:
			typeName = typ.Name
		default:
			return analyzeResult, fmt.Errorf("not implemented Generics")
		}

		variants = append(variants, Variant{
			Name:      fieldName,
			FieldName: fieldName,
			TypeName:  typeName,
			HasData:   hasData,
		})
	}

	return AnalyzeResult{
		TypeParams: typeParams,
		Variants:   variants,
	}, nil
}
