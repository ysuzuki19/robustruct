package main

import (
	"fmt"
	"os"

	"github.com/golangci/golangci-lint/pkg/commands"
)

func main() {
	os.Args = []string{"", "custom", "-vv"} // dummy args for building custom linter
	if err := commands.Execute(commands.BuildInfo{
		GoVersion: "",
		Version:   "",
		Commit:    "",
		Date:      "",
	}); err != nil {
		fmt.Printf("error: %v", err)
	}
}
