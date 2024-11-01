package main

import (
	"fmt"
	"os"

	"github.com/golangci/golangci-lint/pkg/commands"
)

func main() {
	os.Args = []string{"", "custom"} // dummy args for building custom linter
	if err := commands.Execute(commands.BuildInfo{}); err != nil {
		fmt.Printf("error: %v", err)
	}
}
