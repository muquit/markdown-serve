# Table Of Contents
- [v1.0.6](#v106)
- [v1.0.5](#v105)
- [v1.0.4](#v104)
- [v1.0.3](#v103)
- [v1.0.2](#v102)
- [v1.0.1](#v101)


# v1.0.6

* Fix Markdown pages occasionally loading indefinitely in the browser. Live
  reload now uses lightweight revision polling instead of holding one SSE
  connection open per tab, which could exhaust Chrome's HTTP/1.1 connection
  pool.
* Load Highlight.js after the document has finished loading so CDN delays do
  not block the initial page render.

The bug was indentified and fixed by [Codex](https://openai.com/codex)


(Aug-25-2026)

# v1.0.5
* Add `mdsr.sh` in brew formula file. Uses `go-xbuilg-go v1.0.10`. No code
  change.
(Aug-10-2026)

# v1.0.4
* Add project link at the Rendered by line at the top.

(Aug-09-2026)

# v1.0.3

* Add copy-to-clipboard icon button on code blocks
* Reduce frame thickness around code blocks

(Jul-24-2026)

# v1.0.2

* Fix blockquotes rendering without any visual styling
* Add `-dark` flag to render pages in dark mode

(Jul-13-2026)

# v1.0.1

* Initial release

(Jun-26-2026)
