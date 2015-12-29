// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vim_test

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/garyburd/neovim-go/vim"
	"github.com/garyburd/neovim-go/vim/vimtest"
)

func helloHandler(v *vim.Vim, s string) (string, error) {
	return "Hello, " + s, nil
}

func TestAPI(t *testing.T) {
	v, cleanup := vimtest.New(t, false)
	defer cleanup()

	cid, err := v.ChannelID()
	if err != nil {
		t.Fatal(err)
	}

	{
		// Simple handler.

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

	// Buffers

	bufs, err := v.Buffers()
	if err != nil {
		t.Fatal(err)
	}
	if len(bufs) != 1 {
		t.Errorf("expected one buf, found %d bufs", len(bufs))
	}
	if bufs[0] == 0 {
		t.Errorf("bufs[0] == 0")
	}
	buf, err := v.CurrentBuffer()
	if err != nil {
		t.Fatal(err)
	}
	if buf != bufs[0] {
		t.Fatal("buf %v != bufs[0] %v", buf, bufs[0])
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
	if wins[0] == 0 {
		t.Errorf("wins[0] == 0")
	}
	win, err := v.CurrentWindow()
	if err != nil {
		t.Fatal(err)
	}
	if win != wins[0] {
		t.Fatal("win %v != wins[0] %v", win, wins[0])
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
	if pages[0] == 0 {
		t.Errorf("pages[0] == 0")
	}
	page, err := v.CurrentTabpage()
	if err != nil {
		t.Fatal(err)
	}
	if page != pages[0] {
		t.Fatal("page %v != pages[0] %v", page, pages[0])
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

	// Vars

	if err := v.SetVar("foo", "bar", nil); err != nil {
		t.Fatal(err)
	}
	var foo interface{}
	if err := v.Var("foo", &foo); err != nil {
		t.Fatal(err)
	}
	if foo != "bar" {
		t.Errorf("got %v, want %q", foo, "bar")
	}
	if err := v.SetVar("foo", "", nil); err != nil {
		t.Fatal(err)
	}
	foo = nil
	if err := v.Var("foo", &foo); err != nil {
		t.Fatal(err)
	}
	if foo != "" {
		t.Errorf("got %v, want %q", foo, "")
	}

	{
		// Pipeline

		p := v.NewPipeline()
		results := make([]int, 128)

		for i := range results {
			p.SetVar(fmt.Sprintf("v%d", i), i, nil)
		}

		for i := range results {
			p.Var(fmt.Sprintf("v%d", i), &results[i])
		}

		if err := p.Wait(); err != nil {
			t.Fatal(err)
		}

		for i := range results {
			if results[i] != i {
				t.Fatalf("result = %d, want %d", results[i], i)
			}
		}

		// Reuse pipeline

		var i int
		p.Var("v3", &i)
		if err := p.Wait(); err != nil {
			log.Fatal(err)
		}
		if i != 3 {
			t.Fatalf("i = %d, want %d", i, 3)
		}
	}

	{
		// Call with no args.

		var wd string
		err := v.Call("getcwd", &wd)
		if err != nil {
			t.Fatal(err)
		}
	}

}
