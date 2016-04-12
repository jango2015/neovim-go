// Code generated by 'go generate'

package vim

import (
	"fmt"
	"reflect"

	"github.com/garyburd/neovim-go/msgpack"
	"github.com/garyburd/neovim-go/msgpack/rpc"
)

const (
	exceptionError  = 0
	validationError = 1
)

func withExtensions() rpc.Option {
	return rpc.WithExtensions(msgpack.ExtensionMap{

		0: func(p []byte) (interface{}, error) {
			x, err := decodeExt(p)
			return Buffer(x), err
		},

		1: func(p []byte) (interface{}, error) {
			x, err := decodeExt(p)
			return Window(x), err
		},

		2: func(p []byte) (interface{}, error) {
			x, err := decodeExt(p)
			return Tabpage(x), err
		},
	})
}

// Buffer represents a remote Neovim buffer.
type Buffer int

func (x *Buffer) UnmarshalMsgPack(dec *msgpack.Decoder) error {
	if dec.Type() != msgpack.Extension || dec.Extension() != 0 {
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
	return enc.PackExtension(0, encodeExt(int(x)))
}

func (x Buffer) String() string {
	return fmt.Sprintf("Buffer:%d", int(x))
}

// Window represents a remote Neovim window.
type Window int

func (x *Window) UnmarshalMsgPack(dec *msgpack.Decoder) error {
	if dec.Type() != msgpack.Extension || dec.Extension() != 1 {
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
	return enc.PackExtension(1, encodeExt(int(x)))
}

func (x Window) String() string {
	return fmt.Sprintf("Window:%d", int(x))
}

// Tabpage represents a remote Neovim tabpage.
type Tabpage int

func (x *Tabpage) UnmarshalMsgPack(dec *msgpack.Decoder) error {
	if dec.Type() != msgpack.Extension || dec.Extension() != 2 {
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
	return enc.PackExtension(2, encodeExt(int(x)))
}

func (x Tabpage) String() string {
	return fmt.Sprintf("Tabpage:%d", int(x))
}

// BufferLineCount returns the number of lines in the buffer.
func (v *Vim) BufferLineCount(buffer Buffer) (int, error) {
	var result int
	err := v.call("buffer_line_count", &result, buffer)
	return result, err
}

// BufferLineCount returns the number of lines in the buffer.
func (p *Pipeline) BufferLineCount(buffer Buffer, result *int) {
	p.call("buffer_line_count", result, buffer)
}

// BufferLines retrieves a line range from a buffer.
//
// Indexing is zero-based, end-exclusive. Negative indices are interpreted as
// length+1+index, i e -1 refers to the index past the end. So to get the last
// element set start=-2 and end=-1.
//
// Out-of-bounds indices are clamped to the nearest valid value, unless strict
// = true.
func (v *Vim) BufferLines(buffer Buffer, start int, end int, strict bool) ([][]byte, error) {
	var result [][]byte
	err := v.call("buffer_get_lines", &result, buffer, start, end, strict)
	return result, err
}

// BufferLines retrieves a line range from a buffer.
//
// Indexing is zero-based, end-exclusive. Negative indices are interpreted as
// length+1+index, i e -1 refers to the index past the end. So to get the last
// element set start=-2 and end=-1.
//
// Out-of-bounds indices are clamped to the nearest valid value, unless strict
// = true.
func (p *Pipeline) BufferLines(buffer Buffer, start int, end int, strict bool, result *[][]byte) {
	p.call("buffer_get_lines", result, buffer, start, end, strict)
}

// SetBufferLines replaces a line range on a buffer.
//
// Indexing is zero-based, end-exclusive. Negative indices are interpreted as
// length+1+index, ie -1 refers to the index past the end. So to change or
// delete the last element set start=-2 and end=-1.
//
// To insert lines at a given index, set both start and end to the same index.
// To delete a range of lines, set replacement to an empty array.
//
// Out-of-bounds indices are clamped to the nearest valid value, unless strict
// = true.
func (v *Vim) SetBufferLines(buffer Buffer, start int, end int, strict bool, replacement [][]byte) error {
	return v.call("buffer_set_lines", nil, buffer, start, end, strict, replacement)
}

// SetBufferLines replaces a line range on a buffer.
//
// Indexing is zero-based, end-exclusive. Negative indices are interpreted as
// length+1+index, ie -1 refers to the index past the end. So to change or
// delete the last element set start=-2 and end=-1.
//
// To insert lines at a given index, set both start and end to the same index.
// To delete a range of lines, set replacement to an empty array.
//
// Out-of-bounds indices are clamped to the nearest valid value, unless strict
// = true.
func (p *Pipeline) SetBufferLines(buffer Buffer, start int, end int, strict bool, replacement [][]byte) {
	p.call("buffer_set_lines", nil, buffer, start, end, strict, replacement)
}

// BufferVar gets a buffer-scoped (b:) variable.
func (v *Vim) BufferVar(buffer Buffer, name string, result interface{}) error {
	return v.call("buffer_get_var", result, buffer, name)
}

// BufferVar gets a buffer-scoped (b:) variable.
func (p *Pipeline) BufferVar(buffer Buffer, name string, result interface{}) {
	p.call("buffer_get_var", result, buffer, name)
}

// SetBufferVar sets a buffer-scoped (b:) variable. The value nil deletes the
// variable. Result is the previous value of the variable.
func (v *Vim) SetBufferVar(buffer Buffer, name string, value interface{}, result interface{}) error {
	return v.call("buffer_set_var", result, buffer, name, value)
}

// SetBufferVar sets a buffer-scoped (b:) variable. The value nil deletes the
// variable. Result is the previous value of the variable.
func (p *Pipeline) SetBufferVar(buffer Buffer, name string, value interface{}, result interface{}) {
	p.call("buffer_set_var", result, buffer, name, value)
}

// BufferOption gets a buffer option value.
func (v *Vim) BufferOption(buffer Buffer, name string, result interface{}) error {
	return v.call("buffer_get_option", result, buffer, name)
}

// BufferOption gets a buffer option value.
func (p *Pipeline) BufferOption(buffer Buffer, name string, result interface{}) {
	p.call("buffer_get_option", result, buffer, name)
}

// SetBufferOption sets a buffer option value. The value nil deletes the option
// in the case where there's a global fallback.
func (v *Vim) SetBufferOption(buffer Buffer, name string, value interface{}) error {
	return v.call("buffer_set_option", nil, buffer, name, value)
}

// SetBufferOption sets a buffer option value. The value nil deletes the option
// in the case where there's a global fallback.
func (p *Pipeline) SetBufferOption(buffer Buffer, name string, value interface{}) {
	p.call("buffer_set_option", nil, buffer, name, value)
}

// BufferNumber gets a buffer's number.
func (v *Vim) BufferNumber(buffer Buffer) (int, error) {
	var result int
	err := v.call("buffer_get_number", &result, buffer)
	return result, err
}

// BufferNumber gets a buffer's number.
func (p *Pipeline) BufferNumber(buffer Buffer, result *int) {
	p.call("buffer_get_number", result, buffer)
}

// BufferName gets the full file name of a buffer.
func (v *Vim) BufferName(buffer Buffer) (string, error) {
	var result string
	err := v.call("buffer_get_name", &result, buffer)
	return result, err
}

// BufferName gets the full file name of a buffer.
func (p *Pipeline) BufferName(buffer Buffer, result *string) {
	p.call("buffer_get_name", result, buffer)
}

// SetBufferName sets the full file name of a buffer.
// BufFilePre/BufFilePost are triggered.
func (v *Vim) SetBufferName(buffer Buffer, name string) error {
	return v.call("buffer_set_name", nil, buffer, name)
}

// SetBufferName sets the full file name of a buffer.
// BufFilePre/BufFilePost are triggered.
func (p *Pipeline) SetBufferName(buffer Buffer, name string) {
	p.call("buffer_set_name", nil, buffer, name)
}

// IsBufferValid returns true if the buffer is valid.
func (v *Vim) IsBufferValid(buffer Buffer) (bool, error) {
	var result bool
	err := v.call("buffer_is_valid", &result, buffer)
	return result, err
}

// IsBufferValid returns true if the buffer is valid.
func (p *Pipeline) IsBufferValid(buffer Buffer, result *bool) {
	p.call("buffer_is_valid", result, buffer)
}

// BufferMark returns the (row,col) of the named mark.
func (v *Vim) BufferMark(buffer Buffer, name string) ([2]int, error) {
	var result [2]int
	err := v.call("buffer_get_mark", &result, buffer, name)
	return result, err
}

// BufferMark returns the (row,col) of the named mark.
func (p *Pipeline) BufferMark(buffer Buffer, name string, result *[2]int) {
	p.call("buffer_get_mark", result, buffer, name)
}

// AddBufferHighlight adds a highlight to buffer and returns the source id of
// the highlight.
//
// AddBufferHighlight can be used for plugins which dynamically generate
// highlights to a buffer (like a semantic highlighter or linter). The function
// adds a single highlight to a buffer. Unlike matchaddpos() highlights follow
// changes to line numbering (as lines are inserted/removed above the
// highlighted line), like signs and marks do.
//
// The srcID is useful for batch deletion/updating of a set of highlights. When
// called with srcID = 0, an unique source id is generated and returned.
// Succesive calls can pass in it as srcID to add new highlights to the same
// source group. All highlights in the same group can then be cleared with
// ClearBufferHighlight. If the highlight never will be manually deleted pass
// in -1 for srcID.
//
// If hlGroup is the empty string no highlight is added, but a new srcID is
// still returned. This is useful for an external plugin to synchrounously
// request an unique srcID at initialization, and later asynchronously add and
// clear highlights in response to buffer changes.
//
// The startCol and endCol parameters specify the range of columns to
// highlight. Use endCol = -1 to highlight to the end of the line.
func (v *Vim) AddBufferHighlight(buffer Buffer, srcID int, hlGroup string, line int, startCol int, endCol int) (int, error) {
	var result int
	err := v.call("buffer_add_highlight", &result, buffer, srcID, hlGroup, line, startCol, endCol)
	return result, err
}

// AddBufferHighlight adds a highlight to buffer and returns the source id of
// the highlight.
//
// AddBufferHighlight can be used for plugins which dynamically generate
// highlights to a buffer (like a semantic highlighter or linter). The function
// adds a single highlight to a buffer. Unlike matchaddpos() highlights follow
// changes to line numbering (as lines are inserted/removed above the
// highlighted line), like signs and marks do.
//
// The srcID is useful for batch deletion/updating of a set of highlights. When
// called with srcID = 0, an unique source id is generated and returned.
// Succesive calls can pass in it as srcID to add new highlights to the same
// source group. All highlights in the same group can then be cleared with
// ClearBufferHighlight. If the highlight never will be manually deleted pass
// in -1 for srcID.
//
// If hlGroup is the empty string no highlight is added, but a new srcID is
// still returned. This is useful for an external plugin to synchrounously
// request an unique srcID at initialization, and later asynchronously add and
// clear highlights in response to buffer changes.
//
// The startCol and endCol parameters specify the range of columns to
// highlight. Use endCol = -1 to highlight to the end of the line.
func (p *Pipeline) AddBufferHighlight(buffer Buffer, srcID int, hlGroup string, line int, startCol int, endCol int, result *int) {
	p.call("buffer_add_highlight", result, buffer, srcID, hlGroup, line, startCol, endCol)
}

// ClearBufferHighlight clears highlights from a given source group and a range
// of lines.
//
// To clear a source group in the entire buffer, pass in 1 and -1 to startLine
// and endLine respectively.
//
// The lineStart and lineEnd parameters specify the range of lines to clear.
// The end of range is exclusive. Specify -1 to clear to the end of the file.
func (v *Vim) ClearBufferHighlight(buffer Buffer, srcID int, startLine int, endLine int) error {
	return v.call(" buffer_clear_highlight", nil, buffer, srcID, startLine, endLine)
}

// ClearBufferHighlight clears highlights from a given source group and a range
// of lines.
//
// To clear a source group in the entire buffer, pass in 1 and -1 to startLine
// and endLine respectively.
//
// The lineStart and lineEnd parameters specify the range of lines to clear.
// The end of range is exclusive. Specify -1 to clear to the end of the file.
func (p *Pipeline) ClearBufferHighlight(buffer Buffer, srcID int, startLine int, endLine int) {
	p.call(" buffer_clear_highlight", nil, buffer, srcID, startLine, endLine)
}

// TabpageWindows returns the windows in a tabpage.
func (v *Vim) TabpageWindows(tabpage Tabpage) ([]Window, error) {
	var result []Window
	err := v.call("tabpage_get_windows", &result, tabpage)
	return result, err
}

// TabpageWindows returns the windows in a tabpage.
func (p *Pipeline) TabpageWindows(tabpage Tabpage, result *[]Window) {
	p.call("tabpage_get_windows", result, tabpage)
}

// TabpageVar gets a tab-scoped (t:) variable.
func (v *Vim) TabpageVar(tabpage Tabpage, name string, result interface{}) error {
	return v.call("tabpage_get_var", result, tabpage, name)
}

// TabpageVar gets a tab-scoped (t:) variable.
func (p *Pipeline) TabpageVar(tabpage Tabpage, name string, result interface{}) {
	p.call("tabpage_get_var", result, tabpage, name)
}

// SetTabpageVar isets a tab-scoped (t:) variable. A nil value deletes the variable.
func (v *Vim) SetTabpageVar(tabpage Tabpage, name string, value interface{}, result interface{}) error {
	return v.call("tabpage_set_var", result, tabpage, name, value)
}

// SetTabpageVar isets a tab-scoped (t:) variable. A nil value deletes the variable.
func (p *Pipeline) SetTabpageVar(tabpage Tabpage, name string, value interface{}, result interface{}) {
	p.call("tabpage_set_var", result, tabpage, name, value)
}

// TabpageWindow gets the current window in a tab page.
func (v *Vim) TabpageWindow(tabpage Tabpage) (Window, error) {
	var result Window
	err := v.call("tabpage_get_window", &result, tabpage)
	return result, err
}

// TabpageWindow gets the current window in a tab page.
func (p *Pipeline) TabpageWindow(tabpage Tabpage, result *Window) {
	p.call("tabpage_get_window", result, tabpage)
}

// IsTabpageValid checks if a tab page is valid.
func (v *Vim) IsTabpageValid(tabpage Tabpage) (bool, error) {
	var result bool
	err := v.call("tabpage_is_valid", &result, tabpage)
	return result, err
}

// IsTabpageValid checks if a tab page is valid.
func (p *Pipeline) IsTabpageValid(tabpage Tabpage, result *bool) {
	p.call("tabpage_is_valid", result, tabpage)
}

// Command executes a single ex command.
func (v *Vim) Command(str string) error {
	return v.call("vim_command", nil, str)
}

// Command executes a single ex command.
func (p *Pipeline) Command(str string) {
	p.call("vim_command", nil, str)
}

// FeedKeys Pushes keys to the Neovim user input buffer. Options can be a string
// with the following character flags:
//
//  m:  Remap keys. This is default.
//  n:  Do not remap keys.
//  t:  Handle keys as if typed; otherwise they are handled as if coming from a
//     mapping. This matters for undo, opening folds, etc.
func (v *Vim) FeedKeys(keys string, mode string, escapeCsi bool) error {
	return v.call("vim_feedkeys", nil, keys, mode, escapeCsi)
}

// FeedKeys Pushes keys to the Neovim user input buffer. Options can be a string
// with the following character flags:
//
//  m:  Remap keys. This is default.
//  n:  Do not remap keys.
//  t:  Handle keys as if typed; otherwise they are handled as if coming from a
//     mapping. This matters for undo, opening folds, etc.
func (p *Pipeline) FeedKeys(keys string, mode string, escapeCsi bool) {
	p.call("vim_feedkeys", nil, keys, mode, escapeCsi)
}

// Input pushes bytes to the Neovim low level input buffer.
//
// Unlike FeedKeys, this uses the lowest level input buffer and the call is not
// deferred. It returns the number of bytes actually written(which can be less
// than what was requested if the buffer is full).
func (v *Vim) Input(keys string) (int, error) {
	var result int
	err := v.call("vim_input", &result, keys)
	return result, err
}

// Input pushes bytes to the Neovim low level input buffer.
//
// Unlike FeedKeys, this uses the lowest level input buffer and the call is not
// deferred. It returns the number of bytes actually written(which can be less
// than what was requested if the buffer is full).
func (p *Pipeline) Input(keys string, result *int) {
	p.call("vim_input", result, keys)
}

// ReplaceTermcodes replaces any terminal code strings by byte sequences. The
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
	err := v.call("vim_replace_termcodes", &result, str, fromPart, doLt, special)
	return result, err
}

// ReplaceTermcodes replaces any terminal code strings by byte sequences. The
// returned sequences are Nvim's internal representation of keys, for example:
//
//  <esc> -> '\x1b'
//  <cr>  -> '\r'
//  <c-l> -> '\x0c'
//  <up>  -> '\x80ku'
//
// The returned sequences can be used as input to feedkeys.
func (p *Pipeline) ReplaceTermcodes(str string, fromPart bool, doLt bool, special bool, result *string) {
	p.call("vim_replace_termcodes", result, str, fromPart, doLt, special)
}

// CommandOutput executes a single ex command and returns the output.
func (v *Vim) CommandOutput(str string) (string, error) {
	var result string
	err := v.call("vim_command_output", &result, str)
	return result, err
}

// CommandOutput executes a single ex command and returns the output.
func (p *Pipeline) CommandOutput(str string, result *string) {
	p.call("vim_command_output", result, str)
}

// Eval evaluates the expression str using the Vim internal expression
// evaluator.
//
//  :help expression
func (v *Vim) Eval(str string, result interface{}) error {
	return v.call("vim_eval", result, str)
}

// Eval evaluates the expression str using the Vim internal expression
// evaluator.
//
//  :help expression
func (p *Pipeline) Eval(str string, result interface{}) {
	p.call("vim_eval", result, str)
}

// Strwidth returns the number of display cells the string occupies. Tab is
// counted as one cell.
func (v *Vim) Strwidth(str string) (int, error) {
	var result int
	err := v.call("vim_strwidth", &result, str)
	return result, err
}

// Strwidth returns the number of display cells the string occupies. Tab is
// counted as one cell.
func (p *Pipeline) Strwidth(str string, result *int) {
	p.call("vim_strwidth", result, str)
}

// RuntimePaths returns a list of paths contained in the runtimepath option.
func (v *Vim) RuntimePaths() ([]string, error) {
	var result []string
	err := v.call("vim_list_runtime_paths", &result)
	return result, err
}

// RuntimePaths returns a list of paths contained in the runtimepath option.
func (p *Pipeline) RuntimePaths(result *[]string) {
	p.call("vim_list_runtime_paths", result)
}

// ChangeDirectory changes Vim working directory.
func (v *Vim) ChangeDirectory(dir string) error {
	return v.call("vim_change_directory", nil, dir)
}

// ChangeDirectory changes Vim working directory.
func (p *Pipeline) ChangeDirectory(dir string) {
	p.call("vim_change_directory", nil, dir)
}

// CurrentLine gets the current line in the current buffer.
func (v *Vim) CurrentLine() ([]byte, error) {
	var result []byte
	err := v.call("vim_get_current_line", &result)
	return result, err
}

// CurrentLine gets the current line in the current buffer.
func (p *Pipeline) CurrentLine(result *[]byte) {
	p.call("vim_get_current_line", result)
}

// SetCurrentLine sets the current line in the current buffer.
func (v *Vim) SetCurrentLine(line []byte) error {
	return v.call("vim_set_current_line", nil, line)
}

// SetCurrentLine sets the current line in the current buffer.
func (p *Pipeline) SetCurrentLine(line []byte) {
	p.call("vim_set_current_line", nil, line)
}

// DeleteCurrentLine deletes the current line in the current buffer.
func (v *Vim) DeleteCurrentLine() error {
	return v.call("vim_del_current_line", nil)
}

// DeleteCurrentLine deletes the current line in the current buffer.
func (p *Pipeline) DeleteCurrentLine() {
	p.call("vim_del_current_line", nil)
}

// Var gets a global (g:) variable.
func (v *Vim) Var(name string, result interface{}) error {
	return v.call("vim_get_var", result, name)
}

// Var gets a global (g:) variable.
func (p *Pipeline) Var(name string, result interface{}) {
	p.call("vim_get_var", result, name)
}

// SetVar sets a global (g:) variable. The value nil deletes the variable.
// Result is the previous value of the variable.
func (v *Vim) SetVar(name string, value interface{}, result interface{}) error {
	return v.call("vim_set_var", result, name, value)
}

// SetVar sets a global (g:) variable. The value nil deletes the variable.
// Result is the previous value of the variable.
func (p *Pipeline) SetVar(name string, value interface{}, result interface{}) {
	p.call("vim_set_var", result, name, value)
}

// Vvar gets a vim (v:) variable.
func (v *Vim) Vvar(name string, result interface{}) error {
	return v.call("vim_get_vvar", result, name)
}

// Vvar gets a vim (v:) variable.
func (p *Pipeline) Vvar(name string, result interface{}) {
	p.call("vim_get_vvar", result, name)
}

// Option gets an option.
func (v *Vim) Option(name string, result interface{}) error {
	return v.call("vim_get_option", result, name)
}

// Option gets an option.
func (p *Pipeline) Option(name string, result interface{}) {
	p.call("vim_get_option", result, name)
}

// SetOption sets an option.
func (v *Vim) SetOption(name string, value interface{}) error {
	return v.call("vim_set_option", nil, name, value)
}

// SetOption sets an option.
func (p *Pipeline) SetOption(name string, value interface{}) {
	p.call("vim_set_option", nil, name, value)
}

// WriteOut writes a message to vim output buffer. The string is split and
// flushed after each newline. Incomplete lines are kept for writing later.
func (v *Vim) WriteOut(str string) error {
	return v.call("vim_out_write", nil, str)
}

// WriteOut writes a message to vim output buffer. The string is split and
// flushed after each newline. Incomplete lines are kept for writing later.
func (p *Pipeline) WriteOut(str string) {
	p.call("vim_out_write", nil, str)
}

// WriteErr writes a message to vim error buffer. The string is split and
// flushed after each newline. Incomplete lines are kept for writing later.
func (v *Vim) WriteErr(str string) error {
	return v.call("vim_err_write", nil, str)
}

// WriteErr writes a message to vim error buffer. The string is split and
// flushed after each newline. Incomplete lines are kept for writing later.
func (p *Pipeline) WriteErr(str string) {
	p.call("vim_err_write", nil, str)
}

// ReportError writes prints str and a newline as an error message.
func (v *Vim) ReportError(str string) error {
	return v.call("vim_report_error", nil, str)
}

// ReportError writes prints str and a newline as an error message.
func (p *Pipeline) ReportError(str string) {
	p.call("vim_report_error", nil, str)
}

// Buffers returns the current list of buffers.
func (v *Vim) Buffers() ([]Buffer, error) {
	var result []Buffer
	err := v.call("vim_get_buffers", &result)
	return result, err
}

// Buffers returns the current list of buffers.
func (p *Pipeline) Buffers(result *[]Buffer) {
	p.call("vim_get_buffers", result)
}

// CurrentBuffer returns the current buffer.
func (v *Vim) CurrentBuffer() (Buffer, error) {
	var result Buffer
	err := v.call("vim_get_current_buffer", &result)
	return result, err
}

// CurrentBuffer returns the current buffer.
func (p *Pipeline) CurrentBuffer(result *Buffer) {
	p.call("vim_get_current_buffer", result)
}

// SetCurrentBuffer sets the current buffer.
func (v *Vim) SetCurrentBuffer(buffer Buffer) error {
	return v.call("vim_set_current_buffer", nil, buffer)
}

// SetCurrentBuffer sets the current buffer.
func (p *Pipeline) SetCurrentBuffer(buffer Buffer) {
	p.call("vim_set_current_buffer", nil, buffer)
}

// Windows returns the current list of windows.
func (v *Vim) Windows() ([]Window, error) {
	var result []Window
	err := v.call("vim_get_windows", &result)
	return result, err
}

// Windows returns the current list of windows.
func (p *Pipeline) Windows(result *[]Window) {
	p.call("vim_get_windows", result)
}

// CurrentWindow returns the current window.
func (v *Vim) CurrentWindow() (Window, error) {
	var result Window
	err := v.call("vim_get_current_window", &result)
	return result, err
}

// CurrentWindow returns the current window.
func (p *Pipeline) CurrentWindow(result *Window) {
	p.call("vim_get_current_window", result)
}

// SetCurrentWindow sets the current window.
func (v *Vim) SetCurrentWindow(window Window) error {
	return v.call("vim_set_current_window", nil, window)
}

// SetCurrentWindow sets the current window.
func (p *Pipeline) SetCurrentWindow(window Window) {
	p.call("vim_set_current_window", nil, window)
}

// Tabpages returns the current list of tabpages.
func (v *Vim) Tabpages() ([]Tabpage, error) {
	var result []Tabpage
	err := v.call("vim_get_tabpages", &result)
	return result, err
}

// Tabpages returns the current list of tabpages.
func (p *Pipeline) Tabpages(result *[]Tabpage) {
	p.call("vim_get_tabpages", result)
}

// CurrentTabpage returns the current tabpage.
func (v *Vim) CurrentTabpage() (Tabpage, error) {
	var result Tabpage
	err := v.call("vim_get_current_tabpage", &result)
	return result, err
}

// CurrentTabpage returns the current tabpage.
func (p *Pipeline) CurrentTabpage(result *Tabpage) {
	p.call("vim_get_current_tabpage", result)
}

// SetCurrentTabpage sets the current tabpage.
func (v *Vim) SetCurrentTabpage(tabpage Tabpage) error {
	return v.call("vim_set_current_tabpage", nil, tabpage)
}

// SetCurrentTabpage sets the current tabpage.
func (p *Pipeline) SetCurrentTabpage(tabpage Tabpage) {
	p.call("vim_set_current_tabpage", nil, tabpage)
}

// Subscribe subscribes to a Neovim event.
func (v *Vim) Subscribe(event string) error {
	return v.call("vim_subscribe", nil, event)
}

// Subscribe subscribes to a Neovim event.
func (p *Pipeline) Subscribe(event string) {
	p.call("vim_subscribe", nil, event)
}

// Unsubscribe unsubscribes to a Neovim event.
func (v *Vim) Unsubscribe(event string) error {
	return v.call("vim_unsubscribe", nil, event)
}

// Unsubscribe unsubscribes to a Neovim event.
func (p *Pipeline) Unsubscribe(event string) {
	p.call("vim_unsubscribe", nil, event)
}

func (v *Vim) NameToColor(name string) (int, error) {
	var result int
	err := v.call("vim_name_to_color", &result, name)
	return result, err
}

func (p *Pipeline) NameToColor(name string, result *int) {
	p.call("vim_name_to_color", result, name)
}

func (v *Vim) ColorMap() (map[string]interface{}, error) {
	var result map[string]interface{}
	err := v.call("vim_get_color_map", &result)
	return result, err
}

func (p *Pipeline) ColorMap(result *map[string]interface{}) {
	p.call("vim_get_color_map", result)
}

func (v *Vim) APIInfo() ([]interface{}, error) {
	var result []interface{}
	err := v.call("vim_get_api_info", &result)
	return result, err
}

func (p *Pipeline) APIInfo(result *[]interface{}) {
	p.call("vim_get_api_info", result)
}

// WindowBuffer returns the current buffer in a window.
func (v *Vim) WindowBuffer(window Window) (Buffer, error) {
	var result Buffer
	err := v.call("window_get_buffer", &result, window)
	return result, err
}

// WindowBuffer returns the current buffer in a window.
func (p *Pipeline) WindowBuffer(window Window, result *Buffer) {
	p.call("window_get_buffer", result, window)
}

// WindowCursor returns the cursor position in the window.
func (v *Vim) WindowCursor(window Window) ([2]int, error) {
	var result [2]int
	err := v.call("window_get_cursor", &result, window)
	return result, err
}

// WindowCursor returns the cursor position in the window.
func (p *Pipeline) WindowCursor(window Window, result *[2]int) {
	p.call("window_get_cursor", result, window)
}

// SetWindowCursor sets the cursor position in the window to the given position.
func (v *Vim) SetWindowCursor(window Window, pos [2]int) error {
	return v.call("window_set_cursor", nil, window, pos)
}

// SetWindowCursor sets the cursor position in the window to the given position.
func (p *Pipeline) SetWindowCursor(window Window, pos [2]int) {
	p.call("window_set_cursor", nil, window, pos)
}

// WindowHeight returns the window height.
func (v *Vim) WindowHeight(window Window) (int, error) {
	var result int
	err := v.call("window_get_height", &result, window)
	return result, err
}

// WindowHeight returns the window height.
func (p *Pipeline) WindowHeight(window Window, result *int) {
	p.call("window_get_height", result, window)
}

// SetWindowHeight sets the window height.
func (v *Vim) SetWindowHeight(window Window, height int) error {
	return v.call("window_set_height", nil, window, height)
}

// SetWindowHeight sets the window height.
func (p *Pipeline) SetWindowHeight(window Window, height int) {
	p.call("window_set_height", nil, window, height)
}

// WindowWidth returns the window width.
func (v *Vim) WindowWidth(window Window) (int, error) {
	var result int
	err := v.call("window_get_width", &result, window)
	return result, err
}

// WindowWidth returns the window width.
func (p *Pipeline) WindowWidth(window Window, result *int) {
	p.call("window_get_width", result, window)
}

// SetWindowWidth sets the window width.
func (v *Vim) SetWindowWidth(window Window, width int) error {
	return v.call("window_set_width", nil, window, width)
}

// SetWindowWidth sets the window width.
func (p *Pipeline) SetWindowWidth(window Window, width int) {
	p.call("window_set_width", nil, window, width)
}

// WindowVar gets a window-scoped (w:) variable.
func (v *Vim) WindowVar(window Window, name string, result interface{}) error {
	return v.call("window_get_var", result, window, name)
}

// WindowVar gets a window-scoped (w:) variable.
func (p *Pipeline) WindowVar(window Window, name string, result interface{}) {
	p.call("window_get_var", result, window, name)
}

// SetWindowVar sets a window-scoped (w:) variable.
func (v *Vim) SetWindowVar(window Window, name string, value interface{}, result interface{}) error {
	return v.call("window_set_var", result, window, name, value)
}

// SetWindowVar sets a window-scoped (w:) variable.
func (p *Pipeline) SetWindowVar(window Window, name string, value interface{}, result interface{}) {
	p.call("window_set_var", result, window, name, value)
}

// WindowOption gets a window option.
func (v *Vim) WindowOption(window Window, name string, result interface{}) error {
	return v.call("window_get_option", result, window, name)
}

// WindowOption gets a window option.
func (p *Pipeline) WindowOption(window Window, name string, result interface{}) {
	p.call("window_get_option", result, window, name)
}

// SetWindowOption sets a window option.
func (v *Vim) SetWindowOption(window Window, name string, value interface{}) error {
	return v.call("window_set_option", nil, window, name, value)
}

// SetWindowOption sets a window option.
func (p *Pipeline) SetWindowOption(window Window, name string, value interface{}) {
	p.call("window_set_option", nil, window, name, value)
}

// WindowPosition gets the window position in display cells. First position is zero.
func (v *Vim) WindowPosition(window Window) ([2]int, error) {
	var result [2]int
	err := v.call("window_get_position", &result, window)
	return result, err
}

// WindowPosition gets the window position in display cells. First position is zero.
func (p *Pipeline) WindowPosition(window Window, result *[2]int) {
	p.call("window_get_position", result, window)
}

// WindowTabpage gets the tab page that contains the window.
func (v *Vim) WindowTabpage(window Window) (Tabpage, error) {
	var result Tabpage
	err := v.call("window_get_tabpage", &result, window)
	return result, err
}

// WindowTabpage gets the tab page that contains the window.
func (p *Pipeline) WindowTabpage(window Window, result *Tabpage) {
	p.call("window_get_tabpage", result, window)
}

// IsWindowValid returns true if the window is valid.
func (v *Vim) IsWindowValid(window Window) (bool, error) {
	var result bool
	err := v.call("window_is_valid", &result, window)
	return result, err
}

// IsWindowValid returns true if the window is valid.
func (p *Pipeline) IsWindowValid(window Window, result *bool) {
	p.call("window_is_valid", result, window)
}
