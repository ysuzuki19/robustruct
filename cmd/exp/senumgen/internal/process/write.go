package process

import (
	"log"
	"os"
)

type Write interface {
	Write(buf []byte) error
}

var _ Write = &MemoryWriter{} //ignore:fields_require
var _ Write = &FileWriter{}   //ignore:fields_require

type MemoryWriter struct {
	Content string
}

func (w *MemoryWriter) Write(buf []byte) error {
	w.Content = string(buf)
	return nil
}

type FileWriter struct {
	FilePath string
}

func (w *FileWriter) Write(buf []byte) error {
	if err := os.WriteFile(w.FilePath, buf, 0644); err != nil {
		log.Fatal(err)
	}
	return nil
}
