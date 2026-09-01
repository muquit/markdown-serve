# Live Reload

When `-watch` is enabled (the default), the server uses @FSNOTIFY@ to watch the served directory tree for changes. When a `.md` file is written, the browser reloads automatically via a @SERVER_SENT_EVENTS@ connection at `/events`. No WebSocket or external tooling is needed.

To disable live reload:
```
markdown-serve -watch=false
```
