// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vim

import (
	"fmt"
	"reflect"

	"github.com/garyburd/nvimgo/msgpack"
	"github.com/garyburd/nvimgo/msgpack/rpc"
)

const (
	bufferExt  = 0
	windowExt  = 1
	tabpageExt = 2
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
