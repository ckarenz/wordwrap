package wordwrap

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteRune(t *testing.T) {
	const r = '😀'
	b := runeBuffer{}
	n, err := b.WriteRune(r)
	require.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Len(t, b.String(), 4)
	assert.Equal(t, 1, b.Count())
	assert.Equal(t, string(r), b.String())
}

func TestWriteEmptyString(t *testing.T) {
	b := runeBuffer{}
	n, err := b.WriteString("")
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Len(t, b.String(), 0)
	assert.Equal(t, 0, b.Count())
	assert.Equal(t, "", b.String())
}

func TestWriteString(t *testing.T) {
	const s = "Käse"
	b := runeBuffer{}
	n, err := b.WriteString(s)
	require.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Len(t, b.String(), 5)
	assert.Equal(t, 4, b.Count())
	assert.Equal(t, s, b.String())
	fmt.Println(b.String())
}

func TestBufWriteTo(t *testing.T) {
	const s = "Test"
	b1 := runeBuffer{}
	b1.WriteString(s)
	require.Equal(t, s, b1.String())

	b2 := runeBuffer{}
	n, err := b1.WriteTo(&b2)
	require.NoError(t, err)
	assert.Equal(t, int64(4), n)

	assert.Equal(t, 0, b1.Count())
	assert.Equal(t, 4, b2.Count())
	assert.Equal(t, "", b1.String())
	assert.Equal(t, s, b2.String())
}
