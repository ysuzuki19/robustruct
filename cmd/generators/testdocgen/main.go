package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ysuzuki19/robustruct/cmd/generators/internal/writer"
	"github.com/ysuzuki19/robustruct/cmd/generators/testdocgen/internal/engine"
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

	if err := engine.Run(engine.Args{
		CodePath: codePath,
		// Writer: &writer.MemoryWriter{
		// 	Content: "",
		// },
		Writer: &writer.FileWriter{
			FilePath: codePath,
		},
	}); err != nil {
		log.Fatal(err)
	}
}
