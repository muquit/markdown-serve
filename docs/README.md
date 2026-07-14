# Introduction

`markdown-serve` is a cross-platform command-line web server that serves a directory of Markdown files as browsable HTML. Just point it at a directory, it lists the `.md` files as a collapsible tree, and renders them in the browser when clicked. 

I prefer `vi/vim` in a terminal over @VSCODE@ and such for writing Markdown, and 
like to see how the output is rendered in the browser.  The server watches 
for file changes by default and reloads the browser automatically via @SERVER_SENT_EVENTS@. 

I find it much more pleasurable to work that way. Hope you find it useful as well.

Suggestions, pull requests are welcome but please keep in mind that I like to keep things simple.

# Latest Version (v1.0.2 - Jul-13-2026)

The latest version is v1.0.2 Please look at @CHANGELOG@ for details.

@[:markdown](features.md)


# Installation

## Download pre-compiled binaries

Download a pre-built binary for your platform from @MARKDOWN_SERVE_RELEASES@ page.

Extract the archive and copy the binary as `markdown-serve` (Linux/macOS) or
`markdown-serve.exe` (Windows) to somewhere in your PATH.

@[:markdown](brew_install.md)

## Building from source

Make sure @GO@ is installed.

```
git clone https://github.com/muquit/markdown-serve.git
cd markdown-serve
go build .
```
or Look at @MAKEFILE@ and type:
```
make
```
Requires @GO_XBUILD_GO@ for compiling cross-platform binaries

@[:markdown](usage.md)

@[:markdown](clis.md)

# Accessing Home Network from anywhere

Whenever needed, I run  `markdown-serve` on a machine at home and
access it from anywhere over @TAILSCALE@ using a browser to see how the
Markdown is rendered as HTML. As long as both devices are on
the same @TAILSCALE@ network, it just works. Browse and edit Markdown
files remotely as if I were sitting at home.

## Tailscale

@[:markdown](tailscale.md)

# Live Reload

When `-watch` is enabled (the default), the server uses @FSNOTIFY@ to watch the served directory tree for changes. When a `.md` file is written, the browser reloads automatically via a @SERVER_SENT_EVENTS@ connection at `/events`. No WebSocket or external tooling is needed.

To disable live reload:
```
markdown-serve -watch=false
```

# Dependencies

- @GOMARKDOWN@ for rendering Markdown as HTML
- @FSNOTIFY@ for filesystem change detection
- @HIGHLIGHTJS@ loaded from CDN for syntax highlighting at render time

The project uses the @GO@ standard library for HTTP serving.

# Screenshots

Listing of all Markdown files:

![alt markdown-server ss1](images/ss1.png)

Clicked on the [README.md](README.md)

![alt markdown-server ss2](images/ss2.png)



# License

MIT. Look at @LICENSE@ for details.

# Author

Built with the help from @CLAUDE_CODE@. Look at [CLAUDE.md](CLAUDE.md) for the prompt used
for implementation.

