package writer

import (
	"log"
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
