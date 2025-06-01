package engine

import (
	"fmt"

	"github.com/ysuzuki19/robustruct/cmd/gen/internal/postgenerate"
	"github.com/ysuzuki19/robustruct/cmd/gen/internal/writer"
)

type Args struct {
	CodePath string
	Writer   writer.Writer
}

// Run executes the testdocgen engine.
// 1. Loads the source and test files.
// 2. Parses the testdoc annotations from the test file.
// 3. Plans the updates to the Go doc based on the parsed testdoc annotations.
// 4. Format the updated source code.
// 5. Writes the formatted source code to the specified writer.
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
