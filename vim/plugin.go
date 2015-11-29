// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vim

import (
	"fmt"
	"io"
	"log"
	"os"
)

var setupFuncs []func(*Vim) error

func RegisterPluginSetup(f func(v *Vim) error) {
	setupFuncs = append(setupFuncs, f)
}

// PluginMain implements the main function for a Neovim remote plugin.
func PluginMain() {
	if fname := os.Getenv("NVIMGO_LOG_FILE"); fname != "" {
		f, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.SetOutput(f)
		log.SetPrefix(fmt.Sprintf("%8d ", os.Getpid()))
		log.Print("Plugin Start")
		defer log.Print("Plugin Exit")
	}
	v, err := New(struct {
		io.Reader
		io.WriteCloser
	}{
		os.Stdin,
		os.Stdout,
	}, log.Printf)
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) >= 2 {
		v.PluginPath = os.Args[1]
	}
	for _, f := range setupFuncs {
		if err := f(v); err != nil {
			log.Fatal(err)
		}
	}
	if err := v.Serve(); err != nil {
		log.Fatal(err)
	}
}
