# Introduction

`markdown-serve` is a cross-platform command-line web server that serves a directory of Markdown files as browsable HTML. Just point it at a directory, it lists the `.md` files as a collapsible tree, and renders them in the browser when clicked. 

I prefer `vi/vim` in a terminal over @VSCODE@ adns such for writing Markdown, and 
like to see how the output is rendered in the browser.  The server watches 
for file changes by default and reloads the browser automatically via @SERVER_SENT_EVENTS@. 

I find it much more pleasurable to work that way. Hope you find it useful as well.

Suggestions, pull requests are welcome but please keep in mind that I like to keep things simple.


# Features

- Lists `.md` files as a collapsible tree (directories and files)
- Renders Markdown as clean HTML
- Syntax highlighting via @HIGHLIGHTJS@ loaded from CDN
- GitHub-flavored Markdown extensions: tables, strikethrough, task lists, fenced code blocks, auto-heading IDs
- Live reload via @SERVER_SENT_EVENTS@ when files change on disk (on by default)
- Recursive directory support with empty directory pruning
- Path traversal protection
- Binds to `0.0.0.0` by default so you can access it remotely

# Accessing Home Network from anywhere

Whenever needed, I run the markdown-serve on a machine at home and
access it from anywhere over @TAILSCALE@ using a browser to see how the
Markdown is rendered as HTML. As long as both devices are on
the same @TAILSCALE@ network, it just works. Browse and edit Markdown
files remotely as if I were sitting at home.

@[:markdown](tailscale.md)

# Installation

## Homebrew (macOS and Linux)

```
brew tap muquit/markdown-serve
brew install markdown-serve
```

## Download binary

Download a pre-built binaries for your platform from @MARKDOWN_SERVE_RELEASES@
page.

## Build from source

Make sure @GO@ is installed. Look at @MAKEFILE@.

```
git clone https://github.com/muquit/markdown-serve.git
cd markdown-serve
make
```

# Synopsis

```
➤ markdown-serve -h
Usage of markdown-serve:
  -host string
    	Host to bind to (default "0.0.0.0")
  -port int
    	Port to listen on (default 8485)
  -version
    	Print version and exit
  -watch
    	Reload browser on file changes (default true)

```

If no directory is given, the current directory is used.

## Options

| Option | Default | Description |
|--------|---------|-------------|
| `-host` | `0.0.0.0` | Host to bind to |
| `-port` | `8485` | Port to listen on |
| `-watch` | `true` | Reload browser on file changes |
| `-version` | | Print version and exit |

## Examples

Serve the current directory:
```
markdown-serve
```

Serve a specific directory on a custom port:
```
markdown-serve -port 8485 ~/notes
```

Serve without live reload:
```
markdown-serve -watch=false /path/to/docs
```

Restrict to localhost only:
```
markdown-serve -host 127.0.0.1
```

Print version:
```
markdown-serve -version
```

# Live Reload

When `-watch` is enabled (the default), the server uses @FSNOTIFY@ to watch the served directory tree for changes. When a `.md` file is written, the browser reloads automatically via a @SERVER_SENT_EVENTS@ connection at `/events`. No WebSocket or external tooling is needed.

To disable live reload:
```
markdown-serve -watch=false
```

# Building from source

A @MAKEFILE@ is provided. It reads the version from the `VERSION` file and stamps it into the binary at compile time via `-ldflags`.

```
make          # build the binary
make clean    # remove the binary
make docs     # regenerate README.md from docs/README.md
```

# Dependencies

- @GOMARKDOWN@ for rendering Markdown as HTML
- @FSNOTIFY@ for filesystem change detection
- @HIGHLIGHTJS@ loaded from CDN for syntax highlighting at render time

The project uses the @GO@ standard library for HTTP serving.

# License

MIT

# Author

Built with @CLAUDE_CODE@. Look at [CLAUDE.md](CLAUDE.md) for the prompt.

