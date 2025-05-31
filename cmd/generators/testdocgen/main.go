package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ysuzuki19/robustruct/cmd/generators/internal/writer"
	"github.com/ysuzuki19/robustruct/cmd/generators/testdocgen/internal/process"
)

type Args struct {
	File string `clip:"file"`
}

func LoadArgs() (*Args, error) {
	var args Args
	for _, arg := range os.Args {
		parts := strings.Split(arg, "=")
		if len(parts) != 2 {
			continue
		}
		switch parts[0] {
		case "--file", "-file":
			args.File = parts[1]
		default:
			log.Fatalf("unknown argument: %s", parts[0])
		}
	}

	if args.File == "" {
		log.Fatal("file argument is required")
		return nil, fmt.Errorf("file argument is required")
	}

	return &args, nil
}

func main() {
	args, err := LoadArgs()
	if err != nil {
		log.Fatal(err)
	}

	codePath, err := filepath.Abs(args.File)
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}
	fmt.Println("Absolute path:", codePath)

	if err := process.Process(process.Args{
		CodePath: codePath,
		Writer: &writer.MemoryWriter{
			Content: "",
		},
		// Writer: &writer.FileWriter{
		// 	FilePath: process.OutputFilePath(*dirPath),
		// },
	}); err != nil {
		log.Fatal(err)
	}
}
