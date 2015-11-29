// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command nvimgo is a Neovim remote plogin.
package main

import (
	_ "github.com/garyburd/nvimgo/gofmt"
	"github.com/garyburd/nvimgo/vim"
)

func main() {
	vim.PluginMain()
}
