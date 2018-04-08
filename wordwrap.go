// Package wordwrap provide a utility to wrap text on word boundaries.
package wordwrap

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

// Scanner wraps UTF-8 encoded text at word boundaries when lines exceed a limit
// number of characters. Newlines are preserved, including consecutive and
// trailing newlines, though trailing whitespace is stripped from each line.
//
// Clients should not assume Scanner is safe for parallel execution.
type Scanner struct {
	r        io.RuneScanner
	limit    int
	prefix   string
	tabWidth int

	// Scan state
	err         error
	line        bytes.Buffer
	word        bytes.Buffer
	space       bytes.Buffer
	lineChars   int
	wordChars   int
	spaceChars  int
	needNewline bool
	skipNextWS  bool // Skip non-newline whitespace if true.
}

// NewScanner creates and initializes a new Scanner given a reader and fixed
// line limit. The new Scanner takes ownership of the reader, and the caller
// should not use it after this call.
func NewScanner(r io.Reader, limit int) *Scanner {
	rs, ok := r.(io.RuneScanner)
	if !ok {
		rs = bufio.NewReader(r)
	}
	return &Scanner{r: rs, limit: limit, tabWidth: 4}
}

// SetPrefix sets a string to prefix each future line. The prefix is not applied
// to empty lines and the prefix's length is not included in the character limit
// specified in NewScanner.
//
// It's safe to call SetPrefix between calls to ReadLine.
func (s *Scanner) SetPrefix(prefix string) {
	s.prefix = prefix
}

// SetTabWidth sets the width of tab characters.
//
// It's safe to call SetTabWidth between calls to ReadLine.
func (s *Scanner) SetTabWidth(width int) {
	s.tabWidth = width
}

// ReadLine reads a single wrapped line, not including end-of-line characters
// ("\n"). Trailing newlines are preserved. At EOF, the result will be an empty
// string and the error will be io.EOF.
//
// ReadLine always attempts to return at least one line, even on empty input.
//
// ReadLine attempts to handle tab characters gracefully, converting them to
// spaces aligned on the boundary define in SetTabWidth.
func (s *Scanner) ReadLine() (string, error) {
	if s.err != nil {
		return "", s.err
	}

	for {
		var char rune
		char, _, err := s.r.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			s.err = err
			return "", err
		}

		if unicode.IsSpace(char) {
			if _, err := s.flushWord(&s.line); err != nil {
				s.err = err
				return "", err
			}

			if char == '\n' {
				ret := s.line.String()
				s.lineChars = 0
				s.skipNextWS = false
				s.line.Reset()
				s.space.Reset()
				return ret, nil
			}

			if s.skipNextWS {
				continue
			}

			if char == '\t' {
				// Replace tabs with spaces while preserving alignment.
				count := s.tabWidth - s.lineChars%s.tabWidth
				s.space.WriteString(strings.Repeat(" ", count))
				s.spaceChars += count
			} else {
				if _, err := s.space.WriteRune(char); err != nil {
					s.err = err
					return "", err
				}
				s.spaceChars++
			}
		} else {
			s.word.WriteRune(char)
			s.wordChars++
			s.skipNextWS = false

			if s.needNewline {
				ret := s.line.String()
				s.needNewline = false
				s.line.Reset()
				return ret, nil
			}
		}

		// Commit the line if we've reached the maximum width.
		if s.lineChars+s.wordChars+s.spaceChars >= s.limit {
			next, nextSize, err := peekRune(s.r)
			if err != nil && err != io.EOF {
				s.err = err
				return "", err
			}

			// Flush if the next character constitutes a word break.
			if s.wordChars == s.limit || unicode.IsSpace(next) || nextSize == 0 {
				if _, err := s.flushWord(&s.line); err != nil {
					s.err = err
					return "", err
				}
			}

			s.lineChars = 0
			if nextSize != 0 && next != '\n' && s.spaceChars < s.limit {
				// We had some non-whitespace chars, so start a new line for the next write.
				s.needNewline = true
			}

			s.skipNextWS = true
			s.space.Reset()
		}
	}

	if _, err := s.flushWord(&s.line); err != nil {
		s.err = err
		return "", err
	}

	ret := s.line.String()
	s.line.Reset()
	s.err = io.EOF
	return ret, nil
}

// WriteTo implements io.WriterTo. This may make multiple calls to the Read
// method of the underlying Reader.
func (s *Scanner) WriteTo(w io.Writer) (n int64, err error) {
	firstLine := true
	newline := []byte("\n")
	for {
		line, err := s.ReadLine()
		if err == io.EOF {
			return n, nil
		} else if err != nil {
			return n, err
		}

		if !firstLine {
			written, err := w.Write(newline)
			n += int64(written)
			if err != nil {
				return n, err
			}
		}

		written, err := io.WriteString(w, line)
		n += int64(written)
		if err != nil {
			return n, err
		}

		firstLine = false
	}
}

func (s *Scanner) flushWord(w io.Writer) (int, error) {
	var written int
	if s.wordChars > 0 {
		if s.lineChars == 0 {
			n, err := io.WriteString(w, s.prefix)
			written += n
			if err != nil {
				return written, err
			}
		}

		n, err := s.space.WriteTo(w)
		written += int(n)
		if err != nil {
			return written, err
		}

		n, err = s.word.WriteTo(w)
		written += int(n)
		if err != nil {
			return written, err
		}

		s.lineChars += s.spaceChars + s.wordChars
		s.spaceChars, s.wordChars = 0, 0
	}
	return written, nil
}

func peekRune(r io.RuneScanner) (rune, int, error) {
	ch, size, err := r.ReadRune()
	if err != nil {
		return ch, size, err
	}
	if err := r.UnreadRune(); err != nil {
		return 0, 0, err
	}
	return ch, size, nil
}
