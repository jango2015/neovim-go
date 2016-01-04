// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// This program generates Neovim API methods in api.go.
//
// The program generates the code from data declared in this file instead of
// using the output from nvim --api-info. This approach allows names and types
// to be modified to create a more idiomatic and convenient API for Go
// programmers.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"strings"
	"text/template"
)

var extensions = []*struct {
	Type string
	Code int
	Doc  string
}{
	{"Buffer", 0, `// Buffer represents a remote Neovim buffer.`},
	{"Window", 1, `// Window represents a remote Neovim window.`},
	{"Tabpage", 2, `// Tabpage represents a remote Neovim tabpage.`},
}

type param struct{ Name, Type string }

var methods = []*struct {
	Name   string
	Sm     string
	Return string
	Doc    string
	Params []param
}{
	{
		Name:   "BufferLineCount",
		Sm:     "buffer_line_count",
		Return: "int",
		Params: []param{{"buffer", "Buffer"}},
		Doc:    `// BufferLineCount returns the number of lines in the buffer.`,
	},
	{
		Name:   "BufferLine",
		Sm:     "buffer_get_line",
		Return: "[]byte",
		Params: []param{{"buffer", "Buffer"}, {"index", "int"}},
		Doc:    `// BufferLine returns the line at the given index.`,
	},
	{
		Name:   "SetBufferLine",
		Sm:     "buffer_set_line",
		Params: []param{{"buffer", "Buffer"}, {"index", "int"}, {"line", "[]byte"}},
		Doc:    `// SetBufferLine sets the line at the given index.`,
	},
	{
		Name:   "DeleteBufferLine",
		Sm:     "buffer_del_line",
		Params: []param{{"buffer", "Buffer"}, {"index", "int"}},
		Doc:    `// DeleteBufferLine deletes the line at the given index.`,
	},
	{
		Name:   "BufferLineSlice",
		Sm:     "buffer_get_line_slice",
		Return: "[][]byte",
		Params: []param{{"buffer", "Buffer"}, {"start", "int"}, {"end", "int"}, {"includeStart", "bool"}, {"includeEnd", "bool"}},
		Doc:    `// BufferLineSlice retrieves a line range from a buffer.`,
	},
	{
		Name:   "SetBufferLineSlice",
		Sm:     "buffer_set_line_slice",
		Params: []param{{"buffer", "Buffer"}, {"start", "int"}, {"end", "int"}, {"includeStart", "bool"}, {"includeEnd", "bool"}, {"replacement", "[][]byte"}},
		Doc:    `// SetBufferLineSlice replaces a line range on a buffer.`,
	},
	{
		Name:   "BufferVar",
		Sm:     "buffer_get_var",
		Return: "interface{}",
		Params: []param{{"buffer", "Buffer"}, {"name", "string"}},
		Doc:    `// BufferVar gets a buffer-scoped (b:) variable.`,
	},
	{
		Name:   "SetBufferVar",
		Sm:     "buffer_set_var",
		Return: "interface{}",
		Params: []param{{"buffer", "Buffer"}, {"name", "string"}, {"value", "interface{}"}},
		Doc: `
// SetBufferVar sets a buffer-scoped (b:) variable. The value nil deletes the
// variable. Result is the previous value of the variable.
`,
	},
	{
		Name:   "BufferOption",
		Sm:     "buffer_get_option",
		Return: "interface{}",
		Params: []param{{"buffer", "Buffer"}, {"name", "string"}},
		Doc:    `// BufferOption gets a buffer option value.`,
	},
	{
		Name:   "SetBufferOption",
		Sm:     "buffer_set_option",
		Params: []param{{"buffer", "Buffer"}, {"name", "string"}, {"value", "interface{}"}},
		Doc: `
// SetBufferOption sets a buffer option value. The value nil deletes the option
// in the case where there's a global fallback.
`,
	},
	{
		Name:   "BufferNumber",
		Sm:     "buffer_get_number",
		Return: "int",
		Params: []param{{"buffer", "Buffer"}},
		Doc:    `// BufferNumber gets a buffer's number.`,
	},
	{
		Name:   "BufferName",
		Sm:     "buffer_get_name",
		Return: "string",
		Params: []param{{"buffer", "Buffer"}},
		Doc:    `// BufferName gets the full file name of a buffer.`,
	},
	{
		Name:   "SetBufferName",
		Sm:     "buffer_set_name",
		Params: []param{{"buffer", "Buffer"}, {"name", "string"}},
		Doc: `
// SetBufferName sets the full file name of a buffer.
// BufFilePre/BufFilePost are triggered.
`,
	},
	{
		Name:   "IsBufferValid",
		Sm:     "buffer_is_valid",
		Return: "bool",
		Params: []param{{"buffer", "Buffer"}},
		Doc:    `// IsBufferValid returns true if the buffer is valid.`,
	},
	{
		Name:   "InsertBuffer",
		Sm:     "buffer_insert",
		Params: []param{{"buffer", "Buffer"}, {"lnum", "int"}, {"lines", "[][]byte"}},
		Doc:    `// InsertBUffer inserts a range of lines to a buffer at the specified index.`,
	},
	{
		Name:   "BufferMark",
		Sm:     "buffer_get_mark",
		Return: "[2]int",
		Params: []param{{"buffer", "Buffer"}, {"name", "string"}},
		Doc:    `// BufferMark returns the (row,col) of the named mark.`,
	},
	{
		Name:   "TabpageWindows",
		Sm:     "tabpage_get_windows",
		Return: "[]Window",
		Params: []param{{"tabpage", "Tabpage"}},
	},
	{
		Name:   "TabpageVar",
		Sm:     "tabpage_get_var",
		Return: "interface{}",
		Params: []param{{"tabpage", "Tabpage"}, {"name", "string"}},
	},
	{
		Name:   "SetTabpageVar",
		Sm:     "tabpage_set_var",
		Return: "interface{}",
		Params: []param{{"tabpage", "Tabpage"}, {"name", "string"}, {"value", "interface{}"}},
	},
	{
		Name:   "TabpageWindow",
		Sm:     "tabpage_get_window",
		Return: "Window",
		Params: []param{{"tabpage", "Tabpage"}},
	},
	{
		Name:   "IsTabpageValid",
		Sm:     "tabpage_is_valid",
		Return: "bool",
		Params: []param{{"tabpage", "Tabpage"}},
	},
	{
		Name:   "Command",
		Sm:     "vim_command",
		Params: []param{{"str", "string"}},
		Doc: `
// Command executes a single ex command.
`,
	},
	{
		Name:   "Feedkeys",
		Sm:     "vim_feedkeys",
		Params: []param{{"keys", "string"}, {"mode", "string"}, {"escapeCsi", "bool"}},
		Doc: `
// FeedKeys Pushes keys to the Neovim user input buffer. Options can be a string
// with the following character flags:
//
//  m:  Remap keys. This is default.
//  n:  Do not remap keys.
//  t:  Handle keys as if typed; otherwise they are handled as if coming from a
//     mapping. This matters for undo, opening folds, etc. 
`,
	},
	{
		Name:   "Input",
		Sm:     "vim_input",
		Return: "int",
		Params: []param{{"keys", "string"}},
		Doc: `
// Input pushes bytes to the Neovim low level input buffer.
// 
// Unlike FeedKeys, this uses the lowest level input buffer and the call is not
// deferred. It returns the number of bytes actually written(which can be less
// than what was requested if the buffer is full).
`,
	},
	{
		Name:   "ReplaceTermcodes",
		Sm:     "vim_replace_termcodes",
		Return: "string",
		Params: []param{{"str", "string"}, {"fromPart", "bool"}, {"doLt", "bool"}, {"special", "bool"}},
		Doc: `
// ReplaceTermcodes replaces any terminal code strings by byte sequences. The
// returned sequences are Nvim's internal representation of keys, for example:
//
//  <esc> -> '\x1b'
//  <cr>  -> '\r'
//  <c-l> -> '\x0c'
//  <up>  -> '\x80ku'
//
// The returned sequences can be used as input to feedkeys.
`,
	},
	{
		Name:   "CommandOutput",
		Sm:     "vim_command_output",
		Return: "string",
		Params: []param{{"str", "string"}},
		Doc: `
// CommandOutput executes a single ex command and returns the output.
`,
	},
	{
		Name:   "Eval",
		Sm:     "vim_eval",
		Return: "interface{}",
		Params: []param{{"str", "string"}},
		Doc: `
// Eval evaluates the expression str using the Vim internal expression
// evaluator. 
//
//  :help expression
`,
	},
	{
		Name:   "Strwidth",
		Sm:     "vim_strwidth",
		Return: "int",
		Params: []param{{"str", "string"}},
		Doc: `
// Strwidth returns the number of display cells the string occupies. Tab is
// counted as one cell.
`,
	},
	{
		Name:   "RuntimePaths",
		Sm:     "vim_list_runtime_paths",
		Return: "[]string",
		Doc: `
// RuntimePaths returns a list of paths contained in the runtimepath option.
`,
	},
	{
		Name:   "ChangeDirectory",
		Sm:     "vim_change_directory",
		Params: []param{{"dir", "string"}},
		Doc:    `// ChangeDirectory changes Vim working directory.`,
	},
	{
		Name:   "CurrentLine",
		Sm:     "vim_get_current_line",
		Return: "[]byte",
		Doc:    `// CurrentLine gets the current line in the current buffer.`,
	},
	{
		Name:   "SetCurrentLine",
		Sm:     "vim_set_current_line",
		Params: []param{{"line", "[]byte"}},
		Doc:    `// SetCurrentLine sets the current line in the current buffer.`,
	},
	{
		Name: "DeleteCurrentLine",
		Sm:   "vim_del_current_line",
		Doc:  `// DeleteCurrentLine deletes the current line in the current buffer.`,
	},
	{
		Name:   "Var",
		Sm:     "vim_get_var",
		Return: "interface{}",
		Params: []param{{"name", "string"}},
		Doc:    `// Var gets a global variable.`,
	},
	{
		Name:   "SetVar",
		Sm:     "vim_set_var",
		Return: "interface{}",
		Params: []param{{"name", "string"}, {"value", "interface{}"}},
		Doc: `
// SetVar sets a global variable. The value nil deletes the variable. Result is
// the previous value of the variable.
`,
	},
	{
		Name:   "Vvar",
		Sm:     "vim_get_vvar",
		Return: "interface{}",
		Params: []param{{"name", "string"}},
		Doc:    `// Vvar gets a vim variable.`,
	},
	{
		Name:   "Option",
		Sm:     "vim_get_option",
		Return: "interface{}",
		Params: []param{{"name", "string"}},
		Doc:    `// Option gets an option.`,
	},
	{
		Name:   "SetOption",
		Sm:     "vim_set_option",
		Params: []param{{"name", "string"}, {"value", "interface{}"}},
		Doc:    `// SetOption sets an option.`,
	},
	{
		Name:   "WriteOut",
		Sm:     "vim_out_write",
		Params: []param{{"str", "string"}},
		Doc:    `// WriteOut prints str as a normal message.`,
	},
	{
		Name:   "WriteErr",
		Sm:     "vim_err_write",
		Params: []param{{"str", "string"}},
		Doc:    `// WriteErr prints str as an error message.`,
	},
	{
		Name:   "ReportError",
		Sm:     "vim_report_error",
		Params: []param{{"str", "string"}},
		Doc:    `// ReportError writes prints str and a newline as an error message.`,
	},
	{
		Name:   "Buffers",
		Sm:     "vim_get_buffers",
		Return: "[]Buffer",
	},
	{
		Name:   "CurrentBuffer",
		Sm:     "vim_get_current_buffer",
		Return: "Buffer",
	},
	{
		Name:   "SetCurrentBuffer",
		Sm:     "vim_set_current_buffer",
		Params: []param{{"buffer", "Buffer"}},
	},
	{
		Name:   "Windows",
		Sm:     "vim_get_windows",
		Return: "[]Window",
	},
	{
		Name:   "CurrentWindow",
		Sm:     "vim_get_current_window",
		Return: "Window",
	},
	{
		Name:   "SetCurrentWindow",
		Sm:     "vim_set_current_window",
		Params: []param{{"window", "Window"}},
	},
	{
		Name:   "Tabpages",
		Sm:     "vim_get_tabpages",
		Return: "[]Tabpage",
	},
	{
		Name:   "CurrentTabpage",
		Sm:     "vim_get_current_tabpage",
		Return: "Tabpage",
	},
	{
		Name:   "SetCurrentTabpage",
		Sm:     "vim_set_current_tabpage",
		Params: []param{{"tabpage", "Tabpage"}},
	},
	{
		Name:   "Subscribe",
		Sm:     "vim_subscribe",
		Params: []param{{"event", "string"}},
		Doc:    `// Subscribe subscribes to a Neovim event.`,
	},
	{
		Name:   "Unsubscribe",
		Sm:     "vim_unsubscribe",
		Params: []param{{"event", "string"}},
		Doc:    `// Unsubscribe unsubscribes to a Neovim event.`,
	},
	{
		Name:   "NameToColor",
		Sm:     "vim_name_to_color",
		Return: "int",
		Params: []param{{"name", "string"}},
	},
	{
		Name:   "ColorMap",
		Sm:     "vim_get_color_map",
		Return: "map[string]interface{}",
	},
	{
		Name:   "APIInfo",
		Sm:     "vim_get_api_info",
		Return: "[]interface{}",
	},
	{
		Name:   "WindowBuffer",
		Sm:     "window_get_buffer",
		Return: "Buffer",
		Params: []param{{"window", "Window"}},
	},
	{
		Name:   "WindowCursor",
		Sm:     "window_get_cursor",
		Return: "[2]int",
		Params: []param{{"window", "Window"}},
	},
	{
		Name:   "SetWindowCursor",
		Sm:     "window_set_cursor",
		Params: []param{{"window", "Window"}, {"pos", "[2]int"}},
	},
	{
		Name:   "WindowHeight",
		Sm:     "window_get_height",
		Return: "int",
		Params: []param{{"window", "Window"}},
	},
	{
		Name:   "SetWindowHeight",
		Sm:     "window_set_height",
		Params: []param{{"window", "Window"}, {"height", "int"}},
	},
	{
		Name:   "WindowWidth",
		Sm:     "window_get_width",
		Return: "int",
		Params: []param{{"window", "Window"}},
	},
	{
		Name:   "SetWindowWidth",
		Sm:     "window_set_width",
		Params: []param{{"window", "Window"}, {"width", "int"}},
	},
	{
		Name:   "WindowVar",
		Sm:     "window_get_var",
		Return: "interface{}",
		Params: []param{{"window", "Window"}, {"name", "string"}},
	},
	{
		Name:   "SetWindowVar",
		Sm:     "window_set_var",
		Return: "interface{}",
		Params: []param{{"window", "Window"}, {"name", "string"}, {"value", "interface{}"}},
	},
	{
		Name:   "WindowOption",
		Sm:     "window_get_option",
		Return: "interface{}",
		Params: []param{{"window", "Window"}, {"name", "string"}},
	},
	{
		Name:   "SetWindowOption",
		Sm:     "window_set_option",
		Params: []param{{"window", "Window"}, {"name", "string"}, {"value", "interface{}"}},
	},
	{
		Name:   "WindowPosition",
		Sm:     "window_get_position",
		Return: "[2]int",
		Params: []param{{"window", "Window"}},
	},
	{
		Name:   "WindowTabpage",
		Sm:     "window_get_tabpage",
		Return: "Tabpage",
		Params: []param{{"window", "Window"}},
	},
	{
		Name:   "IsWindowValid",
		Sm:     "window_is_valid",
		Return: "bool",
		Params: []param{{"window", "Window"}},
	},
}

