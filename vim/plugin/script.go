// Copyright 2016 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plugin

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type byServiceMethod []*pluginSpec

func (a byServiceMethod) Len() int           { return len(a) }
func (a byServiceMethod) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byServiceMethod) Less(i, j int) bool { return a[i].ServiceMethod < a[j].ServiceMethod }

func stringifyOpts(opts map[string]string, keys ...string) string {
	var p []byte
	for _, k := range keys {
		v, ok := opts[k]
		if ok {
			p = append(p, " -"...)
			p = append(p, k...)
			if v != "" {
				p = append(p, '=')
				p = append(p, v...)
			}
		}
	}
	return string(p)
}

func writePluginScript(w io.Writer, name string) error {
	// Sort for consistent order on output.
	sort.Sort(byServiceMethod(pluginSpecs))

	escapeEval := strings.NewReplacer(`"`, `\"`).Replace

	fmt.Fprintf(w, `if exists(g:loaded_%s)
  finish
endif
let g:loaded_%s = 1
`, name, name)

	for i, s := range pluginSpecs {
		args := []string{"s:channel()", `"` + s.ServiceMethod + `"`}
		rpcCall := "rpcnotify"
		if s.Sync {
			rpcCall = "rpcrequest"
		}
		switch s.Type {
		case "autocmd":
			group, ok := s.Opts["group"]
			if !ok {
				group = fmt.Sprintf("%s_%d", name, i)
			}
			pattern, ok := s.Opts["pattern"]
			if !ok {
				pattern = "*"
			}
			nested := ""
			if _, ok := s.Opts["nested"]; ok {
				nested = " nested"
			}
			if e, ok := s.Opts["eval"]; ok {
				args = append(args, fmt.Sprintf(`eval("%s")`, escapeEval(e)))
			}
			fmt.Fprintf(w, "autocmd! %s %s %s%s call %s(%s)\n",
				group, s.Name, pattern, nested, rpcCall, strings.Join(args, ", "))
		case "command":
			if _, ok := s.Opts["nargs"]; ok {
				args = append(args, "[<f-args>]")
			}
			if r, ok := s.Opts["range"]; ok {
				if r == "" || r == "%" {
					args = append(args, "[<line1>, <line2>]")
				} else if '0' <= r[0] && r[0] <= '9' {
					args = append(args, "<count>")
				}
			} else if _, ok := s.Opts["count"]; ok {
				args = append(args, "<count>")
			}
			if _, ok := s.Opts["bang"]; ok {
				args = append(args, `<q-bang> == "!"`)
			}
			if _, ok := s.Opts["register"]; ok {
				args = append(args, "<q-reg>")
			}
			if e, ok := s.Opts["eval"]; ok {
				args = append(args, fmt.Sprintf(`eval("%s")`, escapeEval(e)))
			}
			fmt.Fprintf(w, "command!%s %s call %s(%s)\n",
				stringifyOpts(s.Opts, "nargs", "complete", "range", "count", "bang", "bar", "register"),
				s.Name,
				rpcCall,
				strings.Join(args, ", "))
		case "function":
			args = append(args, "a:000")
			if e, ok := s.Opts["eval"]; ok {
				args = append(args, fmt.Sprintf(`eval("%s")`, escapeEval(e)))
			}
			fmt.Fprintf(w, "function! %s (...)\n  return %s(%s)\nendfunction\n", s.Name, rpcCall, strings.Join(args, ", "))
		default:
			panic("unknown spec type " + s.Type)
		}
	}
	return nil
}
