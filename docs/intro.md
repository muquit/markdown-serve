# Introduction

`markdown-serve` is a cross-platform command-line web server that serves a directory of Markdown files as browsable HTML. Just point it at a directory, it lists the `.md` files as a collapsible tree, and renders them in the browser when clicked. 

I prefer `vi/vim` in a terminal over @VSCODE@ and such for writing Markdown, and 
like to see how the output is rendered in the browser.  The server watches 
for file changes by default and reloads the browser automatically via @SERVER_SENT_EVENTS@. 

I find it much more pleasurable to work that way. Hope you find it useful as well.

Suggestions, pull requests are welcome but please keep in mind that I like to keep things simple.
