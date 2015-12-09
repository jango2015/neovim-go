// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vimtest

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/garyburd/neovim-go/vim"
)

func newEmbeddedVim(args []string, logf func(string, ...interface{})) (*vim.Vim, *os.Process, error) {
	cmdPath, err := exec.LookPath(args[0])
	if err != nil {
		return nil, nil, err
	}

	outr, outw, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	defer outw.Close()

	inr, inw, err := os.Pipe()
	if err != nil {
		outr.Close()
		return nil, nil, err
	}
	defer inr.Close()

	v, err := vim.New(struct {
		io.Reader
		io.WriteCloser
	}{
		outr,
		inw,
	}, logf)
	if err != nil {
		outr.Close()
		inw.Close()
		return nil, nil, err
	}

	p, err := os.StartProcess(cmdPath, args, &os.ProcAttr{Env: []string{}, Files: []*os.File{inr, outw}})
	if err != nil {
		outr.Close()
		inw.Close()
		return nil, nil, err
	}

	return v, p, nil
}

func New(t *testing.T, setupFuncs ...func(*vim.Vim) error) (v *vim.Vim, cleanup func()) {
	v, p, err := newEmbeddedVim(
		[]string{"nvim", "--embed", "-u", "NONE", "-n"},
		t.Logf,
	)
	if err != nil {
		t.Fatal(err)
	}

	v.PluginPath = "x"

	for _, f := range setupFuncs {
		if err := f(v); err != nil {
			t.Fatal(err)
		}
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		if err := v.Serve(); err != nil {
			t.Fatal(err)
		}

	}()

	if len(setupFuncs) > 0 {
		cid, err := v.ChannelID()
		if err != nil {
			t.Fatal(err)
		}

		if err := v.Command(":call remote#host#RegisterPlugin('nvimgo', 'x', rpcrequest(1, 'specs', 'x'))"); err != nil {
			t.Error(err)
		}

		if err := v.Command(fmt.Sprintf(":call remote#host#Register('nvimgo', 'x', %d)", cid)); err != nil {
			t.Error(err)
		}
	}

	return v, func() {
		v.Close()
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Errorf("timeout waiting for close")
		}
		p.Kill()
		p.Release()
	}
}
