package main

import (
	"flag"
	"log"
	"os"

	"github.com/ysuzuki19/robustruct/cmd/exp/internal/writer"
	"github.com/ysuzuki19/robustruct/cmd/exp/senumgen/internal/process"
)

func main() {
	dirPath := flag.String("dir", "", "input dir path")
	flag.Parse()

	if *dirPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		dirPath = &wd
	}

	if err := process.Process(process.Args{
		DirPath: *dirPath,
		Writer: &writer.FileWriter{
			FilePath: process.OutputFilePath(*dirPath),
		},
	}); err != nil {
		log.Fatal(err)
	}
}
