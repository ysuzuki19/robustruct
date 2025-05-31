package process

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/ysuzuki19/robustruct/cmd/generators/internal/postgenerate"
	"github.com/ysuzuki19/robustruct/cmd/generators/internal/writer"
	"github.com/ysuzuki19/robustruct/pkg/option"
)

var tdRegex = regexp.MustCompile(`^\s*//\s*testdoc\s+`)
var tdBeginRegex = regexp.MustCompile(`^\s*begin\s+`)
var tdEndRegex = regexp.MustCompile(`^\s*end$`)

type Args struct {
	CodePath string
	Writer   writer.Writer
}

func LoadFilePair(codePath string) (source string, test string, err error) {
	b, err := os.ReadFile(codePath)
	if err != nil {
		return
	}
	source = string(b)

	testPath := strings.Replace(codePath, ".go", "_test.go", 1)
	b, err = os.ReadFile(testPath)
	if err != nil {
		return
	}
	test = string(b)
	return
}

type Plan struct {
	Begin int
	End   int
	Lines []string
}

func FindExampleRange(fset *token.FileSet, docList []*ast.Comment) (int, int, error) {
	exampleAnnotation := option.None[int]()
	begin := option.None[int]()
	var searched int
	for _, comment := range docList {
		searched = fset.Position(comment.Pos()).Line
		if exampleAnnotation.IsSome() {
			if comment.Text == "//" {
				if begin, ok := begin.Get(); ok {
					return *begin, searched, nil
				}
			} else {
				if begin.IsNone() {
					pos := fset.Position(comment.Pos())
					begin = option.NewSome(pos.Line)
				}
			}
		} else if regexp.MustCompile(`^\s*//\s*Example:?.*`).MatchString(comment.Text) {
			if exampleAnnotation.IsSome() {
				return 0, 0, fmt.Errorf("Nested Example Detected %d", *begin.Ptr())
			}
			exampleAnnotation = option.NewSome(searched)
		}
	}
	if exm, ok := exampleAnnotation.Get(); ok {
		if *exm == searched {
			return *exm, searched, nil
		}
		return *exm, searched, nil
	}
	return 0, 0, fmt.Errorf("Example not found in doc comments")
}

func ApplyGoDoc(source string, plans []Plan) string {
	sort.Slice(plans, func(i, j int) bool {
		return plans[i].Begin > plans[j].Begin
	})

	lines := strings.Split(source, "\n")
	for _, plan := range plans {
		fmt.Println("Applying plan:", plan.Begin, plan.End, plan.Lines)
		above := lines[:plan.Begin]
		below := lines[plan.End:]
		lines = append(above, append(plan.Lines, below...)...)
	}

	return strings.Join(lines, "\n")
}

func Process(args Args) error {
	source, test, err := LoadFilePair(args.CodePath)
	if err != nil {
		return fmt.Errorf("failed to load file pair: %w", err)
	}

	tds, err := ParseTestDocs(test)
	if err != nil {
		return fmt.Errorf("failed to parse test docs: %w", err)
	}

	plans, err := PlanGoDoc(source, tds)
	if err != nil {
		return fmt.Errorf("failed to update Go doc: %w", err)
	}

	updated := ApplyGoDoc(source, plans)
	if updated == source {
		return nil
	}

	formatted, err := postgenerate.PostGenerate(
		postgenerate.PostGenerateArgs{
			Buf: []byte(updated),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to format updated source: %w", err)
	}

	err = args.Writer.Write([]byte(formatted))
	if err != nil {
		return fmt.Errorf("failed to write updated source: %w", err)
	}

	return nil
}

func matchAndStrip(re *regexp.Regexp, input string) (string, bool) {
	loc := re.FindStringIndex(input)
	if loc == nil {
		return input, false
	}
	return input[:loc[0]] + input[loc[1]:], true
}

// go run ../../cmd/generators/testdocgen/main.go -- -file=$GOFILE
