package process

import (
	"fmt"

	"github.com/ysuzuki19/robustruct/cmd/exp/internal/postgenerate"
	"github.com/ysuzuki19/robustruct/cmd/exp/internal/writer"
)

type Args struct {
	DirPath string
	Writer  writer.Writer
}

/**
 * Process generates the enum code from the given directory path.
 * 1. Inspect the AST of the given directory path.
 * 2. Parse the AST and extract the type parameters and variants.
 * 3. Generate the enum code using the extracted type parameters and variants.
 * 4. Format and resolve imports the generated code.
 * 5. Write the generated code to the output file.
 */
func Process(args Args) error {
	definitionComment, err := FindDefinition(args.DirPath)
	if err != nil {
		return err
	}

	analyzeResult, err := Analyze(definitionComment.CommentLines)
	if err != nil {
		return err
	}

	generated, err := Generate(GenerateArgs{
		DirPath:       args.DirPath,
		PackageName:   definitionComment.Name,
		AnalyzeResult: analyzeResult,
	})
	if err != nil {
		return err
	}

	code := string(generated)
	fmt.Println(code)

	output, err := postgenerate.PostGenerate(postgenerate.PostGenerateArgs{
		OutputFilePath: OutputFilePath(args.DirPath),
		Buf:            generated,
	})
	if err != nil {
		return err
	}

	if err := args.Writer.Write(output); err != nil {
		return err
	}

	return nil
}
