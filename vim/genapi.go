// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// This program generates api.go.

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
	},
	{
		Name:   "BufferLine",
		Sm:     "buffer_get_line",
		Return: "[]byte",
		Params: []param{{"buffer", "Buffer"}, {"index", "int"}},
	},
	{
		Name:   "SetBufferLine",
		Sm:     "buffer_set_line",
		Params: []param{{"buffer", "Buffer"}, {"index", "int"}, {"line", "[]byte"}},
	},
	{
		Name:   "DeleteBufferLine",
		Sm:     "buffer_del_line",
		Params: []param{{"buffer", "Buffer"}, {"index", "int"}},
	},
	{
		Name:   "BufferLineSlice",
		Sm:     "buffer_get_line_slice",
		Return: "[][]byte",
		Params: []param{{"buffer", "Buffer"}, {"start", "int"}, {"end", "int"}, {"includeStart", "bool"}, {"includeEnd", "bool"}},
	},
	{
		Name:   "SetBufferLineSlice",
		Sm:     "buffer_set_line_slice",
		Params: []param{{"buffer", "Buffer"}, {"start", "int"}, {"end", "int"}, {"includeStart", "bool"}, {"includeEnd", "bool"}, {"replacement", "[][]byte"}},
	},
	{
		Name:   "BufferVar",
		Sm:     "buffer_get_var",
		Return: "interface{}",
		Params: []param{{"buffer", "Buffer"}, {"name", "string"}},
	},
	{
		Name:   "SetBufferVar",
		Sm:     "buffer_set_var",
		Return: "interface{}",
		Params: []param{{"buffer", "Buffer"}, {"name", "string"}, {"value", "interface{}"}},
	},
	{
		Name:   "BufferOption",
		Sm:     "buffer_get_option",
		Return: "interface{}",
		Params: []param{{"buffer", "Buffer"}, {"name", "string"}},
	},
	{
		Name:   "SetBufferOption",
		Sm:     "buffer_set_option",
		Params: []param{{"buffer", "Buffer"}, {"name", "string"}, {"value", "interface{}"}},
	},
	{
		Name:   "BufferNumber",
		Sm:     "buffer_get_number",
		Return: "int",
		Params: []param{{"buffer", "Buffer"}},
	},
	{
		Name:   "BufferName",
		Sm:     "buffer_get_name",
		Return: "string",
		Params: []param{{"buffer", "Buffer"}},
	},
	{
		Name:   "SetBufferName",
		Sm:     "buffer_set_name",
		Params: []param{{"buffer", "Buffer"}, {"name", "string"}},
	},
	{
		Name:   "IsBufferValid",
		Sm:     "buffer_is_valid",
		Return: "bool",
		Params: []param{{"buffer", "Buffer"}},
	},
	{
		Name:   "InsertBuffer",
		Sm:     "buffer_insert",
		Params: []param{{"buffer", "Buffer"}, {"lnum", "int"}, {"lines", "[][]byte"}},
	},
	{
		Name:   "BufferMark",
		Sm:     "buffer_get_mark",
		Return: "[2]int",
		Params: []param{{"buffer", "Buffer"}, {"name", "string"}},
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
// Eval evaluates a vimscript expression.
`,
	},
	{
		Name:   "Strwidth",
		Sm:     "vim_strwidth",
		Return: "int",
		Params: []param{{"str", "string"}},
		Doc: `
// Strwidth returns the number of display cells string occupies. Tab is counted
// as one cell.
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
	},
	{
		Name:   "CurrentLine",
		Sm:     "vim_get_current_line",
		Return: "[]byte",
	},
	{
		Name:   "SetCurrentLine",
		Sm:     "vim_set_current_line",
		Params: []param{{"line", "[]byte"}},
	},
	{
		Name: "DeleteCurrentLine",
		Sm:   "vim_del_current_line",
	},
	{
		Name:   "Var",
		Sm:     "vim_get_var",
		Return: "interface{}",
		Params: []param{{"name", "string"}},
	},
	{
		Name:   "SetVar",
		Sm:     "vim_set_var",
		Return: "interface{}",
		Params: []param{{"name", "string"}, {"value", "interface{}"}},
	},
	{
		Name:   "Vvar",
		Sm:     "vim_get_vvar",
		Return: "interface{}",
		Params: []param{{"name", "string"}},
	},
	{
		Name:   "Option",
		Sm:     "vim_get_option",
		Return: "interface{}",
		Params: []param{{"name", "string"}},
	},
	{
		Name:   "SetOption",
		Sm:     "vim_set_option",
		Params: []param{{"name", "string"}, {"value", "interface{}"}},
	},
	{
		Name:   "WriteOut",
		Sm:     "vim_out_write",
		Params: []param{{"str", "string"}},
		Doc: `
// WriteOut writes a message to the output buffer.
`,
	},
	{
		Name:   "WriteErr",
		Sm:     "vim_err_write",
		Params: []param{{"str", "string"}},
		Doc: `
// WriteErr writes a message to the error buffer.
`,
	},
	{
		Name:   "ReportError",
		Sm:     "vim_report_error",
		Params: []param{{"str", "string"}},
		Doc: `
// ReportError writes a message and a newline to the error buffer.
`,
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

{{range .}}
{{if eq "interface{}" .Return}}
{{.Doc}}
func (v *Vim) {{.Name}}({{range .Params}}{{.Name}} {{.Type}},{{end}} result interface{}) error {
    return v.call("{{.Sm}}", result, {{range .Params}}{{.Name}},{{end}})
}

{{.Doc}}
func (b *Batch) {{.Name}}({{range .Params}}{{.Name}} {{.Type}},{{end}} result interface{}) {
    b.call("{{.Sm}}", result, {{range .Params}}{{.Name}},{{end}})
}
{{else if .Return}}
{{.Doc}}
func (v *Vim) {{.Name}}({{range .Params}}{{.Name}} {{.Type}},{{end}}) ({{.Return}}, error) {
    var result {{.Return}}
    err := v.call("{{.Sm}}", {{if .Return}}&result{{else}}nil{{end}}, {{range .Params}}{{.Name}},{{end}})
    return result, err
}
{{.Doc}}
func (b *Batch) {{.Name}}({{range .Params}}{{.Name}} {{.Type}},{{end}} result *{{.Return}}) {
    b.call("{{.Sm}}", result, {{range .Params}}{{.Name}},{{end}})
}
{{else}}
{{.Doc}}
func (v *Vim) {{.Name}}({{range .Params}}{{.Name}} {{.Type}},{{end}}) error {
    return v.call("{{.Sm}}", nil, {{range .Params}}{{.Name}},{{end}})
}
{{.Doc}}
func (b *Batch) {{.Name}}({{range .Params}}{{.Name}} {{.Type}},{{end}}) {
    b.call("{{.Sm}}", nil, {{range .Params}}{{.Name}},{{end}})
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
	if err := templ.Execute(&buf, methods); err != nil {
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