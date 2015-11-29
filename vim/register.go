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

var noOptions = make(map[string]string)

func (v *Vim) RegisterFunction(name string, handler interface{}) error {
	if err := v.ep.RegisterHandler(v.PluginPath+":function:"+name, handler); err != nil {
		return err
	}
	v.pluginSpecs = append(v.pluginSpecs, &pluginSpec{
		Type: "function",
		Name: name,
		Sync: isSync(handler),
		Opts: noOptions,
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

	// Eval is evaluated in Neovim and the result is provided in CommandArgs.
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

func mapFromCommandOptions(o *CommandOptions) map[string]string {
	m := make(map[string]string)
	if o == nil {
		return m
	}

	if o.NArgs != "" {
		m["nargs"] = o.NArgs
	}

	if o.Range != "" {
		if o.Range == "." {
			o.Range = ""
		}
		m["range"] = o.Range
	} else if o.Count != "" {
		m["count"] = o.Count
	}

	if o.Bang {
		m["bang"] = ""
	}

	if o.Register {
		m["register"] = ""
	}

	if o.Eval != "" {
		m["eval"] = o.Eval
	}

	if o.Addr != "" {
		m["addr"] = o.Addr
	}

	if o.Bar {
		m["bar"] = ""
	}

	if o.Complete != "" {
		m["complete"] = o.Complete
	}

	return m
}

func (v *Vim) RegisterCommand(name string, options *CommandOptions, handler interface{}) error {
	sm := fmt.Sprintf("%s:command:%s", v.PluginPath, name)
	if err := v.ep.RegisterHandler(sm, handler); err != nil {
		return err
	}

	v.pluginSpecs = append(v.pluginSpecs, &pluginSpec{
		Type: "command",
		Name: name,
		Sync: isSync(handler),
		Opts: mapFromCommandOptions(options),
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

func mapFromAutocmdOptions(o *AutocmdOptions) map[string]string {
	m := make(map[string]string)
	if o == nil {
		return m
	}

	if o.Pattern != "" {
		m["pattern"] = o.Pattern
	}

	if o.Nested {
		m["nested"] = ""
	}

	if o.Eval != "" {
		m["eval"] = o.Eval
	}

	return m
}

func (v *Vim) RegisterAutocmd(event string, options *AutocmdOptions, handler interface{}) error {
	pattern := ""
	if options != nil {
		pattern = options.Pattern
	}

	sm := fmt.Sprintf("%s:autocmd:%s:%s", v.PluginPath, event, pattern)
	if err := v.ep.RegisterHandler(sm, handler); err != nil {
		return err
	}

	v.pluginSpecs = append(v.pluginSpecs, &pluginSpec{
		Type: "autocmd",
		Name: event,
		Sync: isSync(handler),
		Opts: mapFromAutocmdOptions(options),
	})

	return nil
}
