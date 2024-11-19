package process

import (
	"go/format"
	"log"

	"golang.org/x/tools/imports"
)

type PostGenerateArgs struct {
	OutputFilePath string
	Buf            []byte
}

func PostGenerate(args PostGenerateArgs) ([]byte, error) {
	formattedCode, err := format.Source(args.Buf)
	if err != nil {
		log.Fatal(err)
	}

	output, err := imports.Process(args.OutputFilePath, formattedCode, nil)
	if err != nil {
		log.Fatal(err)
	}

	return output, nil
}
