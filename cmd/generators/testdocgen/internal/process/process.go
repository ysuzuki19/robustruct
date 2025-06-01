package process

import (
	"fmt"
	"os"

	"github.com/ysuzuki19/robustruct/cmd/generators/internal/postgenerate"
	"github.com/ysuzuki19/robustruct/cmd/generators/internal/writer"
	"github.com/ysuzuki19/robustruct/cmd/generators/testdocgen/internal/strchain"
)

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
