package wordwrap

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	message  string
	text     string
	width    int
	prefix   string
	expected string
}

var allCases = map[string][]testCase{
	"BasicText": []testCase{
		{
			"Long words should be broken up.",
			"stupendous", 4, "",
			"stup\nendo\nus",
		},
		{
			"Words should be broken on spaces.",
			"foo bar baz", 4, "",
			"foo\nbar\nbaz",
		},
		{
			"Leading/trailing space should be trimmed on wrap.",
			"foo bar baz  ", 3, "",
			"foo\nbar\nbaz",
		},
		{
			"Contiguous spaces should be preserved.",
			"foo  bar  baz", 9, "",
			"foo  bar\nbaz",
		},
		{
			"Words that would run over should be wrapped.",
			"foo bar", 5, "",
			"foo\nbar",
		},
		{
			"Multiple words can fit on one line.",
			"This should split to two lines.", 20, "",
			"This should split to\ntwo lines.",
		},
		{
			"Multiple words of exact line width should fit.",
			"Nineteen characters", 19, "",
			"Nineteen characters",
		},
		{
			"Long runs of spaces should be trimmed.",
			"foo            bar", 5, "",
			"foo\nbar",
		},
	},
	"Newlines": {
		{
			"Newlines should always wrap.",
			"foo\nbar baz", 8, "",
			"foo\nbar baz",
		},
		{
			"Newline after full line should no-op.",
			"foo\nbar", 3, "",
			"foo\nbar",
		},
		{
			"Trailing space before newline should be trimmed.",
			"foo \nbar", 5, "",
			"foo\nbar",
		},
		{
			"Explicit leading space should be preserved.",
			"foo\n  bar", 8, "",
			"foo\n  bar",
		},
		{
			"Empty lines should be preserved.",
			"foo\n\n\nbar\n", 4, "",
			"foo\n\n\nbar\n",
		},
		{
			"Lines of all whitespace should be trimmed.",
			"first\n  \nlast\n  ", 8, "",
			"first\n\nlast\n",
		},
	},
	"Tabs": {
		{
			"Leading tabs should be trimmed like other whitespace.",
			"foo\tbar", 3, "",
			"foo\nbar",
		},
		{
			"Tabs after newlines should be preserved like other whitespace.",
			"foo\n\tbar", 8, "",
			"foo\n    bar",
		},
		{
			"Split tabs should be trimmed on both lines.",
			"foo\tbar", 5, "",
			"foo\nbar",
		},
		{
			"Tabs should maintain alignment.",
			"1\tfoo", 8, "",
			"1   foo",
		},
		{
			"Tabs should maintain alignment.",
			"22\tfoo", 8, "",
			"22  foo",
		},
		{
			"Tabs should maintain alignment.",
			"333\tfoo", 8, "",
			"333 foo",
		},
		{
			"Tabs should maintain alignment.",
			"4444\tfoo", 12, "",
			"4444    foo",
		},
	},
	"Prefix": {
		{
			"Prefix should be applied to wrapped lines.",
			"foo bar baz", 4, "--",
			"--foo\n--bar\n--baz",
		},
		{
			"Prefix should be applied to split words.",
			"reallylongword", 4, "  ",
			"  real\n  lylo\n  ngwo\n  rd",
		},
		{
			"Prefix should be applied to explicit newlines.",
			"foo\nbar", 8, "  ",
			"  foo\n  bar",
		},
		{
			"Prefix should not be applied to empty lines.",
			"foo\n\nbar\n", 8, "++",
			"++foo\n\n++bar\n",
		},
		{
			"Prefix should be applied to single lines.",
			"foo", 4, "  ",
			"  foo",
		},
	},
	"Degenerate": {
		{
			"Empty string",
			"", 4, "++",
			"",
		},
		{
			"String length is exactly width.",
			"foo", 3, "++",
			"++foo",
		},
		{
			"Input is all spaces.",
			"   ", 4, "",
			"",
		},
		{
			"Space crossing multiple line boundaries",
			"           ", 4, "",
			"",
		},
		{
			// There's no right way to handle this case, so this test is
			// arbitrary and exists only to enforce fixed behavior.
			"Newline followed by too much indentation",
			"foo\n     bar", 4, "",
			"foo\nbar",
		},
	},
}

func TestReadLine(t *testing.T) {
	for name, cases := range allCases {
		t.Run(name, func(t *testing.T) {
			for _, c := range cases {
				s := NewScanner(strings.NewReader(c.text), c.width)
				s.SetPrefix(c.prefix)

				expected := strings.Split(c.expected, "\n")

				var lines []string
				for {
					line, err := s.ReadLine()
					if err == io.EOF {
						break
					}

					require.NoError(t, err)
					lines = append(lines, line)
				}

				assert.Equal(t, expected, lines, c.message)
			}
		})
	}
}

func TestWriteTo(t *testing.T) {
	for name, cases := range allCases {
		t.Run(name, func(t *testing.T) {
			for _, c := range cases {
				s := NewScanner(strings.NewReader(c.text), c.width)
				s.SetPrefix(c.prefix)

				buf := new(bytes.Buffer)
				n, err := s.WriteTo(buf)
				require.NoError(t, err)
				assert.Equal(t, c.expected, buf.String(), c.message)
				assert.Equal(t, len(c.expected), int(n), c.message)
			}
		})
	}
}
