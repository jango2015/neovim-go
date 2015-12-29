// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vim implements a Neovim client.
package vim

import (
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/garyburd/neovim-go/msgpack"
	"github.com/garyburd/neovim-go/msgpack/rpc"
)

//go:generate go run genapi.go -out api.go

// Vim represents a remote instance of Neovim.
type Vim struct {
	ep        *rpc.Endpoint
	mu        sync.Mutex
	channelID int
}

func (v *Vim) Serve() error {
	return v.ep.Serve()
}

func (v *Vim) Close() error {
	return v.ep.Close()
}

// New create a Neovim client. When connecting to Neovim over stdio, use stdin
// as r and stdout as wc. When connecting to Neovim over a network connection,
// use the connection for both r and wc.
//
//  :help msgpack-rpc-connecting
func New(r io.Reader, wc io.WriteCloser, logf func(string, ...interface{})) (*Vim, error) {
	v := &Vim{}

	rwc := struct {
		io.Reader
		io.WriteCloser
	}{r, wc}

	var err error
	v.ep, err = rpc.NewEndpoint(rwc, withExtensions(), rpc.WithLogf(logf), rpc.WithFirstArg(v))
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (v *Vim) RegisterHandler(serviceMethod string, function interface{}) error {
	return v.ep.RegisterHandler(serviceMethod, function)
}

// ChannelID returns the peer's channel id in Neovim.
func (v *Vim) ChannelID() (int, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.channelID != 0 {
		return v.channelID, nil
	}
	var info struct {
		ChannelID int `msgpack:",array"`
		Info      interface{}
	}
	if err := v.ep.Call("vim_get_api_info", &info); err != nil {
		return 0, err
	}
	v.channelID = info.ChannelID
	return v.channelID, nil
}

func (v *Vim) call(sm string, result interface{}, args ...interface{}) error {
	return fixError(sm, v.ep.Call(sm, result, args...))
}

// NewPipeline creates a new pipeline.
func (v *Vim) NewPipeline() *Pipeline {
	return &Pipeline{ep: v.ep}
}

// Pipeline pipelines calls to Neovim. Call the Wait method to wait for the calls
// to complete.
//
// Pipelines do not support concurrent access.
type Pipeline struct {
	ep    *rpc.Endpoint
	n     int
	done  chan *rpc.Call
	chans []chan *rpc.Call
}

const doneChunkSize = 32

func (p *Pipeline) call(sm string, result interface{}, args ...interface{}) {
	if p.n%doneChunkSize == 0 {
		done := make(chan *rpc.Call, doneChunkSize)
		p.done = done
		p.chans = append(p.chans, done)
	}
	p.n++
	p.ep.Go(sm, p.done, result, args...)
}

// Wait waits for all calls in the pipeline to complete.
func (p *Pipeline) Wait() error {
	var el ErrorList
	var done chan *rpc.Call
	for i := 0; i < p.n; i++ {
		if i%doneChunkSize == 0 {
			done = p.chans[0]
			p.chans = p.chans[1:]
		}
		c := <-done
		if c.Err != nil {
			el = append(el, fixError(c.ServiceMethod, c.Err))
		}
	}
	p.n = 0
	p.done = nil
	p.chans = nil
	if len(el) == 0 {
		return nil
	}
	return el
}

func fixError(sm string, err error) error {
	if e, ok := err.(rpc.Error); ok {
		if a, ok := e.Value.([]interface{}); ok && len(a) == 2 {
			switch a[0] {
			case int64(exceptionError), uint64(exceptionError):
				return fmt.Errorf("nvim:%s exception: %v", sm, a[1])
			case int64(validationError), uint64(validationError):
				return fmt.Errorf("nvim:%s validation: %v", sm, a[1])
			}
		}
	}
	return err
}

type ErrorList []error

func (el ErrorList) Error() string {
	return el[0].Error()
}

// Call calls a vimscript function.
func (v *Vim) Call(fname string, result interface{}, args ...interface{}) error {
	if args == nil {
		args = []interface{}{}
	}
	return v.call("vim_call_function", result, fname, args)
}

// Call calls a vimscript function.
func (p *Pipeline) Call(fname string, result interface{}, args ...interface{}) {
	if args == nil {
		args = []interface{}{}
	}
	p.call("vim_call_function", result, fname, args)
}

const (
	bufferExt  = 0
	windowExt  = 1
	tabpageExt = 2

	exceptionError  = 0
	validationError = 1
)

func withExtensions() rpc.Option {
	return rpc.WithExtensions(msgpack.ExtensionMap{
		bufferExt: func(p []byte) (interface{}, error) {
			x, err := decodeExt(p)
			return Buffer(x), err
		},
		windowExt: func(p []byte) (interface{}, error) {
			x, err := decodeExt(p)
			return Window(x), err
		},
		tabpageExt: func(p []byte) (interface{}, error) {
			x, err := decodeExt(p)
			return Tabpage(x), err
		},
	})
}

// decodeInt decodes a MsgPack encoded number to an integer.
func decodeExt(p []byte) (int, error) {
	switch {
	case len(p) == 1 && p[0] <= 0x7f:
		return int(p[0]), nil
	case len(p) == 2 && p[0] == 0xcc:
		return int(p[1]), nil
	case len(p) == 3 && p[0] == 0xcd:
		return int(uint16(p[2]) | uint16(p[1])<<8), nil
	case len(p) == 5 && p[0] == 0xce:
		return int(uint32(p[4]) | uint32(p[3])<<8 | uint32(p[2])<<16 | uint32(p[1])<<24), nil
	case len(p) == 2 && p[0] == 0xd0:
		return int(int8(p[1])), nil
	case len(p) == 3 && p[0] == 0xd1:
		return int(int16(uint16(p[2]) | uint16(p[1])<<8)), nil
	case len(p) == 5 && p[0] == 0xd2:
		return int(int32(uint32(p[4]) | uint32(p[3])<<8 | uint32(p[2])<<16 | uint32(p[1])<<24)), nil
	case len(p) == 1 && p[0] >= 0xe0:
		return int(int8(p[0])), nil
	default:
		return 0, fmt.Errorf("nvimgo: error decoding extension bytes %x", p)
	}
}

// encodeInt encodes n to MsgPack format.
func encodeExt(n int) []byte {
	return []byte{0xd2, byte(n >> 24), byte(n >> 16), byte(n >> 8), byte(n)}
}

type Buffer int

func (x *Buffer) UnmarshalMsgPack(dec *msgpack.Decoder) error {
	if dec.Type() != msgpack.Extension || dec.Extension() != bufferExt {
		err := &msgpack.DecodeConvertError{
			SrcType:  dec.Type(),
			DestType: reflect.TypeOf(x),
		}
		dec.Skip()
		return err
	}
	n, err := decodeExt(dec.BytesNoCopy())
	*x = Buffer(n)
	return err
}

func (x Buffer) MarshalMsgPack(enc *msgpack.Encoder) error {
	return enc.PackExtension(bufferExt, encodeExt(int(x)))
}

func (x Buffer) String() string {
	return fmt.Sprintf("Buffer:%d", int(x))
}

type Window int

func (x *Window) UnmarshalMsgPack(dec *msgpack.Decoder) error {
	if dec.Type() != msgpack.Extension || dec.Extension() != windowExt {
		err := &msgpack.DecodeConvertError{
			SrcType:  dec.Type(),
			DestType: reflect.TypeOf(x),
		}
		dec.Skip()
		return err
	}
	n, err := decodeExt(dec.BytesNoCopy())
	*x = Window(n)
	return err
}

func (x Window) MarshalMsgPack(enc *msgpack.Encoder) error {
	return enc.PackExtension(windowExt, encodeExt(int(x)))
}

func (x Window) String() string {
	return fmt.Sprintf("Window:%d", int(x))
}

type Tabpage int

func (x *Tabpage) UnmarshalMsgPack(dec *msgpack.Decoder) error {
	if dec.Type() != msgpack.Extension || dec.Extension() != tabpageExt {
		err := &msgpack.DecodeConvertError{
			SrcType:  dec.Type(),
			DestType: reflect.TypeOf(x),
		}
		dec.Skip()
		return err
	}
	n, err := decodeExt(dec.BytesNoCopy())
	*x = Tabpage(n)
	return err
}

func (x Tabpage) MarshalMsgPack(enc *msgpack.Encoder) error {
	return enc.PackExtension(tabpageExt, encodeExt(int(x)))
}

func (x Tabpage) String() string {
	return fmt.Sprintf("Tabpage:%d", int(x))
}
