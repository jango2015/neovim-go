// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vim is a Neovim remote plugin peer.
package vim

import (
	"io"
	"sync"

	"github.com/garyburd/neovim-go/msgpack"
	"github.com/garyburd/neovim-go/msgpack/rpc"
)

// Vim represents a remote instance of Neovim.
type Vim struct {
	ep *rpc.Endpoint

	mu        sync.Mutex
	channelID int

	PluginPath  string // TODO: find way to remove this.
	pluginSpecs []*pluginSpec
}

func (v *Vim) Serve() error {
	return v.ep.Serve()
}

func (v *Vim) Close() error {
	return v.ep.Close()
}

// New creates a new peer.
func New(rwc io.ReadWriteCloser, logf func(string, ...interface{})) (*Vim, error) {
	v := &Vim{pluginSpecs: []*pluginSpec{}}

	var err error
	v.ep, err = rpc.NewEndpoint(rwc, withExtensions(), rpc.WithLogf(logf), rpc.WithFirstArg(v))
	if err != nil {
		return nil, err
	}

	if err := v.ep.RegisterHandler("specs", (*Vim).handleSpecs); err != nil {
		return nil, err
	}

	return v, nil
}

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

type marshalAsString []byte

func (m marshalAsString) MarshalMsgPack(enc *msgpack.Encoder) error {
	return enc.PackStringBytes([]byte(m))
}

type marshalAsStrings [][]byte

func (m marshalAsStrings) MarshalMsgPack(enc *msgpack.Encoder) error {
	if err := enc.PackArrayLen(len(m)); err != nil {
		return err
	}
	for _, p := range m {
		if err := enc.PackStringBytes(p); err != nil {
			return err
		}
	}
	return nil
}

func (v *Vim) DelBufferLine(buffer Buffer, index int) error {
	return v.ep.Call("buffer_del_line", nil, buffer, index)
}

func (v *Vim) BufferLine(buffer Buffer, index int) ([]byte, error) {
	var result []byte
	err := v.ep.Call("buffer_get_line", &result, buffer, index)
	return result, err
}

func (v *Vim) BufferLineSlice(buffer Buffer, start int, end int, includeStart bool, includeEnd bool) ([][]byte, error) {
	var result [][]byte
	err := v.ep.Call("buffer_get_line_slice", &result, buffer, start, end, includeStart, includeEnd)
	return result, err
}

func (v *Vim) BufferMark(buffer Buffer, name string) ([2]int, error) {
	var result [2]int
	err := v.ep.Call("buffer_get_mark", &result, buffer, name)
	return result, err
}

func (v *Vim) BufferName(buffer Buffer) (string, error) {
	var result string
	err := v.ep.Call("buffer_get_name", &result, buffer)
	return result, err
}

func (v *Vim) BufferNumber(buffer Buffer) (int, error) {
	var result int
	err := v.ep.Call("buffer_get_number", &result, buffer)
	return result, err
}

func (v *Vim) BufferOption(buffer Buffer, name string, result interface{}) error {
	return v.ep.Call("buffer_get_option", result, buffer, name)
}

func (v *Vim) BufferVar(buffer Buffer, name string, result interface{}) error {
	return v.ep.Call("buffer_get_var", result, buffer, name)
}

func (v *Vim) InsertBufferLines(buffer Buffer, lnum int, lines [][]byte) error {
	return v.ep.Call("buffer_insert", nil, buffer, lnum, marshalAsStrings(lines))
}

func (v *Vim) IsBufferValid(buffer Buffer) (bool, error) {
	var result bool
	err := v.ep.Call("buffer_is_valid", &result, buffer)
	return result, err
}

func (v *Vim) BufferLineCount(buffer Buffer) (int, error) {
	var result int
	err := v.ep.Call("buffer_line_count", &result, buffer)
	return result, err
}

func (v *Vim) SetBufferLine(buffer Buffer, index int, line []byte) error {
	return v.ep.Call("buffer_set_line", nil, buffer, index, marshalAsString(line))
}

func (v *Vim) SetBufferLineSlice(buffer Buffer, start int, end int, includeStart bool, includeEnd bool, replacement [][]byte) error {
	return v.ep.Call("buffer_set_line_slice", nil, buffer, start, end, includeStart, includeEnd, marshalAsStrings(replacement))
}

func (v *Vim) SetBufferName(buffer Buffer, name string) error {
	return v.ep.Call("buffer_set_name", nil, buffer, name)
}

func (v *Vim) SetBufferOption(buffer Buffer, name string, value interface{}) error {
	return v.ep.Call("buffer_set_option", nil, buffer, name, value)
}

func (v *Vim) SetBufferVar(buffer Buffer, name string, value interface{}, result interface{}) error {
	return v.ep.Call("buffer_set_var", result, buffer, name, value)
}

