package process

import "fmt"

func OutputFilePath(dirPath string) string {
	return fmt.Sprintf("%s/senum.gen.go", dirPath)
}
