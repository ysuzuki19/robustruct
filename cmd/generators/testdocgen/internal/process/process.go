package process

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/ysuzuki19/robustruct/cmd/generators/internal/writer"
	"github.com/ysuzuki19/robustruct/pkg/option"
)

var tdRegex = regexp.MustCompile(`^\s*//\s*testdoc\s+`)
var tdBeginRegex = regexp.MustCompile(`^\s*begin\s+`)
var tdEndRegex = regexp.MustCompile(`^\s*end$`)

type TestDocOpening struct {
	Index         int
	StructureName option.Option[string]
	FuncName      string
}

type TestDoc struct {
	StructName option.Option[string]
	FuncName   string
	Content    string
}

type Args struct {
	CodePath string
	TestPath string
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

func ParseTestDocs(test string) ([]TestDoc, error) {
	lines := strings.Split(test, "\n")
	var tds []TestDoc
	opened := option.None[TestDocOpening]()

	for idx, line := range lines {
		if rest, ok := matchAndStrip(tdRegex, line); ok {
			if rest, ok := matchAndStrip(tdBeginRegex, rest); ok {
				if opened.IsSome() {
					return nil, fmt.Errorf("testdoc begin found but already opened at line %v", *opened.Ptr())
				}
				trimed := strings.TrimSpace(rest)
				parts := strings.Split(trimed, ".")
				switch len(parts) {
				case 1:
					opened = option.NewSome(
						TestDocOpening{
							Index:         idx,
							StructureName: option.None[string](),
							FuncName:      parts[0],
						})
				case 2:
					opened = option.NewSome(
						TestDocOpening{
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
		return *exm, searched + 1, nil
	}
	return 0, 0, fmt.Errorf("Example not found in doc comments")
}

func RecvTypeName(fn *ast.FuncDecl) option.Option[string] {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return option.None[string]()
	}
	switch t := fn.Recv.List[0].Type.(type) {
	case *ast.Ident:
		return option.Some(&t.Name)
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return option.Some(&ident.Name)
		}
	case *ast.IndexExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return option.Some(&ident.Name)
		}
	case *ast.IndexListExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return option.Some(&ident.Name)
		}
	}
	return option.None[string]()
}

func PlanGoDoc(source string, tds []TestDoc) ([]Plan, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", source, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	plans := []Plan{}

	for _, td := range tds {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			if structName, ok := td.StructName.Get(); ok {
				recvTypeName, ok := RecvTypeName(fn).Get()
				if !ok {
					continue
				}

				fnName := fn.Name.Name
				if recvTypeName == structName && fnName == td.FuncName {
					begin, end, err := FindExampleRange(fset, fn.Doc.List)
					if err != nil {
						return nil, fmt.Errorf("failed to find example range: %w", err)
					}
					//TODO use strchain
					// lines := strchain.String(td.Content).Split("\n").Map(func(line string) string {
					// 	return "// " + line
					// })
					lines := strings.Split(td.Content, "\n")
					for i := range lines {
						lines[i] = "// " + lines[i]
					}
					plans = append(plans, Plan{
						Begin: begin,
						End:   end,
						Lines: lines,
					})
				}
			} else {
				if fn.Name.Name == td.FuncName {
					begin, end, err := FindExampleRange(fset, fn.Doc.List)
					if err != nil {
						return nil, fmt.Errorf("failed to find example range: %w", err)
					}
					lines := strings.Split(td.Content, "\n")
					for i := range lines {
						lines[i] = "// " + lines[i]
					}
					plans = append(plans, Plan{
						Begin: begin,
						End:   end,
						Lines: strings.Split(td.Content, "\n"),
					})
				}
			}
		}
	}
	return plans, nil
}

func ApplyGoDoc(source string, plans []Plan) string {
	sort.Slice(plans, func(i, j int) bool {
		return plans[i].Begin > plans[j].Begin
	})

	lines := strings.Split(source, "\n")
	for _, plan := range plans {
		fmt.Println("Applying plan:", plan.Begin, plan.End, plan.Lines)
		lines = append(lines[:plan.Begin-1], append(plan.Lines, lines[plan.End-1:]...)...)
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

	fmt.Println("Updating source code...")
	fmt.Println(updated)

	// formatted, err := postgenerate.PostGenerate(
	// 	postgenerate.PostGenerateArgs{
	// 		Buf: []byte(updated),
	// 	},
	// )
	// if err != nil {
	// 	return fmt.Errorf("failed to format updated source: %w", err)
	// }

	// err = args.Writer.Write(formatted)
	// if err != nil {
	// 	return fmt.Errorf("failed to write updated source: %w", err)
	// }

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
