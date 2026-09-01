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
