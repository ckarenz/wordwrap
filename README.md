# wordwrap

[![GoDoc Reference](https://godoc.org/github.com/ckarenz/wordcount?status.svg)](http://godoc.org/github.com/ckarenz/wordcount)
[![Build Status](https://travis-ci.org/ckarenz/wordcount.svg?branch=master)](https://travis-ci.org/ckarenz/wordcount)

A Go library for word-wrapping text.

## Features

- Text is guaranteed to never exceed line width.
- Support for multi-byte (utf8) text
- Handling for tab width and alignment; tabs are replaced by spaces
- Streaming: text need not be loaded into a buffer.
- Settings (line prefix, tab width) can be changed on the fly.
