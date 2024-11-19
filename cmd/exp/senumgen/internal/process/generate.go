package process

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"log"
	"strings"

	"golang.org/x/tools/imports"
)

type TemplateData struct {
	Package       string
	Name          string
	DefTypeParams string
	UseTypeParams string
	EnumDefName   string
	EnumUseName   string
	Variants      []Variant
}

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
	templateData := TemplateData{
		Package:       args.Name,
		Name:          args.Name,
		DefTypeParams: defTypeParams,
		UseTypeParams: useTypeParams,
		EnumDefName:   fmt.Sprintf("%sEnum%s", strings.ToLower(args.Name), bracket(defTypeParams)),
		EnumUseName:   fmt.Sprintf("%sEnum%s", strings.ToLower(args.Name), bracket(useTypeParams)),
		Variants:      args.AnalyzeResult.Variants,
	}

	tmplBytes, err := structEnumTemplateFS.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err := template.New(fileName).Funcs(template.FuncMap{
		"capitalize": capitalize,
		"bracket":    bracket,
		"csvConnect": csvConnect,
	}).Parse(string(tmplBytes))
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, templateData); err != nil {
		log.Fatal(err)
	}

	return buf.Bytes(), nil
}

type PostGenerateArgs struct {
	OutputFilePath string
	Buf            []byte
}

func PostGenerate(args PostGenerateArgs) ([]byte, error) {
	formattedCode, err := format.Source(args.Buf)
	if err != nil {
		log.Fatal(err)
	}

	output, err := imports.Process(args.OutputFilePath, formattedCode, nil)
	if err != nil {
		log.Fatal(err)
	}

	return output, nil
}
