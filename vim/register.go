// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vim

import (
	"fmt"
	"reflect"
)

type pluginSpec struct {
	Type string            `msgpack:"type"`
	Name string            `msgpack:"name"`
	Sync bool              `msgpack:"sync"`
	Opts map[string]string `msgpack:"opts"`
}

func isSync(f interface{}) bool {
	t := reflect.TypeOf(f)
	return t.Kind() == reflect.Func && t.NumOut() > 0
}

func (v *Vim) handleSpecs(p string) ([]*pluginSpec, error) {
	return v.pluginSpecs, nil
}

func (v *Vim) RegisterHandler(serviceMethod string, handler interface{}) error {
	return v.ep.RegisterHandler(serviceMethod, handler)
}

// FunctionOptions specifies registered function options.
type FunctionOptions struct {
	// Eval is evaluated in Neovim and the result is passed the the handler
	// function.
	Eval string
}

func (v *Vim) RegisterFunction(name string, options *FunctionOptions, handler interface{}) error {
	m := make(map[string]string)
	if options != nil {
		if options.Eval != "" {
			m["eval"] = options.Eval
		}
	}

	if err := v.ep.RegisterHandler(v.PluginPath+":function:"+name, handler); err != nil {
		return err
	}
	v.pluginSpecs = append(v.pluginSpecs, &pluginSpec{
		Type: "function",
		Name: name,
		Sync: isSync(handler),
		Opts: m,
	})
	return nil
}

// CommandOptions specifies command options.
type CommandOptions struct {

	// NArgs specifies the number command arguments.
	//
	//  0   No arguments are allowed
	//  1   Exactly one argument is required, it includes spaces
	//  *   Any number of arguments are allowed (0, 1, or many),
	//      separated by white space
	//  ?   0 or 1 arguments are allowed
	//  +   Arguments must be supplied, but any number are allowed
	NArgs string

	// Range specifies that the command accepts a range.
	//
	//  .   Range allowed, default is current line. The value
	//      "." is converted to "" for Neovim.
	//  %   Range allowed, default is whole file (1,$)
	//  N   A count (default N) which is specified in the line
	//      number position (like |:split|); allows for zero line
	//	    number.
	//
	//  :help :command-range
	Range string

	// Count specfies that thecommand accepts a count.
	//
	//  N   A count (default N) which is specified either in the line
	//	    number position, or as an initial argument (like |:Next|).
	//      Specifying -count (without a default) acts like -count=0
	//
	//  :help :command-count
	Count string

	// Addr sepcifies the domain for the range option
	//
	//  lines           Range of lines (this is the default)
	//  arguments       Range for arguments
	//  buffers         Range for buffers (also not loaded buffers)
	//  loaded_buffers  Range for loaded buffers
	//  windows         Range for windows
	//  tabs            Range for tab pages
	//
	//  :help command-addr
	Addr string

	// Bang specifies that the command can take a ! modifier (like :q or :w).
	Bang bool

	// Register specifes that the first argument to the command can be an
	// optional register name (like :del, :put, :yank).
	Register bool

	// Eval is evaluated in Neovim and the result is passed as an argument.
	Eval string

	// Bar specifies that the command can be followed by a "|" and another
	// command.  A "|" inside the command argument is not allowed then. Also
	// checks for a " to start a comment.
	Bar bool

	// Complete specifies command completion.
	//
	//  :help :command-complete
	Complete string
}

// RegisterCommand registers a handler for a command with the given options.
//
// The arguments to the handler function are:
//
//  v *vim.Vim
//  args []string       when NArgs != ""
//  range [2]int        when Range == "." or Range == "%"
//  range int           when Range == N or Count != ""
//  bang bool           when Bang == true
//  register string     when Register == true
//  eval interface{}    when Eval != ""
func (v *Vim) RegisterCommand(name string, options *CommandOptions, handler interface{}) error {
	m := make(map[string]string)
	if options != nil {

		if options.NArgs != "" {
			m["nargs"] = options.NArgs
		}

		if options.Range != "" {
			if options.Range == "." {
				options.Range = ""
			}
			m["range"] = options.Range
		} else if options.Count != "" {
			m["count"] = options.Count
		}

		if options.Bang {
			m["bang"] = ""
		}

		if options.Register {
			m["register"] = ""
		}

		if options.Eval != "" {
			m["eval"] = options.Eval
		}

		if options.Addr != "" {
			m["addr"] = options.Addr
		}

		if options.Bar {
			m["bar"] = ""
		}

		if options.Complete != "" {
			m["complete"] = options.Complete
		}
	}

	sm := fmt.Sprintf("%s:command:%s", v.PluginPath, name)
	if err := v.ep.RegisterHandler(sm, handler); err != nil {
		return err
	}

	v.pluginSpecs = append(v.pluginSpecs, &pluginSpec{
		Type: "command",
		Name: name,
		Sync: isSync(handler),
		Opts: m,
	})
	return nil
}

// AutocmdOptions specifies autocmd options.
type AutocmdOptions struct {
	// Pattern specifies an autocmd pattern.
	//
	//  :help autocmd-patterns
	Pattern string

	// Nested allows nested autocmds.
	//
	//  :help autocmd-nested
	Nested bool

	// Eval is evaluated in Neovim and the result is passed the the handler
	// function.
	Eval string
}

func (v *Vim) RegisterAutocmd(event string, options *AutocmdOptions, handler interface{}) error {
	pattern := ""
	m := make(map[string]string)
	if options != nil {

		if options.Pattern != "" {
			m["pattern"] = options.Pattern
			pattern = options.Pattern
		}

		if options.Nested {
			m["nested"] = ""
		}

		if options.Eval != "" {
			m["eval"] = options.Eval
		}

	}

	sm := fmt.Sprintf("%s:autocmd:%s:%s", v.PluginPath, event, pattern)
	if err := v.ep.RegisterHandler(sm, handler); err != nil {
		return err
	}

	v.pluginSpecs = append(v.pluginSpecs, &pluginSpec{
		Type: "autocmd",
		Name: event,
		Sync: isSync(handler),
		Opts: m,
	})

	return nil
}
