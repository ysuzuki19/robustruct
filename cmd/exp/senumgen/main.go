package main

import (
	"flag"
	"log"
	"os"

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

	process.Process(process.Args{
		DirPath: *dirPath,
		Writer: &process.FileWriter{
			FilePath: process.OutputFilePath(*dirPath),
		},
	})
}
