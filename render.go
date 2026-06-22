package main

import (
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

	flags := html.UseXHTML | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: flags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
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
