package process

import (
	"fmt"

	"github.com/ysuzuki19/robustruct/cmd/generators/testdocgen/internal/strchain"
	"github.com/ysuzuki19/robustruct/pkg/option"
)

type testDocOpening struct {
	LineNo        int
	StructureName option.Option[string]
	FuncName      string
}

type TestDoc struct {
	StructName option.Option[string]
	FuncName   string
	Content    string
}

func ParseTestDocs(test string) ([]TestDoc, error) {
	lines := strchain.From(test).Split("\n")
	var tds []TestDoc
	opened := option.None[testDocOpening]()

	for idx, line := range lines.Collect() {
		if rest, ok := matchAndStrip(tdRegex, line); ok {
			if rest, ok := matchAndStrip(tdBeginRegex, rest); ok {
				if opened.IsSome() {
					return nil, fmt.Errorf("testdoc begin found but already opened at line %v", opened.Ptr().LineNo)
				}
				parts := strchain.From(rest).TrimSpace().Split(".").Collect()
				switch len(parts) {
				case 1:
					opened = option.NewSome(
						testDocOpening{
							LineNo:        idx,
							StructureName: option.None[string](),
							FuncName:      parts[0],
						})
				case 2:
					opened = option.NewSome(
						testDocOpening{
							LineNo:        idx,
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
						Content:    lines.Slice(begin.LineNo+1, idx).Join("\n").String(),
					})
				} else {
					return nil, fmt.Errorf("testdoc end found but not opened")
				}
			}
		}
	}

	return tds, nil
}
