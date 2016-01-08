// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vimutil implements utilities for working with Vim.
package vimutil

import (
	"bytes"
	"io"
	"strconv"

	"github.com/garyburd/neovim-go/vim"
)

// QuickfixError represents an item in a quickfix list.
type QuickfixError struct {
	// Buffer number
	Bufnr int `msgpack:"bufnr,omitempty"`

	// Line number in the file.
	LNum int `msgpack:"lnum,omitempty"`

	// Search pattern used to locate the error.
	Pattern string `msgpack:"pattern,omitempty"`

	// Column number (first column is 1).
	Col int `msgpack:"col,omitempty"`

	// When Vcol is != 0,  Col is visual column.
	VCol int `msgpack:"vcol,omitempty"`

	// Error number.
	Nr int `msgpack:"nr,omitempty"`

	// Description of the error.
	Text string `msgpack:"text,omitempty"`

	// Single-character error type, 'E', 'W', etc.
	Type string `msgpack:"type,omitempty"`

	// Name of a file; only used when bufnr is not present or it is invalid.
	FileName string `msgpack:"filename,omitempty"`

	// Valid is non-zero if this is a recognized error message.
	Valid int `msgpack:"valid,omitempty"`
}

// CommandCompletionArgs represents the arguments to a custom command line
// completion function.
//
//  :help :command-completion-custom
type CommandCompletionArgs struct {
	// ArgLead is the leading portion of the argument currently being completed
	// on.
	ArgLead string `msgpack:",array"`

	// CmdLine is the entire command line.
	CmdLine string

	// CursorPosString is decimal representation of the cursor position in
	// bytes.
	CursorPosString string
}

// CursorPos returns the cursor position.
func (a *CommandCompletionArgs) CursorPos() int {
	n, _ := strconv.Atoi(a.CursorPosString)
	return n
}

type bufferReader struct {
	v *vim.Vim
	r io.Reader
}

// CurrentBufferReader returns a reader for the current buffer.
func CurrentBufferReader(v *vim.Vim) io.Reader {
	return &bufferReader{v: v}
}

func (br *bufferReader) Read(p []byte) (int, error) {
	if br.r == nil {
		b, err := br.v.CurrentBuffer()
		if err != nil {
			return 0, err
		}
		lines, err := br.v.BufferLineSlice(b, 0, -1, true, true)
		if err != nil {
			return 0, err
		}
		br.r = bytes.NewReader(bytes.Join(lines, []byte{'\n'}))
	}
	return br.r.Read(p)
}
