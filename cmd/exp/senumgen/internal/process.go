package internal

import (
	"bytes"
	"embed"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"log"
	"strings"

	"golang.org/x/tools/imports"
)

const fileName = "templates/senum.go.tmpl"

//go:embed templates/senum.go.tmpl
var structEnumTemplateFS embed.FS

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

type TemplateData struct {
	Package       string
	Name          string
	DefTypeParams string
	UseTypeParams string
	EnumDefName   string
	EnumUseName   string
	Variants      []Variant
}

type Args struct {
	DirPath string
	Writer  Write
}

func Process(args Args) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, args.DirPath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	var name string
	var typeParams []TypeParam
	var variants []Variant
	for _, pkg := range pkgs {
		for fname, f := range pkg.Files {
			if strings.HasSuffix(fname, "_test.go") {
				continue
			}
			name = pkg.Name

			ast.Inspect(f, func(n ast.Node) bool {
				genDecl, ok := n.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					return true
				}

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok || typeSpec.Name.Name != name {
						continue
					}

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
						continue
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
									panic("not implemented")
								}
							}
						case *ast.Ident:
							typeName = typ.Name
						default:
							panic("not implemented Generics")
						}

						variants = append(variants, Variant{
							Name:      fieldName,
							FieldName: fieldName,
							TypeName:  typeName,
							HasData:   hasData,
						})
					}
				}
				return true
			})
		}
	}

	var defTypeParams string
	for idx, tp := range typeParams {
		if idx != 0 {
			defTypeParams += ", "
		}
		defTypeParams += fmt.Sprintf("%s %s", tp.Name, tp.TypeName)
	}

	var useTypeParams string
	for idx, tp := range typeParams {
		if idx != 0 {
			useTypeParams += ", "
		}
		useTypeParams += tp.Name
	}

	templateData := TemplateData{
		Package:       name,
		Name:          name,
		DefTypeParams: defTypeParams,
		UseTypeParams: useTypeParams,
		EnumDefName:   fmt.Sprintf("%sEnum%s", strings.ToLower(name), bracket(defTypeParams)),
		EnumUseName:   fmt.Sprintf("%sEnum%s", strings.ToLower(name), bracket(useTypeParams)),
		Variants:      variants,
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

	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	output, err := imports.Process(OutputFilePath(args.DirPath), formattedCode, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := args.Writer.Write(output); err != nil {
		log.Fatal(err)
	}
}
