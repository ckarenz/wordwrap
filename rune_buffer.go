package wordwrap

import (
	"bytes"
	"unicode/utf8"
)

type runeBuffer struct {
	buf       bytes.Buffer
	runeCount int
}

func (b *runeBuffer) Count() int     { return b.runeCount }
func (b *runeBuffer) String() string { return b.buf.String() }

func (b *runeBuffer) WriteRune(r rune) (n int, err error) {
	n, err = b.buf.WriteRune(r)
	if err == nil {
		b.runeCount++
	}
	return
}

func (b *runeBuffer) WriteString(s string) (n int, err error) {
	n, err = b.buf.WriteString(s)
	b.runeCount += utf8.RuneCount([]byte(s[:n]))
	return
}

func (b *runeBuffer) WriteTo(w *runeBuffer) (n int64, err error) {
	// These counts will be wrong on error, but the buffer shouldn't be used anyway.
	w.runeCount += b.runeCount
	b.runeCount = 0
	return b.buf.WriteTo(&w.buf)
}

func (b *runeBuffer) Reset() {
	b.buf.Reset()
	b.runeCount = 0
}
