#!/usr/bin/env bash
# Restart wrapper for markdown-serve: kills any already-running instance
# (by resolved binary path) and starts a new one with the given args.
set -euo pipefail

BINARY="${MARKDOWN_SERVE_BIN:-markdown-serve}"

if ! command -v "$BINARY" >/dev/null 2>&1; then
    echo "error: '$BINARY' not found in PATH (set MARKDOWN_SERVE_BIN to override)" >&2
    exit 1
fi

resolved="$(command -v "$BINARY")"
proc_name="$(basename "$resolved")"

RED='\033[0;31m'
NC='\033[0m'

find_pids() {
    pgrep -x "$proc_name" 2>/dev/null | grep -v "^$$\$" || true
}

pids="$(find_pids)"

if [ -n "$pids" ]; then
    for pid in $pids; do
        cmd="$(ps -o args= -p "$pid" 2>/dev/null || true)"
        printf "${RED}Killing running markdown-serve (pid %s): %s${NC}\n" "$pid" "$cmd"
    done
    kill $pids 2>/dev/null || true

    for _ in $(seq 1 20); do
        [ -z "$(find_pids)" ] && break
        sleep 0.2
    done

    pids="$(find_pids)"
    if [ -n "$pids" ]; then
        echo "Force killing: $pids"
        kill -9 $pids 2>/dev/null || true
    fi
fi

echo "Starting: $resolved $*"
exec "$resolved" "$@"
