package postgenerate

import (
	"go/format"
	"log"
)

type PostGenerateArgs struct {
	Buf []byte
}

func PostGenerate(args PostGenerateArgs) ([]byte, error) {
	formattedCode, err := format.Source(args.Buf)
	if err != nil {
		log.Fatal(err)
	}

	return formattedCode, nil
}

//go:generate go run ../../testdocgen/main.go -- -file=$GOFILE
