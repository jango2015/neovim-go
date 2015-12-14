// Code generated by 'go generate'

package vim

func (v *Vim) BufferLineCount(buffer Buffer) (int, error) {
	var result int
	err := v.call("buffer_line_count", &result, buffer)
	return result, err
}

func (b *Batch) BufferLineCount(buffer Buffer, result *int) {
	b.call("buffer_line_count", result, buffer)
}

func (v *Vim) BufferLine(buffer Buffer, index int) ([]byte, error) {
	var result []byte
	err := v.call("buffer_get_line", &result, buffer, index)
	return result, err
}

func (b *Batch) BufferLine(buffer Buffer, index int, result *[]byte) {
	b.call("buffer_get_line", result, buffer, index)
}

func (v *Vim) SetBufferLine(buffer Buffer, index int, line []byte) error {
	return v.call("buffer_set_line", nil, buffer, index, line)
}

func (b *Batch) SetBufferLine(buffer Buffer, index int, line []byte) {
	b.call("buffer_set_line", nil, buffer, index, line)
}

func (v *Vim) DeleteBufferLine(buffer Buffer, index int) error {
	return v.call("buffer_del_line", nil, buffer, index)
}

func (b *Batch) DeleteBufferLine(buffer Buffer, index int) {
	b.call("buffer_del_line", nil, buffer, index)
}

func (v *Vim) BufferLineSlice(buffer Buffer, start int, end int, includeStart bool, includeEnd bool) ([][]byte, error) {
	var result [][]byte
	err := v.call("buffer_get_line_slice", &result, buffer, start, end, includeStart, includeEnd)
	return result, err
}

func (b *Batch) BufferLineSlice(buffer Buffer, start int, end int, includeStart bool, includeEnd bool, result *[][]byte) {
	b.call("buffer_get_line_slice", result, buffer, start, end, includeStart, includeEnd)
}

func (v *Vim) SetBufferLineSlice(buffer Buffer, start int, end int, includeStart bool, includeEnd bool, replacement [][]byte) error {
	return v.call("buffer_set_line_slice", nil, buffer, start, end, includeStart, includeEnd, replacement)
}

func (b *Batch) SetBufferLineSlice(buffer Buffer, start int, end int, includeStart bool, includeEnd bool, replacement [][]byte) {
	b.call("buffer_set_line_slice", nil, buffer, start, end, includeStart, includeEnd, replacement)
}

func (v *Vim) BufferVar(buffer Buffer, name string, result interface{}) error {
	return v.call("buffer_get_var", result, buffer, name)
}

func (b *Batch) BufferVar(buffer Buffer, name string, result interface{}) {
	b.call("buffer_get_var", result, buffer, name)
}

func (v *Vim) SetBufferVar(buffer Buffer, name string, value interface{}, result interface{}) error {
	return v.call("buffer_set_var", result, buffer, name, value)
}

func (b *Batch) SetBufferVar(buffer Buffer, name string, value interface{}, result interface{}) {
	b.call("buffer_set_var", result, buffer, name, value)
}

func (v *Vim) BufferOption(buffer Buffer, name string, result interface{}) error {
	return v.call("buffer_get_option", result, buffer, name)
}

func (b *Batch) BufferOption(buffer Buffer, name string, result interface{}) {
	b.call("buffer_get_option", result, buffer, name)
}

func (v *Vim) SetBufferOption(buffer Buffer, name string, value interface{}) error {
	return v.call("buffer_set_option", nil, buffer, name, value)
}

func (b *Batch) SetBufferOption(buffer Buffer, name string, value interface{}) {
	b.call("buffer_set_option", nil, buffer, name, value)
}

func (v *Vim) BufferNumber(buffer Buffer) (int, error) {
	var result int
	err := v.call("buffer_get_number", &result, buffer)
	return result, err
}

func (b *Batch) BufferNumber(buffer Buffer, result *int) {
	b.call("buffer_get_number", result, buffer)
}

func (v *Vim) BufferName(buffer Buffer) (string, error) {
	var result string
	err := v.call("buffer_get_name", &result, buffer)
	return result, err
}

func (b *Batch) BufferName(buffer Buffer, result *string) {
	b.call("buffer_get_name", result, buffer)
}

