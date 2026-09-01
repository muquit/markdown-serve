package main

import (
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func renderMarkdown(src []byte) []byte {
	extensions := parser.CommonExtensions
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(src)

	assignGitHubHeadingIDs(doc)
	alerts := markGitHubAlerts(doc)

	flags := html.UseXHTML | html.HrefTargetBlank
	opts := html.RendererOptions{
		Flags:          flags,
		RenderNodeHook: githubAlertHook(alerts),
	}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

// githubAlertMarkers maps a GitHub-flavored alert marker to its CSS/label
// name. Matching is case-sensitive to mirror GitHub's own parser.
var githubAlertMarkers = map[string]string{
	"[!NOTE]":      "note",
	"[!TIP]":       "tip",
	"[!IMPORTANT]": "important",
	"[!WARNING]":   "warning",
	"[!CAUTION]":   "caution",
}

var githubAlertTitles = map[string]string{
	"note":      "Note",
	"tip":       "Tip",
	"important": "Important",
	"warning":   "Warning",
	"caution":   "Caution",
}

var githubAlertIcons = map[string]string{
	"note":      `<svg viewBox="0 0 16 16" width="16" height="16" fill="currentColor"><path d="M0 8a8 8 0 1 1 16 0A8 8 0 0 1 0 8Zm8-6.5a6.5 6.5 0 1 0 0 13 6.5 6.5 0 0 0 0-13ZM6.5 7.75A.75.75 0 0 1 7.25 7h1a.75.75 0 0 1 .75.75v2.75h.25a.75.75 0 0 1 0 1.5h-2a.75.75 0 0 1 0-1.5h.25v-2h-.25a.75.75 0 0 1-.75-.75ZM8 6a1 1 0 1 1 0-2 1 1 0 0 1 0 2Z"></path></svg>`,
	"tip":       `<svg viewBox="0 0 16 16" width="16" height="16" fill="currentColor"><path d="M8 1.5c-2.363 0-4 1.69-4 3.75 0 .984.424 1.625.984 2.304l.214.253c.223.264.47.556.673.848.284.411.537.856.621 1.345a.75.75 0 0 1-1.478.25c-.04-.238-.19-.53-.44-.892a10.137 10.137 0 0 0-.582-.734l-.2-.238C3.246 7.96 2.5 7.024 2.5 5.25 2.5 2.31 4.863 0 8 0s5.5 2.31 5.5 5.25c0 1.774-.746 2.71-1.292 3.336l-.2.238c-.243.29-.438.53-.582.734-.25.362-.4.655-.44.892a.75.75 0 0 1-1.478-.25c.084-.49.337-.934.621-1.345.204-.292.45-.584.673-.848l.214-.253c.56-.679.984-1.32.984-2.304 0-2.06-1.637-3.75-4-3.75ZM5.75 12h4.5a.75.75 0 0 1 0 1.5h-4.5a.75.75 0 0 1 0-1.5ZM6 15.25a.75.75 0 0 1 .75-.75h2.5a.75.75 0 0 1 0 1.5h-2.5a.75.75 0 0 1-.75-.75Z"></path></svg>`,
	"important": `<svg viewBox="0 0 16 16" width="16" height="16" fill="currentColor"><path d="M0 1.75C0 .784.784 0 1.75 0h12.5C15.216 0 16 .784 16 1.75v9.5A1.75 1.75 0 0 1 14.25 13H8.06l-2.573 2.573A.25.25 0 0 1 5 15.396V13H1.75A1.75 1.75 0 0 1 0 11.25Zm1.75-.25a.25.25 0 0 0-.25.25v9.5c0 .138.112.25.25.25h4a.75.75 0 0 1 .75.75v1.918l1.72-1.72a.749.749 0 0 1 .53-.22h6.5a.25.25 0 0 0 .25-.25v-9.5a.25.25 0 0 0-.25-.25Zm7 2.25v2.5a.75.75 0 0 1-1.5 0v-2.5a.75.75 0 0 1 1.5 0ZM9 9a1 1 0 1 1-2 0 1 1 0 0 1 2 0Z"></path></svg>`,
	"warning":   `<svg viewBox="0 0 16 16" width="16" height="16" fill="currentColor"><path d="M6.457 1.047c.659-1.234 2.427-1.234 3.086 0l6.082 11.378A1.75 1.75 0 0 1 14.082 15H1.918a1.75 1.75 0 0 1-1.543-2.575Zm1.763.707a.25.25 0 0 0-.44 0L1.698 13.132a.25.25 0 0 0 .22.368h12.164a.25.25 0 0 0 .22-.368Zm.53 3.996v2.5a.75.75 0 0 1-1.5 0v-2.5a.75.75 0 0 1 1.5 0ZM9 11a1 1 0 1 1-2 0 1 1 0 0 1 2 0Z"></path></svg>`,
	"caution":   `<svg viewBox="0 0 16 16" width="16" height="16" fill="currentColor"><path d="M4.47.22A.749.749 0 0 1 5 0h6c.199 0 .389.079.53.22l4.25 4.25c.141.141.22.331.22.53v6a.749.749 0 0 1-.22.53l-4.25 4.25A.749.749 0 0 1 11 16H5a.749.749 0 0 1-.53-.22L.22 11.53A.749.749 0 0 1 0 11V5c0-.199.079-.389.22-.53Zm.84 1.28L1.5 5.31v5.38l3.81 3.81h5.38l3.81-3.81V5.31L10.69 1.5ZM8 4a.75.75 0 0 1 .75.75v3.5a.75.75 0 0 1-1.5 0v-3.5A.75.75 0 0 1 8 4Zm0 8a1 1 0 1 1 0-2 1 1 0 0 1 0 2Z"></path></svg>`,
}

// markGitHubAlerts finds blockquotes containing a GitHub alert marker
// (e.g. "[!NOTE]") as the first line of a paragraph, strips the marker from
// the rendered text, and returns the set of blockquote nodes that should be
// rendered as alert boxes, keyed by alert type.
//
// gomarkdown does not end a blockquote at a blank line that is followed by
// another quoted line, so several "> [!TYPE]" alerts written back-to-back
// (blank line between them, no blank line's worth of unquoted content) parse
// as a single ast.BlockQuote containing one paragraph per alert rather than
// several blockquotes. To render each as its own alert box, a blockquote's
// children are split into segments at each marker paragraph, and each
// segment becomes its own *ast.BlockQuote in the tree.
func markGitHubAlerts(doc ast.Node) map[ast.Node]string {
	alerts := map[ast.Node]string{}
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if bq, ok := node.(*ast.BlockQuote); ok && entering {
			splitAlertBlockQuote(bq, alerts)
		}
		return ast.GoToNext
	})
	return alerts
}

