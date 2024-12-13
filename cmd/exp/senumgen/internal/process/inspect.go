package process

import (
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

type Definition struct {
	Name         string
	CommentLines []string
}

func FindDefinition(dirPath string) (Definition, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dirPath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	var pkgName string
	var commentLines []string
	for _, pkg := range pkgs {
		pkgName = pkg.Name
		for fname, f := range pkg.Files {
			fmt.Println(fname)
			if !strings.HasSuffix(fname, fmt.Sprintf("%s.go", pkgName)) {
				continue
			}

			for _, commentGroup := range f.Comments {
				comment := commentGroup.Text()
				splited := strings.Split(comment, "\n")
				commentLines = append(commentLines, splited...)
			}
		}
	}

	if len(commentLines) == 0 {
		return Definition{}, fmt.Errorf("definition comment not found") // ignore:fields_require
	}

	return Definition{
		Name:         pkgName,
		CommentLines: commentLines,
	}, nil
}
