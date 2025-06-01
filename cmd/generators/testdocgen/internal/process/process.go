package process

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"regexp"
	"sort"

	"github.com/ysuzuki19/robustruct/cmd/generators/internal/postgenerate"
	"github.com/ysuzuki19/robustruct/cmd/generators/internal/writer"
	"github.com/ysuzuki19/robustruct/cmd/generators/testdocgen/internal/strchain"
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

	testPath := strchain.From(codePath).Replace(".go", "_test.go", 1).String()
	b, err = os.ReadFile(testPath)
	if err != nil {
		return
	}
	test = string(b)
	return
}

type Plan struct {
	InsertIndex  int
	ReplaceCount int
	Lines        []string
}

func FindExamplePosition(fset *token.FileSet, fn *ast.FuncDecl) (int, int, error) {
	if fn.Doc == nil || len(fn.Doc.List) == 0 {
		return fset.Position(fn.Pos()).Line - 1, 0, nil
	}
	begin := option.None[int]()
	var searched int
	for _, comment := range fn.Doc.List {
		searched = fset.Position(comment.Pos()).Line
		if begin.IsSome() {
			// if exampleAnnotation.IsSome() {
			// if comment.Text == "//" {
			// 	if begin, ok := begin.Get(); ok {
			// 		return *begin, searched, nil
			// 	}
			// } else {
			// 	if begin.IsNone() {
			// 		pos := fset.Position(comment.Pos())
			// 		begin = option.NewSome(pos.Line)
			// 	}
			// }
		} else if regexp.MustCompile(`^\s*//\s*Example:?.*`).MatchString(comment.Text) {
			begin = option.NewSome(searched - 1)
		}
	}
	if exm, ok := begin.Get(); ok {
		return *exm, searched - *exm, nil
	}
	// if `Example:` is not found, we return the last searched line as the end
	return searched, 0, nil
}

func ApplyGoDoc(source string, plans []Plan) string {
	sort.Slice(plans, func(i, j int) bool {
		return plans[i].InsertIndex > plans[j].InsertIndex
	})

	lines := strchain.From(source).Split("\n")
	for _, plan := range plans {
		// fmt.Println("Applying plan:", plan.InsertIndex, plan.ReplaceCount)
		lines = lines.Splice(plan.InsertIndex, plan.ReplaceCount, plan.Lines)
	}

	return lines.Join("\n").String()
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

	fmt.Println("TestDoc Count:", len(tds))
	fmt.Println("Plan Count:", len(plans))

	// fmt.Println(source)
	updated := ApplyGoDoc(source, plans)
	if updated == source {
		return nil
	}
	// fmt.Println(updated)

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
