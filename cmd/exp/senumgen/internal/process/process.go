package process

import (
	"embed"
	"log"
)

const fileName = "templates/senum.go.tmpl"

//go:embed templates/senum.go.tmpl
var structEnumTemplateFS embed.FS

type Args struct {
	DirPath string
	Writer  Write
}

/**
 * Process generates the enum code from the given directory path.
 * 1. Inspect the AST of the given directory path.
 * 2. Parse the AST and extract the type parameters and variants.
 * 3. Generate the enum code using the extracted type parameters and variants.
 * 4. Format and resolve imports the generated code.
 * 5. Write the generated code to the output file.
 */
func Process(args Args) {
	targetDefinition, err := FindTargetDefinition(args.DirPath)
	if err != nil {
		log.Fatal(err)
	}

	name := targetDefinition.Name

	analyzeResult, err := Analyze(targetDefinition)
	if err != nil {
		log.Fatal(err)
	}

	generated, err := Generate(GenerateArgs{
		DirPath:       args.DirPath,
		Name:          name,
		AnalyzeResult: analyzeResult,
	})
	if err != nil {
		log.Fatal(err)
	}

	output, err := PostGenerate(PostGenerateArgs{
		OutputFilePath: OutputFilePath(args.DirPath),
		Buf:            generated,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := args.Writer.Write(output); err != nil {
		log.Fatal(err)
	}
}
