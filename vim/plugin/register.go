// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plugin

import (
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/garyburd/neovim-go/vim"
)

type pluginSpec struct {
	Type string            `msgpack:"type"`
	Name string            `msgpack:"name"`
	Sync bool              `msgpack:"sync"`
	Opts map[string]string `msgpack:"opts"`

	smSuffix string
	fn       interface{}
}

type handler struct {
	sm string
	fn interface{}
}

var (
	pluginSpecs = []*pluginSpec{}
	handlers    []*handler
)

func isSync(f interface{}) bool {
	t := reflect.TypeOf(f)
	return t.Kind() == reflect.Func && t.NumOut() > 0
}

func RegisterHandlers(v *vim.Vim, paths ...string) error {
	var mu sync.Mutex
	done := false
	err := v.RegisterHandler("specs", func(v *vim.Vim, path string) ([]*pluginSpec, error) {
		mu.Lock()
		defer mu.Unlock()
		if done {
			return []*pluginSpec{}, nil
		}
		done = true
		return pluginSpecs, nil
	})
	if err != nil {
		return err
	}
	for _, path := range paths {
		for _, s := range pluginSpecs {
			log.Println(path + s.smSuffix)
			if err := v.RegisterHandler(path+s.smSuffix, s.fn); err != nil {
				return err
			}
		}
	}
	for _, h := range handlers {
		if err := v.RegisterHandler(h.sm, h.fn); err != nil {
			return err
		}
	}
	return nil
}

// Handle registers fn as a MsgPack RPC handler for the specified method name.
func Handle(serviceMethod string, fn interface{}) {
	handlers = append(handlers, &handler{fn: fn, sm: serviceMethod})
}

// FunctionOptions specifies function options.
type FunctionOptions struct {
	// Eval is evaluated in Neovim and the result is passed the the handler
	// function.
	Eval string
}

// HandleFunction registers fn as a handler for a Neovim function with the
// specified name.
func HandleFunction(name string, options *FunctionOptions, fn interface{}) {
	m := make(map[string]string)
	if options != nil {
		if options.Eval != "" {
			m["eval"] = options.Eval
		}
	}
	pluginSpecs = append(pluginSpecs, &pluginSpec{
		Type: "function",
		Name: name,
		Sync: isSync(fn),
		Opts: m,

		fn:       fn,
		smSuffix: ":function:" + name,
	})
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

// HandleCommand registers fn as a handler for a Neovim command with the
// specified name.
//
// The arguments to fn function are:
//
//  v *vim.Vim
//  args []string       when options.NArgs != ""
//  range [2]int        when options.Range == "." or Range == "%"
//  range int           when options.Range == N or Count != ""
//  bang bool           when options.Bang == true
//  register string     when options.Register == true
//  eval interface{}    when options.Eval != ""
func HandleCommand(name string, options *CommandOptions, fn interface{}) error {
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

	pluginSpecs = append(pluginSpecs, &pluginSpec{
		Type: "command",
		Name: name,
		Sync: isSync(fn),
		Opts: m,

		smSuffix: ":command:" + name,
		fn:       fn,
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

// HandleAutocmd registers fn as a handler for the specified autocmnd event.
func HandleAutocmd(event string, options *AutocmdOptions, fn interface{}) {
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
	pluginSpecs = append(pluginSpecs, &pluginSpec{
		Type: "autocmd",
		Name: event,
		Sync: isSync(fn),
		Opts: m,

		fn:       fn,
		smSuffix: fmt.Sprintf(":autocmd:%s:%s", event, pattern),
	})
}
