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
    body.dark { background: #0d1117; color: #c9d1d9; }
    body.dark h1 { color: #8b949e; }
    body.dark summary { color: #c9d1d9; }
    body.dark summary:hover { color: #fff; }
    body.dark a { color: #58a6ff; }
    body.dark footer { color: #8b949e; }
  </style>
</head>
<body{{if .Dark}} class="dark"{{end}}>
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
  {{if .Dark}}<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github-dark.min.css">{{else}}<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github.min.css">{{end}}
  <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
  <script>hljs.highlightAll();</script>
  <style>
    body { max-width: 860px; margin: 40px auto; font-family: Georgia, serif; padding: 0 20px; color: #333; line-height: 1.7; }
    .code-wrap { position: relative; }
    pre { background: #f6f8fa; padding: 4px 6px; border-radius: 4px; overflow-x: auto; }
    code { font-family: monospace; font-size: 0.9em; }
    pre code { background: none; padding: 0; }
    .copy-btn { position: absolute; top: 6px; right: 6px; width: 28px; height: 28px; display: flex; align-items: center; justify-content: center; padding: 0; border: 1px solid rgba(27,31,36,0.15); border-radius: 4px; background: #fff; color: #57606a; cursor: pointer; }
    .copy-btn:hover { background: #f3f4f6; }
    .copy-btn.copied { background: #2da44e; color: #fff; border-color: #2da44e; }
    a { color: #0366d6; }
    nav { display: flex; justify-content: space-between; align-items: baseline; margin-bottom: 24px; font-size: 0.9em; }
    .modtime { color: #888; font-size: 0.85em; }
    table { border-collapse: collapse; width: 100%; }
    th, td { border: 1px solid #ddd; padding: 8px 12px; }
    th { background: #f6f8fa; }
    img { max-width: 100%; }
    blockquote { margin: 0 0 16px; padding: 0 1em; color: #6a737d; border-left: 0.25em solid #dfe2e5; }
    blockquote > :first-child { margin-top: 0; }
    blockquote > :last-child { margin-bottom: 0; }
    body.dark { background: #0d1117; color: #c9d1d9; }
    body.dark pre { background: #161b22; }
    body.dark a { color: #58a6ff; }
    body.dark .modtime { color: #8b949e; }
    body.dark th, body.dark td { border-color: #30363d; }
    body.dark th { background: #161b22; }
    body.dark blockquote { color: #8b949e; border-left-color: #30363d; }
    body.dark .copy-btn { background: #21262d; color: #8b949e; border-color: rgba(240,246,252,0.1); }
    body.dark .copy-btn:hover { background: #30363d; }
    body.dark .copy-btn.copied { background: #238636; color: #fff; border-color: #238636; }
  </style>
</head>
<body{{if .Dark}} class="dark"{{end}}>
  <nav><a href="/">&#8592; Back</a><span class="modtime">{{ if .Version }}Rendered by markdown-serve {{ .Version }} &middot; {{ end }}Last modified: {{ .ModTime }}</span></nav>
  {{ .Body }}
  <script>
function copyText(text) {
  if (navigator.clipboard && window.isSecureContext) {
    return navigator.clipboard.writeText(text);
  }
  const ta = document.createElement('textarea');
  ta.value = text;
  ta.style.position = 'fixed';
  ta.style.opacity = '0';
  document.body.appendChild(ta);
  ta.select();
  document.execCommand('copy');
  document.body.removeChild(ta);
  return Promise.resolve();
}
const copyIcon = '<svg viewBox="0 0 16 16" width="16" height="16" fill="currentColor"><path d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 0 1 0 1.5h-1.5a.25.25 0 0 0-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 0 0 .25-.25v-1.5a.75.75 0 0 1 1.5 0v1.5A1.75 1.75 0 0 1 9.25 16h-7.5A1.75 1.75 0 0 1 0 14.25Z"></path><path d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0 1 14.25 11h-7.5A1.75 1.75 0 0 1 5 9.25Zm1.75-.25a.25.25 0 0 0-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 0 0 .25-.25v-7.5a.25.25 0 0 0-.25-.25Z"></path></svg>';
const checkIcon = '<svg viewBox="0 0 16 16" width="16" height="16" fill="currentColor"><path d="M13.78 4.22a.75.75 0 0 1 0 1.06l-7.25 7.25a.75.75 0 0 1-1.06 0L2.22 9.28a.75.75 0 0 1 1.06-1.06L6 10.94l6.72-6.72a.75.75 0 0 1 1.06 0Z"></path></svg>';
document.querySelectorAll('pre code').forEach((code) => {
  const pre = code.parentElement;
  const wrap = document.createElement('div');
  wrap.className = 'code-wrap';
  pre.parentNode.insertBefore(wrap, pre);
  wrap.appendChild(pre);
  const btn = document.createElement('button');
  btn.type = 'button';
  btn.className = 'copy-btn';
  btn.setAttribute('aria-label', 'Copy code');
  btn.innerHTML = copyIcon;
  btn.addEventListener('click', () => {
    copyText(code.textContent).then(() => {
      btn.innerHTML = checkIcon;
      btn.classList.add('copied');
      setTimeout(() => {
        btn.innerHTML = copyIcon;
        btn.classList.remove('copied');
      }, 2000);
    });
  });
  wrap.appendChild(btn);
});
</script>
  {{if .Watch}}<script>
const es = new EventSource('/events');
es.onmessage = () => location.reload();
es.onerror = () => { es.close(); setTimeout(() => location.reload(), 2000); };
</script>{{end}}
</body>
</html>
`))

func renderIndex(w io.Writer, dir string, root *TreeNode, count int, watch bool, dark bool) {
	data := struct {
		Dir   string
		Root  *TreeNode
		Count int
		Watch bool
		Dark  bool
	}{
		Dir:   dir,
		Root:  root,
		Count: count,
		Watch: watch,
		Dark:  dark,
	}
	if err := indexTmpl.Execute(w, data); err != nil {
		fmt.Fprintf(w, "template error: %v", err)
	}
}

func renderDocument(w io.Writer, title string, body []byte, modTime string, watch bool, version string, dark bool) {
	display := strings.TrimSuffix(title, ".md")
	data := struct {
		Title   string
		Body    template.HTML
		ModTime string
		Watch   bool
		Version string
		Dark    bool
	}{
		Title:   display,
		Body:    template.HTML(body),
		ModTime: modTime,
		Watch:   watch,
		Version: version,
		Dark:    dark,
	}
	if err := docTmpl.Execute(w, data); err != nil {
		fmt.Fprintf(w, "template error: %v", err)
	}
}
