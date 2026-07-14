#!/usr/bin/env bash
#
# Generate docs/usage.md from the CLI's own -h output, so the flag list
# and options table in the docs never drift from main.go. Meant to be
# included from docs/README.md with markdown-toc-go's
# @[:markdown](usage.md) directive during doc generation.
#
# muquit@gmail.com

set -euo pipefail

cd "$(dirname "$0")/.."

OUT_FILE="docs/usage.md"
TMP_BIN=$(mktemp -t markdown-serve-gen-usage.XXXXXX)
trap 'rm -f "$TMP_BIN"' EXIT

go build -o "$TMP_BIN" .

# flag.Usage() writes "Usage of <argv0>:" followed by one block per flag:
#   -name[ type]
#   \tDescription (default X)
FLAG_BODY=$("$TMP_BIN" -h 2>&1 | tail -n +2)

# Build the Options table from the same flag blocks.
TABLE=$(echo "$FLAG_BODY" | awk '
    /^  -/ {
        if (name != "") print_row()
        line = $0
        sub(/^  -/, "", line)
        split(line, parts, " ")
        name = parts[1]
        next
    }
    /^    \t/ {
        desc = $0
        sub(/^    \t/, "", desc)
        next
    }
    END { if (name != "") print_row() }
    function print_row(   def, defcell) {
        def = ""
        if (match(desc, /\(default [^)]*\)/)) {
            def = substr(desc, RSTART + 9, RLENGTH - 10)
            desc = substr(desc, 1, RSTART - 1) substr(desc, RSTART + RLENGTH)
            gsub(/^[ \t]+|[ \t]+$/, "", desc)
        }
        gsub(/"/, "", def)
        defcell = (def == "") ? "" : "`" def "`"
        printf("| `-%s` | %s | %s |\n", name, defcell, desc)
        name = ""
        desc = ""
    }
')

{
    echo "# Usage"
    echo
    echo '* Start `markdown-serve` server with a directory containing `.md` file'
    echo '* Point your browser to the URL e.g. http://localhost:8485'
    echo '* Edit Markdown files and the changes will be automatically refreshed in your'
    echo 'browser'
    echo
    echo '```'
    echo "➤ markdown-serve -h"
    echo "Usage of markdown-serve:"
    echo "$FLAG_BODY"
    echo '```'
    echo
    echo "If no directory is given, the current directory is used."
    echo
    echo "## Options"
    echo
    echo "| Option | Default | Description |"
    echo "|--------|---------|-------------|"
    echo "$TABLE"
    echo
    echo
    echo "## Examples"
    echo
    echo "Serve the current directory:"
    echo '```'
    echo "markdown-serve"
    echo '```'
    echo
    echo "Serve a specific directory on a custom port:"
    echo '```'
    echo "markdown-serve -port 8485 ~/notes"
    echo '```'
    echo
    echo "Serve without live reload:"
    echo '```'
    echo "markdown-serve -watch=false /path/to/docs"
    echo '```'
    echo
    echo "Restrict to localhost only:"
    echo '```'
    echo "markdown-serve -host 127.0.0.1"
    echo '```'
    echo
    echo "Serve in dark mode:"
    echo '```'
    echo "markdown-serve -dark ~/notes"
    echo '```'
} > "$OUT_FILE"

echo ">>>> Wrote $OUT_FILE"
