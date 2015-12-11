// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vimutil implements utilities for working with Vim.
package vimutil

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
