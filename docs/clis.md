# Helper CLIs

`markdown-serve` doesn't track whether another instance is already
running. Start a second one on the same port and it just fails with
`address already in use`. To avoid having to manually find and kill a
stale process before starting a new one, two small wrapper scripts are
included:

- `mdsr.sh` for macOS/Linux
- `mdsr.ps1` for Windows (PowerShell)

Both take exactly the same arguments as `markdown-serve` itself. Before
starting a new instance, they look for an already running
`markdown-serve` process, print what it was serving, kill it, and then
start the new one in its place.

## Usage

The scripts accept the same options and directory argument as `markdown-serve`. Copy them somewhere in your PATH.

macOS/Linux:
```
mdsr.sh -port 8485 ~/notes
```

Windows (PowerShell):
```
mdsr.ps1 -port 8485 C:\notes
```

If `markdown-serve` was already running, you'll see something like:
```
Killing running markdown-serve (pid 79299): markdown-serve /Users/muquit/notes
Starting: /Users/muquit/bin/markdown-serve -port 8485 /Users/muquit/notes
```

`mdsr` is short for mark**d**own **s**erve **r**estart. Easier to type than the
full name when running it often.

The scripts locate `markdown-serve` on `PATH`. To point at a binary
that isn't on `PATH`, set the `MARKDOWN_SERVE_BIN` environment variable.
Both scripts respect it:

```
MARKDOWN_SERVE_BIN=/path/to/markdown-serve mdsr.sh ~/notes
```

```
$env:MARKDOWN_SERVE_BIN = "C:\path\to\markdown-serve.exe"
mdsr.ps1 C:\notes
```