var templ = template.Must(template.New("").Parse(`// Code generated by 'go generate'

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
{{range .Extensions}}
		{{.Code}}: func(p []byte) (interface{}, error) {
			x, err := decodeExt(p)
			return {{.Type}}(x), err
		},
{{end}}
	})
}

{{range .Extensions}}
{{.Doc}}
type {{.Type}} int

func (x *{{.Type}}) UnmarshalMsgPack(dec *msgpack.Decoder) error {
	if dec.Type() != msgpack.Extension || dec.Extension() != {{.Code}} {
		err := &msgpack.DecodeConvertError{
			SrcType:  dec.Type(),
			DestType: reflect.TypeOf(x),
		}
		dec.Skip()
		return err
	}
	n, err := decodeExt(dec.BytesNoCopy())
	*x = {{.Type}}(n)
	return err
}

func (x {{.Type}}) MarshalMsgPack(enc *msgpack.Encoder) error {
	return enc.PackExtension({{.Code}}, encodeExt(int(x)))
}

func (x {{.Type}}) String() string {
	return fmt.Sprintf("{{.Type}}:%d", int(x))
}
{{end}}

{{range .Methods}}
{{if eq "interface{}" .Return}}
{{.Doc}}
func (v *Vim) {{.Name}}({{range .Params}}{{.Name}} {{.Type}},{{end}} result interface{}) error {
    return v.call("{{.Sm}}", result, {{range .Params}}{{.Name}},{{end}})
}

{{.Doc}}
func (p *Pipeline) {{.Name}}({{range .Params}}{{.Name}} {{.Type}},{{end}} result interface{}) {
    p.call("{{.Sm}}", result, {{range .Params}}{{.Name}},{{end}})
}
{{else if .Return}}
{{.Doc}}
func (v *Vim) {{.Name}}({{range .Params}}{{.Name}} {{.Type}},{{end}}) ({{.Return}}, error) {
    var result {{.Return}}
    err := v.call("{{.Sm}}", {{if .Return}}&result{{else}}nil{{end}}, {{range .Params}}{{.Name}},{{end}})
    return result, err
}
{{.Doc}}
func (p *Pipeline) {{.Name}}({{range .Params}}{{.Name}} {{.Type}},{{end}} result *{{.Return}}) {
    p.call("{{.Sm}}", result, {{range .Params}}{{.Name}},{{end}})
}
{{else}}
{{.Doc}}
func (v *Vim) {{.Name}}({{range .Params}}{{.Name}} {{.Type}},{{end}}) error {
    return v.call("{{.Sm}}", nil, {{range .Params}}{{.Name}},{{end}})
}
{{.Doc}}
func (p *Pipeline) {{.Name}}({{range .Params}}{{.Name}} {{.Type}},{{end}}) {
    p.call("{{.Sm}}", nil, {{range .Params}}{{.Name}},{{end}})
}
{{end}}
{{end}}
`))

func main() {
	log.SetFlags(0)
	outFile := flag.String("out", "", "Output file")
	flag.Parse()

	for _, m := range methods {
		m.Doc = strings.TrimSpace(m.Doc)
	}

	var buf bytes.Buffer
	if err := templ.Execute(&buf, map[string]interface{}{
		"Methods":    methods,
		"Extensions": extensions,
	}); err != nil {
		log.Fatalf("error executing template: %v", err)
	}

	out, err := format.Source(buf.Bytes())
	if err != nil {
		for i, p := range bytes.Split(buf.Bytes(), []byte("\n")) {
			fmt.Fprintf(os.Stderr, "%d: %s\n", i+1, p)
		}
		log.Fatalf("error formating source: %v", err)
	}

	f := os.Stdout
	if *outFile != "" {
		f, err = os.Create(*outFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}

	f.Write(out)
}