func (v *Vim) TabpageVar(tabpage Tabpage, name string, result interface{}) error {
	return v.ep.Call("tabpage_get_var", result, tabpage, name)
}

func (v *Vim) TabpageWindow(tabpage Tabpage) (Window, error) {
	var result Window
	err := v.ep.Call("tabpage_get_window", &result, tabpage)
	return result, err
}

func (v *Vim) TabpageWindows(tabpage Tabpage) ([]Window, error) {
	var result []Window
	err := v.ep.Call("tabpage_get_windows", &result, tabpage)
	return result, err
}

func (v *Vim) IsTabpageValid(tabpage Tabpage) (bool, error) {
	var result bool
	err := v.ep.Call("tabpage_is_valid", &result, tabpage)
	return result, err
}

func (v *Vim) SetTabpageVar(tabpage Tabpage, name string, value interface{}, result interface{}) error {
	return v.ep.Call("tabpage_set_var", result, tabpage, name, value)
}

// Call calls a vimscript function.
func (v *Vim) Call(fname string, result interface{}, args ...interface{}) error {
	return v.ep.Call("vim_call_function", result, fname, args)
}

func (v *Vim) ChangeDirectory(dir string) error {
	return v.ep.Call("vim_change_directory", nil, dir)
}

// Command executes a single ex command.
func (v *Vim) Command(str string) error {
	return v.ep.Call("vim_command", nil, str)
}

// Command executes a single ex command and returns the output.
func (v *Vim) CommandOutput(str string) (string, error) {
	var result string
	err := v.ep.Call("vim_command_output", &result, str)
	return result, err
}

func (v *Vim) DelCurrentLine() error {
	return v.ep.Call("vim_del_current_line", nil)
}

// ErrWrite Print msg as an error message.
func (v *Vim) ErrWrite(msg string) error {
	return v.ep.Call("vim_err_write", nil, msg)
}

// Eval evaluates a vimscript expression.
func (v *Vim) Eval(str string, result interface{}) error {
	return v.ep.Call("vim_eval", result, str)
}

// FeedKeys Pushes keys to the Neovim user input buffer.  Options can be a
// string with the following character flags:
//
//  m:  Remap keys. This is default.
//  n:  Do not remap keys.
//  t:  Handle keys as if typed; otherwise they are handled as if coming from a
//      mapping. This matters for undo, opening folds, etc.
func (v *Vim) Feedkeys(keys string, mode string, escapeCsi bool) error {
	return v.ep.Call("vim_feedkeys", nil, keys, mode, escapeCsi)
}

func (v *Vim) APIInfo() ([]interface{}, error) {
	var result []interface{}
	err := v.ep.Call("vim_get_api_info", &result)
	return result, err
}

func (v *Vim) Buffers() ([]Buffer, error) {
	var result []Buffer
	err := v.ep.Call("vim_get_buffers", &result)
	return result, err
}

func (v *Vim) ColorMap() (map[string]interface{}, error) {
	var result map[string]interface{}
	err := v.ep.Call("vim_get_color_map", &result)
	return result, err
}

func (v *Vim) CurrentBuffer() (Buffer, error) {
	var result Buffer
	err := v.ep.Call("vim_get_current_buffer", &result)
	return result, err
}

func (v *Vim) CurrentLine() ([]byte, error) {
	var result []byte
	err := v.ep.Call("vim_get_current_line", &result)
	return result, err
}

func (v *Vim) CurrentTabpage() (Tabpage, error) {
	var result Tabpage
	err := v.ep.Call("vim_get_current_tabpage", &result)
	return result, err
}

func (v *Vim) CurrentWindow() (Window, error) {
	var result Window
	err := v.ep.Call("vim_get_current_window", &result)
	return result, err
}

func (v *Vim) Option(name string, result interface{}) error {
	return v.ep.Call("vim_get_option", result, name)
}

func (v *Vim) Tabpages() ([]Tabpage, error) {
	var result []Tabpage
	err := v.ep.Call("vim_get_tabpages", &result)
	return result, err
}

func (v *Vim) Var(name string, result interface{}) error {
	return v.ep.Call("vim_get_var", result, name)
}

func (v *Vim) Vvar(name string, result interface{}) error {
	return v.ep.Call("vim_get_vvar", result, name)
}

func (v *Vim) Windows() ([]Window, error) {
	var result []Window
	err := v.ep.Call("vim_get_windows", &result)
	return result, err
}

// Input pushes bytes to the Neovim low level input buffer.
//
// Unlike `feedkeys()`, this uses the lowest level input buffer and the call is
// not deferred. It returns the number of bytes actually written(which can be
// less than what was requested if the buffer is full).
func (v *Vim) Input(keys string) (int, error) {
	var result int
	err := v.ep.Call("vim_input", &result, keys)
	return result, err
}