// alertSegment is a run of a blockquote's children that belongs to a single
// rendered box: either a plain quote (alertType == "") or one alert.
type alertSegment struct {
	alertType string
	nodes     []ast.Node
}

// splitAlertBlockQuote breaks bq's children into alertSegments at each alert
// marker paragraph and, if at least one marker was found, replaces bq in its
// parent's children with one *ast.BlockQuote per segment, recording each
// alert segment's node and type in alerts.
func splitAlertBlockQuote(bq *ast.BlockQuote, alerts map[ast.Node]string) {
	var segments []alertSegment
	cur := alertSegment{}
	flush := func() {
		if len(cur.nodes) > 0 || cur.alertType != "" {
			segments = append(segments, cur)
		}
	}
	for _, child := range bq.GetChildren() {
		if alertType, kept, ok := extractAlertMarker(child); ok {
			flush()
			cur = alertSegment{alertType: alertType}
			if kept != nil {
				cur.nodes = append(cur.nodes, kept)
			}
			continue
		}
		cur.nodes = append(cur.nodes, child)
	}
	flush()

	hasAlert := false
	for _, seg := range segments {
		if seg.alertType != "" {
			hasAlert = true
			break
		}
	}
	if !hasAlert {
		return
	}

	parent := bq.GetParent()
	if parent == nil {
		return
	}
	siblings := parent.GetChildren()
	idx := -1
	for i, n := range siblings {
		if n == ast.Node(bq) {
			idx = i
			break
		}
	}
	if idx < 0 {
		return
	}

	replacement := make([]ast.Node, 0, len(segments))
	for _, seg := range segments {
		newBQ := &ast.BlockQuote{}
		newBQ.SetParent(parent)
		for _, n := range seg.nodes {
			n.SetParent(newBQ)
		}
		newBQ.SetChildren(seg.nodes)
		if seg.alertType != "" {
			alerts[ast.Node(newBQ)] = seg.alertType
		}
		replacement = append(replacement, ast.Node(newBQ))
	}

	merged := make([]ast.Node, 0, len(siblings)-1+len(replacement))
	merged = append(merged, siblings[:idx]...)
	merged = append(merged, replacement...)
	merged = append(merged, siblings[idx+1:]...)
	parent.SetChildren(merged)
}

