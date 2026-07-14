#!/usr/bin/env bash
#
# Generate release_notes.md from the archives built by go-xbuild-go in
# bin/, for the version listed in VERSION.
#
# muquit@gmail.com

set -euo pipefail

BIN_DIR="bin"
VERSION_FILE="VERSION"
OUT_FILE="release_notes.md"

if [ ! -f "$VERSION_FILE" ]; then
    echo "*** ERROR: $VERSION_FILE not found" >&2
    exit 1
fi
VERSION=$(cat "$VERSION_FILE")

ZIP_FILE="$BIN_DIR/markdown-serve-${VERSION}-windows-amd64.d.zip"
TAR_FILE="$BIN_DIR/markdown-serve-${VERSION}-linux-amd64.d.tar.gz"

for f in "$ZIP_FILE" "$TAR_FILE"; do
    if [ ! -f "$f" ]; then
        echo "*** ERROR: $f not found, build the release archives first (make release or go-xbuild-go)" >&2
        exit 1
    fi
done

{
    echo "# Release ${VERSION}"
    echo
    echo "Please look at [ChangeLog.md](ChangeLog.md) for details on what has changed in this version."
    echo
    echo
    echo '```'
    echo "➤ unzip -l $(basename "$ZIP_FILE")"
    unzip -l "$ZIP_FILE"
    echo '```'
    echo
    echo '```'
    echo "➤ tar -tvf $(basename "$TAR_FILE")"
    tar -tvf "$TAR_FILE"
    echo '```'
    echo
    echo 'Copy the binary as `markdown-serve` or `markdown-serve.exe` to somewhere in your PATH.'
    echo
    echo "Please look at [Installation](README.md#installation) section in [README.md](README.md) for info on installing via Homebrew or compiling from source etc."
    echo
    echo
    echo "Cross-compiled and released with [go-xbuild-go](https://github.com/muquit/go-xbuild-go)"
} > "$OUT_FILE"

echo ">>>> Wrote $OUT_FILE for $VERSION"
