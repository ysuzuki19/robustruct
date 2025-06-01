package process

import (
	"fmt"
	"regexp"

	"github.com/ysuzuki19/robustruct/cmd/generators/testdocgen/internal/strchain"
	"github.com/ysuzuki19/robustruct/pkg/option"
)

var tdRegex = regexp.MustCompile(`^\s*//\s*testdoc\s+`)
var tdBeginRegex = regexp.MustCompile(`^\s*begin\s+`)
var tdEndRegex = regexp.MustCompile(`^\s*end$`)

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

	for idx, line := range lines.Entries() {
		if rest, ok := line.MatchAndStrip(tdRegex); ok {
			if rest, ok := rest.MatchAndStrip(tdBeginRegex); ok {
				if opened.IsSome() {
					return nil, fmt.Errorf("testdoc begin found but already opened at line %v", opened.Ptr().LineNo)
				}
				parts := rest.TrimSpace().Split(".").Collect()
				switch len(parts) {
				case 1:
					if parts[0] == "" {
						return nil, fmt.Errorf("testdoc begin line must contain either 'begin StructName' or 'begin StructName.FuncName'")
					}
					opened = option.NewSome(
						testDocOpening{
							LineNo:        idx,
							StructureName: option.None[string](),
							FuncName:      parts[0],
						})
				case 2:
					if parts[0] == "" || parts[1] == "" {
						return nil, fmt.Errorf("testdoc begin line must contain either 'begin StructName' or 'begin StructName.FuncName'")
					}
					opened = option.NewSome(
						testDocOpening{
							LineNo:        idx,
							StructureName: option.Some(&parts[0]),
							FuncName:      parts[1],
						})
				default:
					return nil, fmt.Errorf("testdoc begin line must contain either 'begin StructName' or 'begin StructName.FuncName'")
				}
			} else if _, ok := rest.MatchAndStrip(tdEndRegex); ok {
				if begin, ok := opened.Take().Get(); ok {
					tds = append(tds, TestDoc{
						StructName: begin.StructureName,
						FuncName:   begin.FuncName,
						Content:    lines.Slice(begin.LineNo+1, idx).Join("\n").String(),
					})
				} else {
					return nil, fmt.Errorf("testdoc end found but not opened")
				}
			} else {
				return nil, fmt.Errorf("testdoc line must start with '// testdoc begin' or '// testdoc end': %s", line)
			}
		}
	}

	return tds, nil
}
