# markdown-serve

A lightweight Go web server that serves a directory of Markdown files as
browsable HTML.

## Project Overview

`markdown-serve` takes a directory path as an argument, lists all `.md` files
in that directory, and renders them as HTML when clicked. It is a local
developer/personal tool, not a production web server.

## Repository

`github.com/muquit/markdown-serve`

## Tech Stack

- **Language**: Go (latest stable)
- **Markdown library**: `github.com/gomarkdown/markdown`
- **HTTP**: Go standard library `net/http`
- **No frameworks**: keep dependencies minimal

## Features

- Accept a directory path as a CLI argument (default: current directory)
- List all `.md` files in the directory on the index page
- Clicking a file renders it as HTML in the browser
- Clean, readable HTML output with basic CSS styling
- Syntax highlighting via `highlight.js` (CDN)
- GitHub-flavored Markdown extensions: tables, strikethrough, task lists,
  fenced code blocks, auto-heading IDs
- Optional `--watch` flag: live reload via SSE when `.md` files change on disk

## CLI Usage

```
markdown-serve [options] [directory]

Options:
  -port int    Port to listen on (default 8485)
  -host string Host to bind to (default "0.0.0.0")
  -watch       Live reload browser on file changes (default true)

Arguments:
  directory    Path to directory containing .md files (default: current dir)

Examples:
  markdown-serve .
  markdown-serve /path/to/docs
  markdown-serve -port 9090 ~/notes
```

## Project Structure

```
markdown-serve/
├── main.go          # Entry point, CLI flag parsing, server startup
├── server.go        # HTTP handlers: index listing, markdown rendering
├── render.go        # Markdown-to-HTML rendering logic
├── templates.go     # Inline HTML templates (index page, document page)
├── go.mod
├── go.sum
├── README.md
└── CLAUDE.md
```

## Code Style

- Standard `gofmt` formatting (K&R-style braces, tabs for indentation)
- Explicit `return` statements always
- `if err != nil` error handling, no panic in handlers
- No global mutable state; pass config struct to handlers
- All handlers use `http.HandlerFunc` signature

## HTTP Routes

| Route      | Description                                      |
|------------|--------------------------------------------------|
| `GET /`    | Lists all `.md` files in the served directory   |
| `GET /:file` | Renders the named `.md` file as HTML           |

- Return `404` with a clear message if a file is not found
- Return `400` if the path attempts directory traversal (`..`)
- Only serve files with `.md` extension; reject all others

## Rendering

- Use `github.com/gomarkdown/markdown` with `parser.CommonExtensions | parser.AutoHeadingIDs`
- Renderer flags: `html.CommonFlags | html.HrefTargetBlank`
- Wrap rendered output in a full HTML page with:
  - `<meta charset="utf-8">`
  - Inline CSS for readable body (max-width, font, padding)
  - `highlight.js` from CDN for code syntax highlighting
  - A navigation link back to the index (`← Back`)

## Index Page

- Show the directory path at the top
- List `.md` files sorted alphabetically
- Each entry is a clickable link
- Show file count at the bottom
- Simple clean styling, no JavaScript required

## Build & Release

- Build with `go build -o markdown-serve .`
- Use `go-xbuild-go` for cross-platform releases
- Homebrew formula generated via `go-xbuild-go`

## Security Note

The server binds to `0.0.0.0` by default to allow remote access (e.g., accessing
markdown files from outside the local network). Use `-host 127.0.0.1` to restrict
to localhost only. Since this serves local files, ensure the machine's firewall
restricts access as needed.

## Dependencies

```
github.com/gomarkdown/markdown
github.com/fsnotify/fsnotify
```

`highlight.js` loaded from CDN at render time.

## Non-Goals

- No file upload
- No editing
- No authentication

