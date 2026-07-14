# Usage

* Start `markdown-serve` server with a directory containing `.md` file
* Point your browser to the URL e.g. http://localhost:8485
* Edit Markdown files and the changes will be automatically refreshed in your
browser

```
➤ markdown-serve -h
Usage of markdown-serve:
  -dark
    	Render pages in dark mode
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
| `-dark` |  | Render pages in dark mode |
| `-host` | `0.0.0.0` | Host to bind to |
| `-port` | `8485` | Port to listen on |
| `-version` |  | Print version and exit |
| `-watch` | `true` | Reload browser on file changes |


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

Serve in dark mode:
```
markdown-serve -dark ~/notes
```
