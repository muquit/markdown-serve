# Table Of Contents
- [v1.0.7](#v107)
- [v1.0.6](#v106)
- [v1.0.5](#v105)
- [v1.0.4](#v104)
- [v1.0.3](#v103)
- [v1.0.2](#v102)
- [v1.0.1](#v101)


# v1.0.7

* Add rendering support for github-style alerts (`[!NOTE]`, `[!TIP]`,
  `[!IMPORTANT]`, `[!WARNING]`, `[!CAUTION]`).

> [!NOTE]
> Highlights information that users should take into account, even when skimming.

> [!TIP]
> Optional information to help a user be more successful.

> [!IMPORTANT]
> Crucial information necessary for users to succeed.

> [!WARNING]
> Critical content demanding immediate user attention due to potential risks.

> [!CAUTION]
> Negative potential consequences of an action.

(Aug-31-2026)

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
