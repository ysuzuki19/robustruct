package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	once   sync.Once
	result string
)

func findGoMod(absPath string) string {
	if absPath == "" || absPath == "/" {
		return ""
	}
	dir := filepath.Dir(absPath)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		if entry.Name() == "go.mod" {
			return absPath + "/"
		}
	}
	return findGoMod(dir)
}

func goModPath(callerPath string) string {
	once.Do(func() {
		result = findGoMod(callerPath)
	})
	return result
}

func Debug(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Println("Debug: Unable to get caller information")
		return
	}
	prunePrefix := goModPath(file)
	pruned := strings.Replace(file, prunePrefix, "", 1)
	fmt.Printf("[%s:%d] ", pruned, line)
	fmt.Println(args...)
}
