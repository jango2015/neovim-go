// Copyright 2016 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plugin

import (
	"bytes"
	"strings"
	"testing"
)

func testSync() error { return nil }
func testAsync()      {}

var expectedScript = `if exists(g:loaded_zztestzz)
  finish
endif
let g:loaded_zztestzz = 1
autocmd! zztestzz_0 BufRead *.go nested call rpcrequest(s:channel(), ":autocmd:BufRead:*.go", eval("cwd()"))
command! -nargs=1 -range -bang -bar -register Command1 call rpcrequest(s:channel(), ":command:Command1", [<f-args>], [<line1>, <line2>], <q-bang> == "!", <q-reg>, eval("cwd()"))
function! FooAsync (...)
  return rpcnotify(s:channel(), ":function:FooAsync", a:000, eval("printf(\"hello\nworld\")"))
endfunction
function! FooSync (...)
  return rpcrequest(s:channel(), ":function:FooSync", a:000, eval("cwd()"))
endfunction
`

func TestPluginScript(t *testing.T) {
	savePluginSpecs := pluginSpecs
	pluginSpecs = nil
	defer func() {
		pluginSpecs = savePluginSpecs
	}()

	HandleAutocmd("BufRead", &AutocmdOptions{Pattern: "*.go", Nested: true, Eval: "cwd()"}, testSync)
	HandleCommand("Command1", &CommandOptions{NArgs: "1", Range: ".", Bang: true, Register: true, Eval: "cwd()", Bar: true}, testSync)
	HandleFunction("FooSync", &FunctionOptions{Eval: "cwd()"}, testSync)
	HandleFunction("FooAsync", &FunctionOptions{Eval: `printf("hello\nworld")`}, testAsync)

	var buf bytes.Buffer
	if err := writePluginScript(&buf, "zztestzz"); err != nil {
		t.Fatal(err)
	}

	actual := buf.String()
	if actual != expectedScript {
		alines := strings.Split(actual, "\n")
		elines := strings.Split(expectedScript, "\n")
		alines = append(alines, "EOF")
		elines = append(elines, "EOF")
		n := len(alines)
		if len(elines) < n {
			n = len(elines)
		}
		for i := 0; i < n; i++ {
			if alines[i] != elines[i] {
				t.Fatalf("difference at line %d\n  actual: %q\nexpected: %q", i+1, alines[i], elines[i])
			}
		}
	}
}
