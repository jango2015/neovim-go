// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plugin_test

import (
	"testing"

	"github.com/garyburd/neovim-go/vim"
	"github.com/garyburd/neovim-go/vim/plugin"
	"github.com/garyburd/neovim-go/vim/vimtest"
)

func init() {
	plugin.Handle("hello", func(v *vim.Vim, s string) (string, error) {
		return "Hello, " + s, nil
	})
}

func TestRegister(t *testing.T) {
	v, cleanup := vimtest.New(t, true)
	defer cleanup()

	cid, err := v.ChannelID()
	if err != nil {
		t.Fatal(err)
	}

	{
		result, err := v.CommandOutput(":echo Hello('John', 'Doe')")
		if err != nil {
			t.Error(err)
		}
		expected := "\nHello, John Doe"
		if result != expected {
			t.Errorf("Hello returned %q, want %q", result, expected)
		}
	}

	{
		var result string
		if err := v.Call("rpcrequest", &result, cid, "hello", "world"); err != nil {
			t.Fatal(err)
		}

		expected := "Hello, world"
		if result != expected {
			t.Errorf("hello returned %q, want %q", result, expected)
		}
	}
}
