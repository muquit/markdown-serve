# Table Of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
  - [Homebrew (macOS and Linux)](#homebrew-macos-and-linux)
  - [Download binary](#download-binary)
  - [Build from source](#build-from-source)
- [Usage](#usage)
  - [Options](#options)
  - [Examples](#examples)
- [Live Reload](#live-reload)
- [Building from source](#building-from-source)
- [Dependencies](#dependencies)
- [License](#license)
- [Author](#author)

# Introduction

`markdown-serve` is a simple command line web server that serves a directory of Markdown files as browsable HTML. Just point it at a directory, it lists the `.md` files as a collapsible tree, and renders them in the browser when clicked. 

I prefer to use `vim` in a terminal to write Markdown files instead of [vscode](https://code.visualstudio.com/) and like to see how the output is rendered from anywhere when editing files remotely over [Tailscale](https://tailscale.com/). The server watches for file changes by default and reloads the browser automatically via Server-Sent Events.

# Features

- Lists `.md` files as a collapsible tree (directories and files)
- Renders Markdown as clean HTML with a readable layout
- Syntax highlighting via HIGHLIGHTJS loaded from CDN
- GitHub-flavored Markdown extensions: tables, strikethrough, task lists, fenced code blocks, auto-heading IDs
- Live reload via Server-Sent Events when files change on disk (on by default)
- Recursive directory support with empty directory pruning
- Path traversal protection
- Binds to `0.0.0.0` by default so you can access it remotely

# Requirements

- Go 1.22 or later

# Installation

## Homebrew (macOS and Linux)

```
brew tap muquit/markdown-serve
brew install markdown-serve
```

## Download binary

Download a pre-built binary for your platform from MARKDOWN_SERVE_RELEASES.

## Build from source

```
git clone https://github.com/muquit/markdown-serve.git
cd markdown-serve
make
```

# Usage

```
markdown-serve [options] [directory]
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

When `-watch` is enabled (the default), the server uses FSNOTIFY to watch the served directory tree for changes. When a `.md` file is written, the browser reloads automatically via a Server-Sent Events connection at `/events`. No WebSocket or external tooling is needed.

To disable live reload:
```
markdown-serve -watch=false
```

# Building from source

A Makefile is provided. It reads the version from the `VERSION` file and stamps it into the binary at compile time via `-ldflags`.

```
make          # build the binary
make clean    # remove the binary
make docs     # regenerate README.md from docs/README.md
```

# Dependencies

- GOMARKDOWN for Markdown rendering
- FSNOTIFY for filesystem change detection
- HIGHLIGHTJS loaded from CDN for syntax highlighting at render time

The project uses the Go standard library for HTTP serving and has no other runtime dependencies.

# License

MIT

# Author

Built with [Claude Code](https://claude.ai/code). Look at [CLAUDE.md](CLAUDE.md) for the prompt.


---
<sub>TOC is created by https://github.com/muquit/markdown-toc-go on Jun-16-2026</sub>
