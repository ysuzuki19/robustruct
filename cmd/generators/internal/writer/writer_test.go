package writer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ysuzuki19/robustruct/cmd/generators/internal/writer"
)

func TestMemoryWriter(t *testing.T) {
	require := require.New(t)
	// testdoc begin MemoryWriter.Write
	w := &writer.MemoryWriter{}
	err := w.Write([]byte("Hello, World!"))
	require.NoError(err)
	require.Equal("Hello, World!", w.Content)
	// testdoc end
}
