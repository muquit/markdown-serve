# Features

- Lists `.md` files as a collapsible tree (directories and files)
- Renders Markdown as clean HTML
- Syntax highlighting via @HIGHLIGHTJS@ loaded from CDN
- GitHub-flavored Markdown extensions: tables, strikethrough, task lists, fenced code blocks, auto-heading IDs
- Live reload via @SERVER_SENT_EVENTS@ when files change on disk (on by default)
- Recursive directory support with empty directory pruning
- Path traversal protection
- Binds to `0.0.0.0` by default so you can access it remotely
- Dark mode via the `-dark` flag
- `mdsr.sh`/`mdsr.ps1` helper scripts to restart `markdown-serve` without manually killing a running instance first
