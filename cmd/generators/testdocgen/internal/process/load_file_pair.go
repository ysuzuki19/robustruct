package process

import (
	"os"

	"github.com/ysuzuki19/robustruct/cmd/generators/testdocgen/internal/strchain"
)

func LoadFilePair(codePath string) (source string, test string, err error) {
	b, err := os.ReadFile(codePath)
	if err != nil {
		return
	}
	source = string(b)

	testPath := strchain.From(codePath).Replace(".go", "_test.go", 1).String()
	b, err = os.ReadFile(testPath)
	if err != nil {
		return
	}
	test = string(b)
	return
}
