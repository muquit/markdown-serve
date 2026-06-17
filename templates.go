package main

import (
	"fmt"
	"html/template"
	"io"
	"strings"
)

var indexTmpl = template.Must(template.New("index").Parse(`
{{define "node"}}
{{- if .IsDir}}
<li>
  <details open>
    <summary>📁 {{.Name}}/</summary>
    <ul>{{range .Children}}{{template "node" .}}{{end}}</ul>
  </details>
</li>
{{- else}}
<li><a href="/{{.Path}}">{{.Name}}</a></li>
{{- end}}
{{end}}
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>{{ .Dir }}</title>
  <style>
    body { max-width: 800px; margin: 40px auto; font-family: sans-serif; padding: 0 20px; color: #333; }
    h1 { font-size: 1.2em; color: #555; word-break: break-all; }
    ul { list-style: none; padding-left: 1.5em; margin: 2px 0; }
    ul.tree { padding-left: 0; }
    li { padding: 2px 0; }
    summary { cursor: pointer; color: #333; padding: 3px 0; user-select: none; list-style: none; }
    summary::-webkit-details-marker { display: none; }
    summary::marker { display: none; }
    summary:hover { color: #000; }
    a { text-decoration: none; color: #0366d6; }
    a:hover { text-decoration: underline; }
    footer { margin-top: 24px; font-size: 0.85em; color: #888; }
  </style>
</head>
<body>
  <h1>{{ .Dir }}</h1>
  <ul class="tree">
    {{range .Root.Children}}{{template "node" .}}{{end}}
  </ul>
  <footer>{{ .Count }} file{{ if ne .Count 1 }}s{{ end }}</footer>
  {{if .Watch}}<script>
const es = new EventSource('/events');
es.onmessage = () => location.reload();
es.onerror = () => { es.close(); setTimeout(() => location.reload(), 2000); };
</script>{{end}}
</body>
</html>
`))

var docTmpl = template.Must(template.New("doc").Parse(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>{{ .Title }}</title>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github.min.css">
  <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
  <script>hljs.highlightAll();</script>
  <style>
    body { max-width: 860px; margin: 40px auto; font-family: Georgia, serif; padding: 0 20px; color: #333; line-height: 1.7; }
    pre { background: #f6f8fa; padding: 16px; border-radius: 4px; overflow-x: auto; }
    code { font-family: monospace; font-size: 0.9em; }
    pre code { background: none; padding: 0; }
    a { color: #0366d6; }
    nav { display: flex; justify-content: space-between; align-items: baseline; margin-bottom: 24px; font-size: 0.9em; }
    .modtime { color: #888; font-size: 0.85em; }
    table { border-collapse: collapse; width: 100%; }
    th, td { border: 1px solid #ddd; padding: 8px 12px; }
    th { background: #f6f8fa; }
    img { max-width: 100%; }
  </style>
</head>
<body>
  <nav><a href="/">&#8592; Back</a><span class="modtime">Last modified: {{ .ModTime }}</span></nav>
  {{ .Body }}
  {{if .Watch}}<script>
const es = new EventSource('/events');
es.onmessage = () => location.reload();
es.onerror = () => { es.close(); setTimeout(() => location.reload(), 2000); };
</script>{{end}}
</body>
</html>
`))

func renderIndex(w io.Writer, dir string, root *TreeNode, count int, watch bool) {
	data := struct {
		Dir   string
		Root  *TreeNode
		Count int
		Watch bool
	}{
		Dir:   dir,
		Root:  root,
		Count: count,
		Watch: watch,
	}
	if err := indexTmpl.Execute(w, data); err != nil {
		fmt.Fprintf(w, "template error: %v", err)
	}
}

func renderDocument(w io.Writer, title string, body []byte, modTime string, watch bool) {
	display := strings.TrimSuffix(title, ".md")
	data := struct {
		Title   string
		Body    template.HTML
		ModTime string
		Watch   bool
	}{
		Title:   display,
		Body:    template.HTML(body),
		ModTime: modTime,
		Watch:   watch,
	}
	if err := docTmpl.Execute(w, data); err != nil {
		fmt.Fprintf(w, "template error: %v", err)
	}
}
