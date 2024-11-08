package chain_writer_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ysuzuki19/robustruct/internal/chain_writer"
)

func TestChainWriter(t *testing.T) {
	cw := chain_writer.New(nil)
	cw.Push("hell world\n")
	s, err := cw.String()
	require.NoError(t, err)
	require.Equal(t, "hell world\n", s)

	cw.Push([]byte("hello world\n"))
	b, err := cw.Bytes()
	require.NoError(t, err)
	require.Equal(t, []byte("hell world\nhello world\n"), b)

	cw.Push(42)
	_, err = cw.String()
	require.Error(t, err)
}