func (v *Vim) SetBufferName(buffer Buffer, name string) error {
	return v.call("buffer_set_name", nil, buffer, name)
}

func (b *Batch) SetBufferName(buffer Buffer, name string) {
	b.call("buffer_set_name", nil, buffer, name)
}

func (v *Vim) IsBufferValid(buffer Buffer) (bool, error) {
	var result bool
	err := v.call("buffer_is_valid", &result, buffer)
	return result, err
}

func (b *Batch) IsBufferValid(buffer Buffer, result *bool) {
	b.call("buffer_is_valid", result, buffer)
}

func (v *Vim) InsertBuffer(buffer Buffer, lnum int, lines [][]byte) error {
	return v.call("buffer_insert", nil, buffer, lnum, lines)
}

func (b *Batch) InsertBuffer(buffer Buffer, lnum int, lines [][]byte) {
	b.call("buffer_insert", nil, buffer, lnum, lines)
}

func (v *Vim) BufferMark(buffer Buffer, name string) ([2]int, error) {
	var result [2]int
	err := v.call("buffer_get_mark", &result, buffer, name)
	return result, err
}

func (b *Batch) BufferMark(buffer Buffer, name string, result *[2]int) {
	b.call("buffer_get_mark", result, buffer, name)
}

func (v *Vim) TabpageWindows(tabpage Tabpage) ([]Window, error) {
	var result []Window
	err := v.call("tabpage_get_windows", &result, tabpage)
	return result, err
}

func (b *Batch) TabpageWindows(tabpage Tabpage, result *[]Window) {
	b.call("tabpage_get_windows", result, tabpage)
}

func (v *Vim) TabpageVar(tabpage Tabpage, name string, result interface{}) error {
	return v.call("tabpage_get_var", result, tabpage, name)
}

func (b *Batch) TabpageVar(tabpage Tabpage, name string, result interface{}) {
	b.call("tabpage_get_var", result, tabpage, name)
}

func (v *Vim) SetTabpageVar(tabpage Tabpage, name string, value interface{}, result interface{}) error {
	return v.call("tabpage_set_var", result, tabpage, name, value)
}

func (b *Batch) SetTabpageVar(tabpage Tabpage, name string, value interface{}, result interface{}) {
	b.call("tabpage_set_var", result, tabpage, name, value)
}

func (v *Vim) TabpageWindow(tabpage Tabpage) (Window, error) {
	var result Window
	err := v.call("tabpage_get_window", &result, tabpage)
	return result, err
}

func (b *Batch) TabpageWindow(tabpage Tabpage, result *Window) {
	b.call("tabpage_get_window", result, tabpage)
}

func (v *Vim) IsTabpageValid(tabpage Tabpage) (bool, error) {
	var result bool
	err := v.call("tabpage_is_valid", &result, tabpage)
	return result, err
}

func (b *Batch) IsTabpageValid(tabpage Tabpage, result *bool) {
	b.call("tabpage_is_valid", result, tabpage)
}

// Command executes a single ex command.
func (v *Vim) Command(str string) error {
	return v.call("vim_command", nil, str)
}

// Command executes a single ex command.
func (b *Batch) Command(str string) {
	b.call("vim_command", nil, str)
}

// FeedKeys Pushes keys to the Neovim user input buffer. Options can be a string
// with the following character flags:
//
//  m:  Remap keys. This is default.
//  n:  Do not remap keys.
//  t:  Handle keys as if typed; otherwise they are handled as if coming from a
//     mapping. This matters for undo, opening folds, etc.
func (v *Vim) Feedkeys(keys string, mode string, escapeCsi bool) error {
	return v.call("vim_feedkeys", nil, keys, mode, escapeCsi)
}