// ListRuntimePaths returns a list of paths contained in the runtimepath
// option.
func (v *Vim) ListRuntimePaths() ([]string, error) {
	var result []string
	err := v.ep.Call("vim_list_runtime_paths", &result)
	return result, err
}

func (v *Vim) NameToColor(name string) (int, error) {
	var result int
	err := v.ep.Call("vim_name_to_color", &result, name)
	return result, err
}

// OutWrite prints msg as a normal message.
func (v *Vim) OutWrite(msg string) error {
	return v.ep.Call("vim_out_write", nil, msg)
}

// ReplaceTermcodes replaces any terminal code strings by byte sequences.  The
// returned sequences are Nvim's internal representation of keys, for example:
//
//  <esc> -> '\x1b'
//  <cr>  -> '\r'
//  <c-l> -> '\x0c'
//  <up>  -> '\x80ku'
//
// The returned sequences can be used as input to feedkeys.
func (v *Vim) ReplaceTermcodes(str string, fromPart bool, doLt bool, special bool) (string, error) {
	var result string
	err := v.ep.Call("vim_replace_termcodes", &result, str, fromPart, doLt, special)
	return result, err
}

func (v *Vim) ReportError(str string) error {
	return v.ep.Call("vim_report_error", nil, str)
}

func (v *Vim) SetCurrentBuffer(buffer Buffer) error {
	return v.ep.Call("vim_set_current_buffer", nil, buffer)
}

func (v *Vim) SetCurrentLine(line []byte) error {
	return v.ep.Call("vim_set_current_line", nil, marshalAsString(line))
}

func (v *Vim) SetCurrentTabpage(tabpage Tabpage) error {
	return v.ep.Call("vim_set_current_tabpage", nil, tabpage)
}

func (v *Vim) SetCurrentWindow(window Window) error {
	return v.ep.Call("vim_set_current_window", nil, window)
}

func (v *Vim) SetOption(name string, value interface{}) error {
	return v.ep.Call("vim_set_option", nil, name, value)
}

func (v *Vim) SetVar(name string, value interface{}, result interface{}) error {
	return v.ep.Call("vim_set_var", result, name, value)
}

// Strwidth returns the number of display cells string occupies.  Tab is
// counted as one cell.
func (v *Vim) Strwidth(str string) (int, error) {
	var result int
	err := v.ep.Call("vim_strwidth", &result, str)
	return result, err
}

// Subscribe subscribes to a Neovim event.
func (v *Vim) Subscribe(event string) error {
	return v.ep.Call("vim_subscribe", nil, event)
}

// Unsubscribe unsubscribes to a Neovim event.
func (v *Vim) Unsubscribe(event string) error {
	return v.ep.Call("vim_unsubscribe", nil, event)
}

func (v *Vim) WindowBuffer(window Window) (Buffer, error) {
	var result Buffer
	err := v.ep.Call("window_get_buffer", &result, window)
	return result, err
}

func (v *Vim) WindowCursor(window Window) ([2]int, error) {
	var result [2]int
	err := v.ep.Call("window_get_cursor", &result, window)
	return result, err
}

func (v *Vim) WindowHeight(window Window) (int, error) {
	var result int
	err := v.ep.Call("window_get_height", &result, window)
	return result, err
}

func (v *Vim) WindowOption(window Window, name string, result interface{}) error {
	return v.ep.Call("window_get_option", result, window, name)
}

func (v *Vim) WindowPosition(window Window) ([2]int, error) {
	var result [2]int
	err := v.ep.Call("window_get_position", &result, window)
	return result, err
}

func (v *Vim) WindowTabpage(window Window) (Tabpage, error) {
	var result Tabpage
	err := v.ep.Call("window_get_tabpage", &result, window)
	return result, err
}

func (v *Vim) WindowVar(window Window, name string, result interface{}) error {
	return v.ep.Call("window_get_var", result, window, name)
}

func (v *Vim) WindowWidth(window Window) (int, error) {
	var result int
	err := v.ep.Call("window_get_width", &result, window)
	return result, err
}

func (v *Vim) IsWindowValid(window Window) (bool, error) {
	var result bool
	err := v.ep.Call("window_is_valid", &result, window)
	return result, err
}

func (v *Vim) SetWindowCursor(window Window, pos [2]int) error {
	return v.ep.Call("window_set_cursor", nil, window, pos)
}

func (v *Vim) SetWindowHeight(window Window, height int) error {
	return v.ep.Call("window_set_height", nil, window, height)
}

func (v *Vim) SetWindowOption(window Window, name string, value interface{}) error {
	return v.ep.Call("window_set_option", nil, window, name, value)
}

func (v *Vim) SetWindowVar(window Window, name string, value interface{}, result interface{}) error {
	return v.ep.Call("window_set_var", result, window, name, value)
}

func (v *Vim) SetWindowWidth(window Window, width int) error {
	return v.ep.Call("window_set_width", nil, window, width)
}
