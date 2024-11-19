package process

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

type TargetDefinition struct {
	Name     string
	TypeSpec ast.TypeSpec
}

func FindTargetDefinition(dirPath string) (TargetDefinition, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirPath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	var targetDefinition *TargetDefinition = nil
	for _, pkg := range pkgs {
		if targetDefinition != nil {
			break
		}
		pkgName := pkg.Name
		for fname, f := range pkg.Files {
			if targetDefinition != nil {
				break
			}
			if strings.HasSuffix(fname, "_test.go") {
				continue
			}

			ast.Inspect(f, func(n ast.Node) bool {
				genDecl, ok := n.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					return true
				}

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					if typeSpec.Name.Name == pkgName {
						targetDefinition = &TargetDefinition{
							Name:     pkgName,
							TypeSpec: *typeSpec,
						}
						return false
					}
				}
				return true
			})
		}
	}

	if targetDefinition == nil {
		return TargetDefinition{}, fmt.Errorf("target definition not found") //ignore:fields_require
	}

	return *targetDefinition, nil
}
