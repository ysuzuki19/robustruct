package process

import (
	"fmt"
	"strings"

	"github.com/ysuzuki19/robustruct/pkg/option"
)

type testDocOpening struct {
	Index         int
	StructureName option.Option[string]
	FuncName      string
}

type TestDoc struct {
	StructName option.Option[string]
	FuncName   string
	Content    string
}

func ParseTestDocs(test string) ([]TestDoc, error) {
	lines := strings.Split(test, "\n")
	var tds []TestDoc
	opened := option.None[testDocOpening]()

	for idx, line := range lines {
		if rest, ok := matchAndStrip(tdRegex, line); ok {
			if rest, ok := matchAndStrip(tdBeginRegex, rest); ok {
				if opened.IsSome() {
					return nil, fmt.Errorf("testdoc begin found but already opened at line %v", opened.Ptr().Index)
				}
				trimed := strings.TrimSpace(rest)
				parts := strings.Split(trimed, ".")
				switch len(parts) {
				case 1:
					opened = option.NewSome(
						testDocOpening{
							Index:         idx,
							StructureName: option.None[string](),
							FuncName:      parts[0],
						})
				case 2:
					opened = option.NewSome(
						testDocOpening{
							Index:         idx,
							StructureName: option.Some(&parts[0]),
							FuncName:      parts[1],
						})
				default:
					return nil, fmt.Errorf("testdoc begin line must contain either 'begin StructName' or 'begin StructName.FuncName'")
				}
			}
			if _, ok := matchAndStrip(tdEndRegex, rest); ok {
				if begin, ok := opened.Take().Get(); ok {
					tds = append(tds, TestDoc{
						StructName: begin.StructureName,
						FuncName:   begin.FuncName,
						Content:    strings.Join(lines[begin.Index+1:idx], "\n"),
					})
				} else {
					return nil, fmt.Errorf("testdoc end found but not opened")
				}
			}
		}
	}

	return tds, nil
}