// extractAlertMarker reports whether node is a paragraph whose first line is
// a GitHub alert marker. When it is, the marker is stripped from the
// paragraph's text and the paragraph is returned as kept, or kept is nil if
// the marker was the paragraph's entire content.
func extractAlertMarker(node ast.Node) (alertType string, kept ast.Node, ok bool) {
	para, isPara := node.(*ast.Paragraph)
	if !isPara {
		return "", nil, false
	}
	paraChildren := para.GetChildren()
	if len(paraChildren) == 0 {
		return "", nil, false
	}
	text, isText := paraChildren[0].(*ast.Text)
	if !isText {
		return "", nil, false
	}

	literal := string(text.Literal)
	firstLine, rest, hasRest := strings.Cut(literal, "\n")
	if !hasRest {
		firstLine = literal
	}
	alertType, ok = githubAlertMarkers[strings.TrimSpace(firstLine)]
	if !ok {
		return "", nil, false
	}

	if !hasRest && len(paraChildren) == 1 {
		return alertType, nil, true
	}
	text.Literal = []byte(rest)
	return alertType, para, true
}

// githubAlertHook returns a RenderNodeFunc that renders blockquotes marked
// by markGitHubAlerts as "<div class=\"markdown-alert markdown-alert-TYPE\">"
// boxes with an icon and title, matching GitHub's alert rendering.
func githubAlertHook(alerts map[ast.Node]string) html.RenderNodeFunc {
	return func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		alertType, ok := alerts[node]
		if !ok {
			return ast.GoToNext, false
		}

		if entering {
			io.WriteString(w, `<div class="markdown-alert markdown-alert-`+alertType+`">`+"\n")
			io.WriteString(w, `<p class="markdown-alert-title">`+githubAlertIcons[alertType]+" "+githubAlertTitles[alertType]+"</p>\n")
		} else {
			io.WriteString(w, "</div>\n")
		}

		return ast.GoToNext, true
	}
}

// assignGitHubHeadingIDs sets HeadingID on every heading using GitHub's
// anchor-slug algorithm, so links from GitHub-flavored TOC tools resolve
// the same way locally as they do on github.com.
func assignGitHubHeadingIDs(doc ast.Node) {
	taken := map[string]bool{}
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		heading, ok := node.(*ast.Heading)
		if !ok || !entering || heading.HeadingID != "" {
			return ast.GoToNext
		}

		id := githubHeadingSlug(headingText(heading))
		base := id
		for n := 1; taken[id]; n++ {
			id = base + "-" + strconv.Itoa(n)
		}
		taken[id] = true
		heading.HeadingID = id

		return ast.GoToNext
	})
}

// headingText concatenates the literal text of a heading's descendants.
func headingText(node ast.Node) string {
	var sb strings.Builder
	ast.WalkFunc(node, func(n ast.Node, entering bool) ast.WalkStatus {
		if leaf, ok := n.(*ast.Text); ok && entering {
			sb.Write(leaf.Literal)
		}
		return ast.GoToNext
	})
	return sb.String()
}

// githubHeadingSlug reproduces GitHub's heading-anchor algorithm: lowercase,
// turn spaces into hyphens, strip any character that isn't alphanumeric,
// '-', or '_', then collapse repeated hyphens and trim leading/trailing ones.
func githubHeadingSlug(text string) string {
	var sb strings.Builder
	for _, r := range text {
		switch {
		case r == ' ':
			sb.WriteRune('-')
		case unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' || r == '_':
			sb.WriteRune(unicode.ToLower(r))
		}
	}

	slug := sb.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	return strings.Trim(slug, "-")
}
