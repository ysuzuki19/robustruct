package engine

import (
	"fmt"

	"github.com/ysuzuki19/robustruct/cmd/generators/internal/postgenerate"
	"github.com/ysuzuki19/robustruct/cmd/generators/internal/writer"
)

type Args struct {
	CodePath string
	Writer   writer.Writer
}

func Run(args Args) error {
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