// FeedKeys Pushes keys to the Neovim user input buffer. Options can be a string
// with the following character flags:
//
//  m:  Remap keys. This is default.
//  n:  Do not remap keys.
//  t:  Handle keys as if typed; otherwise they are handled as if coming from a
//     mapping. This matters for undo, opening folds, etc.
func (b *Batch) Feedkeys(keys string, mode string, escapeCsi bool) {
	b.call("vim_feedkeys", nil, keys, mode, escapeCsi)
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
func (b *Batch) Input(keys string, result *int) {
	b.call("vim_input", result, keys)
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
func (b *Batch) ReplaceTermcodes(str string, fromPart bool, doLt bool, special bool, result *string) {
	b.call("vim_replace_termcodes", result, str, fromPart, doLt, special)
}

// CommandOutput executes a single ex command and returns the output.
func (v *Vim) CommandOutput(str string) (string, error) {
	var result string
	err := v.call("vim_command_output", &result, str)
	return result, err
}

// CommandOutput executes a single ex command and returns the output.
func (b *Batch) CommandOutput(str string, result *string) {
	b.call("vim_command_output", result, str)
}

// Eval evaluates a vimscript expression.
func (v *Vim) Eval(str string, result interface{}) error {
	return v.call("vim_eval", result, str)
}

// Eval evaluates a vimscript expression.
func (b *Batch) Eval(str string, result interface{}) {
	b.call("vim_eval", result, str)
}

// Strwidth returns the number of display cells string occupies. Tab is counted
// as one cell.
func (v *Vim) Strwidth(str string) (int, error) {
	var result int
	err := v.call("vim_strwidth", &result, str)
	return result, err
}

// Strwidth returns the number of display cells string occupies. Tab is counted
// as one cell.
func (b *Batch) Strwidth(str string, result *int) {
	b.call("vim_strwidth", result, str)
}

// RuntimePaths returns a list of paths contained in the runtimepath option.
func (v *Vim) RuntimePaths() ([]string, error) {
	var result []string
	err := v.call("vim_list_runtime_paths", &result)
	return result, err
}

// RuntimePaths returns a list of paths contained in the runtimepath option.
func (b *Batch) RuntimePaths(result *[]string) {
	b.call("vim_list_runtime_paths", result)
}

func (v *Vim) ChangeDirectory(dir string) error {
	return v.call("vim_change_directory", nil, dir)
}

func (b *Batch) ChangeDirectory(dir string) {
	b.call("vim_change_directory", nil, dir)
}

func (v *Vim) CurrentLine() ([]byte, error) {
	var result []byte
	err := v.call("vim_get_current_line", &result)
	return result, err
}

func (b *Batch) CurrentLine(result *[]byte) {
	b.call("vim_get_current_line", result)
}

func (v *Vim) SetCurrentLine(line []byte) error {
	return v.call("vim_set_current_line", nil, line)
}

func (b *Batch) SetCurrentLine(line []byte) {
	b.call("vim_set_current_line", nil, line)
}

func (v *Vim) DeleteCurrentLine() error {
	return v.call("vim_del_current_line", nil)
}

func (b *Batch) DeleteCurrentLine() {
	b.call("vim_del_current_line", nil)
}

func (v *Vim) Var(name string, result interface{}) error {
	return v.call("vim_get_var", result, name)
}

func (b *Batch) Var(name string, result interface{}) {
	b.call("vim_get_var", result, name)
}

func (v *Vim) SetVar(name string, value interface{}, result interface{}) error {
	return v.call("vim_set_var", result, name, value)
}

func (b *Batch) SetVar(name string, value interface{}, result interface{}) {
	b.call("vim_set_var", result, name, value)
}

func (v *Vim) Vvar(name string, result interface{}) error {
	return v.call("vim_get_vvar", result, name)
}

func (b *Batch) Vvar(name string, result interface{}) {
	b.call("vim_get_vvar", result, name)
}

func (v *Vim) Option(name string, result interface{}) error {
	return v.call("vim_get_option", result, name)
}

func (b *Batch) Option(name string, result interface{}) {
	b.call("vim_get_option", result, name)
}

func (v *Vim) SetOption(name string, value interface{}) error {
	return v.call("vim_set_option", nil, name, value)
}

func (b *Batch) SetOption(name string, value interface{}) {
	b.call("vim_set_option", nil, name, value)
}

// WriteOut writes a message to the output buffer.
func (v *Vim) WriteOut(str string) error {
	return v.call("vim_out_write", nil, str)
}

// WriteOut writes a message to the output buffer.
func (b *Batch) WriteOut(str string) {
	b.call("vim_out_write", nil, str)
}

// WriteErr writes a message to the error buffer.
func (v *Vim) WriteErr(str string) error {
	return v.call("vim_err_write", nil, str)
}

// WriteErr writes a message to the error buffer.
func (b *Batch) WriteErr(str string) {
	b.call("vim_err_write", nil, str)
}

// ReportError writes a message and a newline to the error buffer.
func (v *Vim) ReportError(str string) error {
	return v.call("vim_report_error", nil, str)
}

// ReportError writes a message and a newline to the error buffer.
func (b *Batch) ReportError(str string) {
	b.call("vim_report_error", nil, str)
}

func (v *Vim) Buffers() ([]Buffer, error) {
	var result []Buffer
	err := v.call("vim_get_buffers", &result)
	return result, err
}

func (b *Batch) Buffers(result *[]Buffer) {
	b.call("vim_get_buffers", result)
}

func (v *Vim) CurrentBuffer() (Buffer, error) {
	var result Buffer
	err := v.call("vim_get_current_buffer", &result)
	return result, err
}

func (b *Batch) CurrentBuffer(result *Buffer) {
	b.call("vim_get_current_buffer", result)
}

func (v *Vim) SetCurrentBuffer(buffer Buffer) error {
	return v.call("vim_set_current_buffer", nil, buffer)
}

func (b *Batch) SetCurrentBuffer(buffer Buffer) {
	b.call("vim_set_current_buffer", nil, buffer)
}

func (v *Vim) Windows() ([]Window, error) {
	var result []Window
	err := v.call("vim_get_windows", &result)
	return result, err
}

func (b *Batch) Windows(result *[]Window) {
	b.call("vim_get_windows", result)
}

func (v *Vim) CurrentWindow() (Window, error) {
	var result Window
	err := v.call("vim_get_current_window", &result)
	return result, err
}

func (b *Batch) CurrentWindow(result *Window) {
	b.call("vim_get_current_window", result)
}

func (v *Vim) SetCurrentWindow(window Window) error {
	return v.call("vim_set_current_window", nil, window)
}

func (b *Batch) SetCurrentWindow(window Window) {
	b.call("vim_set_current_window", nil, window)
}

func (v *Vim) Tabpages() ([]Tabpage, error) {
	var result []Tabpage
	err := v.call("vim_get_tabpages", &result)
	return result, err
}

func (b *Batch) Tabpages(result *[]Tabpage) {
	b.call("vim_get_tabpages", result)
}

func (v *Vim) CurrentTabpage() (Tabpage, error) {
	var result Tabpage
	err := v.call("vim_get_current_tabpage", &result)
	return result, err
}

func (b *Batch) CurrentTabpage(result *Tabpage) {
	b.call("vim_get_current_tabpage", result)
}

func (v *Vim) SetCurrentTabpage(tabpage Tabpage) error {
	return v.call("vim_set_current_tabpage", nil, tabpage)
}

func (b *Batch) SetCurrentTabpage(tabpage Tabpage) {
	b.call("vim_set_current_tabpage", nil, tabpage)
}

// Subscribe subscribes to a Neovim event.
func (v *Vim) Subscribe(event string) error {
	return v.call("vim_subscribe", nil, event)
}

// Subscribe subscribes to a Neovim event.
func (b *Batch) Subscribe(event string) {
	b.call("vim_subscribe", nil, event)
}

// Unsubscribe unsubscribes to a Neovim event.
func (v *Vim) Unsubscribe(event string) error {
	return v.call("vim_unsubscribe", nil, event)
}

// Unsubscribe unsubscribes to a Neovim event.
func (b *Batch) Unsubscribe(event string) {
	b.call("vim_unsubscribe", nil, event)
}

func (v *Vim) NameToColor(name string) (int, error) {
	var result int
	err := v.call("vim_name_to_color", &result, name)
	return result, err
}

func (b *Batch) NameToColor(name string, result *int) {
	b.call("vim_name_to_color", result, name)
}

func (v *Vim) ColorMap() (map[string]interface{}, error) {
	var result map[string]interface{}
	err := v.call("vim_get_color_map", &result)
	return result, err
}

func (b *Batch) ColorMap(result *map[string]interface{}) {
	b.call("vim_get_color_map", result)
}

func (v *Vim) APIInfo() ([]interface{}, error) {
	var result []interface{}
	err := v.call("vim_get_api_info", &result)
	return result, err
}

func (b *Batch) APIInfo(result *[]interface{}) {
	b.call("vim_get_api_info", result)
}

func (v *Vim) WindowBuffer(window Window) (Buffer, error) {
	var result Buffer
	err := v.call("window_get_buffer", &result, window)
	return result, err
}

func (b *Batch) WindowBuffer(window Window, result *Buffer) {
	b.call("window_get_buffer", result, window)
}

func (v *Vim) WindowCursor(window Window) ([2]int, error) {
	var result [2]int
	err := v.call("window_get_cursor", &result, window)
	return result, err
}

func (b *Batch) WindowCursor(window Window, result *[2]int) {
	b.call("window_get_cursor", result, window)
}

func (v *Vim) SetWindowCursor(window Window, pos [2]int) error {
	return v.call("window_set_cursor", nil, window, pos)
}

func (b *Batch) SetWindowCursor(window Window, pos [2]int) {
	b.call("window_set_cursor", nil, window, pos)
}

func (v *Vim) WindowHeight(window Window) (int, error) {
	var result int
	err := v.call("window_get_height", &result, window)
	return result, err
}

func (b *Batch) WindowHeight(window Window, result *int) {
	b.call("window_get_height", result, window)
}

func (v *Vim) SetWindowHeight(window Window, height int) error {
	return v.call("window_set_height", nil, window, height)
}

func (b *Batch) SetWindowHeight(window Window, height int) {
	b.call("window_set_height", nil, window, height)
}

func (v *Vim) WindowWidth(window Window) (int, error) {
	var result int
	err := v.call("window_get_width", &result, window)
	return result, err
}

func (b *Batch) WindowWidth(window Window, result *int) {
	b.call("window_get_width", result, window)
}

func (v *Vim) SetWindowWidth(window Window, width int) error {
	return v.call("window_set_width", nil, window, width)
}

func (b *Batch) SetWindowWidth(window Window, width int) {
	b.call("window_set_width", nil, window, width)
}

func (v *Vim) WindowVar(window Window, name string, result interface{}) error {
	return v.call("window_get_var", result, window, name)
}

func (b *Batch) WindowVar(window Window, name string, result interface{}) {
	b.call("window_get_var", result, window, name)
}

func (v *Vim) SetWindowVar(window Window, name string, value interface{}, result interface{}) error {
	return v.call("window_set_var", result, window, name, value)
}

func (b *Batch) SetWindowVar(window Window, name string, value interface{}, result interface{}) {
	b.call("window_set_var", result, window, name, value)
}

func (v *Vim) WindowOption(window Window, name string, result interface{}) error {
	return v.call("window_get_option", result, window, name)
}

func (b *Batch) WindowOption(window Window, name string, result interface{}) {
	b.call("window_get_option", result, window, name)
}

func (v *Vim) SetWindowOption(window Window, name string, value interface{}) error {
	return v.call("window_set_option", nil, window, name, value)
}

func (b *Batch) SetWindowOption(window Window, name string, value interface{}) {
	b.call("window_set_option", nil, window, name, value)
}

func (v *Vim) WindowPosition(window Window) ([2]int, error) {
	var result [2]int
	err := v.call("window_get_position", &result, window)
	return result, err
}

func (b *Batch) WindowPosition(window Window, result *[2]int) {
	b.call("window_get_position", result, window)
}

func (v *Vim) WindowTabpage(window Window) (Tabpage, error) {
	var result Tabpage
	err := v.call("window_get_tabpage", &result, window)
	return result, err
}

func (b *Batch) WindowTabpage(window Window, result *Tabpage) {
	b.call("window_get_tabpage", result, window)
}

func (v *Vim) IsWindowValid(window Window) (bool, error) {
	var result bool
	err := v.call("window_is_valid", &result, window)
	return result, err
}

func (b *Batch) IsWindowValid(window Window, result *bool) {
	b.call("window_is_valid", result, window)
}