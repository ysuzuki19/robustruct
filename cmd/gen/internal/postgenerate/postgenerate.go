package postgenerate

import (
	"go/format"
	"log"
)

type PostGenerateArgs struct {
	Buf []byte
}

// Example:
//
//	buf := []byte("package main\nfunc main() {\nprintln(\"testing\")}")
//	formattedCode, err := postgenerate.PostGenerate(
//		postgenerate.PostGenerateArgs{
//			Buf: buf,
//		},
//	)
func PostGenerate(args PostGenerateArgs) ([]byte, error) {
	formattedCode, err := format.Source(args.Buf)
	if err != nil {
		log.Fatal(err)
	}

	return formattedCode, nil
}

//go:generate go run github.com/ysuzuki19/robustruct/cmd/gen/testdocgen -file=$GOFILE
