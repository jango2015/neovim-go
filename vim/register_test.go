// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vim_test

import (
	"strings"
	"testing"

	"github.com/garyburd/neovim-go/vim"
	"github.com/garyburd/neovim-go/vim/vimtest"
)

func helloHandler(v *vim.Vim, s string) (string, error) {
	return "Hello, " + s, nil
}

func helloFunc(v *vim.Vim, args []string) (string, error) {
	return "Hello, " + strings.Join(args, " "), nil
}

func TestRegister(t *testing.T) {
	v, cleanup := vimtest.New(t, func(v *vim.Vim) error {
		if err := v.RegisterFunction("Hello", nil, helloFunc); err != nil {
			return err
		}
		if err := v.RegisterHandler("hello", helloHandler); err != nil {
			return err
		}
		return nil
	})
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
