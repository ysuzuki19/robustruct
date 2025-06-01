package writer

import (
	"os"
)

type Writer interface {
	Write(buf []byte) error
}

var _ Writer = &MemoryWriter{} //ignore:fields_require
var _ Writer = &FileWriter{}   //ignore:fields_require

type MemoryWriter struct {
	Content string
}

// Example:
//
//	w := &writer.MemoryWriter{}
//	err := w.Write([]byte("Hello, World!"))
//	require.NoError(err)
//	require.Equal("Hello, World!", w.Content)
func (w *MemoryWriter) Write(buf []byte) error {
	w.Content = string(buf)
	return nil
}

type FileWriter struct {
	FilePath string
}

// Example:
//
//	w := &writer.MemoryWriter{}
//	err := w.Write([]byte("Hello, World!"))
//	require.NoError(err)
//	require.Equal("Hello, World!", w.Content)
func (w *FileWriter) Write(buf []byte) error {
	if err := os.WriteFile(w.FilePath, buf, 0644); err != nil {
		return err
	}
	return nil
}

//go:generate go run github.com/ysuzuki19/robustruct/cmd/generators/testdocgen -file=$GOFILE
