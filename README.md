### Nvimgo

Nvimgo is a [Neovim](https://neovim.io/) remote plugin for writing Go programs.
The plugin is implemented in Go.

Nvimgo is a work in progress. Use it with caution.

### Setup

Add the following code to init.vim:

    function! s:RequireNvimgo(host) abort
      let args = []
      let plugins = remote#host#PluginsForHost(a:host.name)
      for plugin in plugins
        call add(args, plugin.path)
      endfor
      return rpcstart('nvimgo', args)
    endfunction

    call remote#host#Register('nvimgo', '*', function('s:RequireNvimgo'))

Add an empty file to rplugin/nvimgo/x

Run `:UpdateRemotePlugins`.

