package main

import (
	"flag"
	"log"
	"os"

	"github.com/ysuzuki19/robustruct/cmd/exp/senumgen/internal"
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

	internal.Process(internal.Args{
		DirPath: *dirPath,
		Writer: &internal.FileWriter{
			FilePath: internal.OutputFilePath(*dirPath),
		},
	})
}
