package process

import (
	"fmt"
	"strings"

	"github.com/ysuzuki19/robustruct/cmd/exp/senumgen/internal/process/tmpl"
)

type GenerateArgs struct {
	DirPath       string
	Name          string
	AnalyzeResult AnalyzeResult
}

func Generate(args GenerateArgs) ([]byte, error) {
	var defTypeParams string
	for idx, tp := range args.AnalyzeResult.TypeParams {
		if idx != 0 {
			defTypeParams += ", "
		}
		defTypeParams += fmt.Sprintf("%s %s", tp.Name, tp.TypeName)
	}

	var useTypeParams string
	for idx, tp := range args.AnalyzeResult.TypeParams {
		if idx != 0 {
			useTypeParams += ", "
		}
		useTypeParams += tp.Name
	}

	templateData := struct {
		Package       string
		Name          string
		DefTypeParams string
		UseTypeParams string
		EnumDefName   string
		EnumUseName   string
		Variants      []Variant
	}{
		Package:       args.Name,
		Name:          args.Name,
		DefTypeParams: defTypeParams,
		UseTypeParams: useTypeParams,
		EnumDefName:   fmt.Sprintf("%sEnum%s", strings.ToLower(args.Name), tmpl.Bracket(defTypeParams)),
		EnumUseName:   fmt.Sprintf("%sEnum%s", strings.ToLower(args.Name), tmpl.Bracket(useTypeParams)),
		Variants:      args.AnalyzeResult.Variants,
	}

	tc := tmpl.NewTemplateCollector()

	tc.Merge(tmpl.Header, templateData).
		Merge(tmpl.Tag, templateData).
		Merge(tmpl.Enum, templateData).
		Merge(tmpl.New, templateData).
		Merge(tmpl.Is, templateData).
		Merge(tmpl.As, templateData).
		Merge(tmpl.Switch, templateData).
		Merge(tmpl.Match, templateData)

	generated, err := tc.Export()
	if err != nil {
		return nil, err
	}

	return generated, nil
}
