// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vim_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/garyburd/nvimgo/vim"
	"github.com/garyburd/nvimgo/vim/vimtest"
)

func helloHandler(v *vim.Vim, s string) (string, error) {
	return "Hello, " + s, nil
}

func helloFunc(v *vim.Vim, args []string) (string, error) {
	return "Hello, " + strings.Join(args, " "), nil
}

func TestRegister(t *testing.T) {
	v, cleanup := vimtest.New(t)
	defer cleanup()

	cid, err := v.ChannelID()
	if err != nil {
		t.Fatal(err)
	}

	if err := v.RegisterHandler("hello", helloHandler); err != nil {
		t.Fatal(err)
	}

	var result string
	if err := v.Call("rpcrequest", &result, cid, "hello", "world"); err != nil {
		t.Fatal(err)
	}

	expected := "Hello, world"
	if result != expected {
		t.Errorf("hello returned %q, want %q", result, expected)
	}
}

func TestAPI(t *testing.T) {
	v, cleanup := vimtest.New(t)
	defer cleanup()

	cid, err := v.ChannelID()
	if err != nil {
		t.Fatal(err)
	}

	if err := v.RegisterHandler("hello", helloHandler); err != nil {
		t.Fatal(err)
	}

	var result string
	if err := v.Call("rpcrequest", &result, cid, "hello", "world"); err != nil {
		t.Fatal(err)
	}

	expected := "Hello, world"
	if result != expected {
		t.Errorf("hello returned %q, want %q", result, expected)
	}

	// Buffers

	bufs, err := v.Buffers()
	if err != nil {
		t.Fatal(err)
	}
	if len(bufs) != 1 {
		t.Errorf("expected one buf, found %d bufs", len(bufs))
	}
	if len(bufs[0]) == 0 {
		t.Errorf("buf len == 0: %s", bufs[0])
	}
	buf, err := v.CurrentBuffer()
	if err != nil {
		t.Fatal(err)
	}
	if buf != bufs[0] {
		t.Fatal("buf %v != bufs[0] %v", buf, buf[0])
	}
	err = v.SetCurrentBuffer(buf)
	if err != nil {
		t.Fatal(err)
	}

	// Windows

	wins, err := v.Windows()
	if err != nil {
		t.Fatal(err)
	}
	if len(wins) != 1 {
		t.Errorf("expected one win, found %d wins", len(wins))
	}
	if len(wins[0]) == 0 {
		t.Errorf("win len == 0: %s", wins[0])
	}
	win, err := v.CurrentWindow()
	if err != nil {
		t.Fatal(err)
	}
	if win != wins[0] {
		t.Fatal("win %v != wins[0] %v", win, win[0])
	}
	err = v.SetCurrentWindow(win)
	if err != nil {
		t.Fatal(err)
	}

	// Tabpage

	pages, err := v.Tabpages()
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) != 1 {
		t.Errorf("expected one page, found %d pages", len(pages))
	}
	if len(pages[0]) == 0 {
		t.Errorf("page len == 0: %s", pages[0])
	}
	page, err := v.CurrentTabpage()
	if err != nil {
		t.Fatal(err)
	}
	if page != pages[0] {
		t.Fatal("page %v != pages[0] %v", page, page[0])
	}
	err = v.SetCurrentTabpage(page)
	if err != nil {
		t.Fatal(err)
	}

	// Lines

	lines := [][]byte{[]byte("hello"), []byte("world")}
	if err := v.SetBufferLineSlice(buf, 0, -1, true, true, lines); err != nil {
		t.Fatal(err)
	}
	lines2, err := v.BufferLineSlice(buf, 0, -1, true, true)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(lines2, lines) {
		t.Fatalf("lines = %+v, want %+v", lines2, lines)
	}

}

func TestPlugin(t *testing.T) {
	v, cleanup := vimtest.New(t, func(v *vim.Vim) error {
		return v.RegisterFunction("Hello", helloFunc)
	})
	defer cleanup()

	result, err := v.CommandOutput(":echo Hello('John', 'Doe')")
	if err != nil {
		t.Error(err)
	}
	expected := "\nHello, John Doe"
	if result != expected {
		t.Errorf("Hello returned %q, want %q", result, expected)
	}
}
