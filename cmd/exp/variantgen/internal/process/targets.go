package process

import "fmt"

func OutputFilePath(dirPath string) string {
	return fmt.Sprintf("%s/variant.gen.go", dirPath)
}
