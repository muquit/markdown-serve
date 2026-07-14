package main

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var imageContentTypes = map[string]string{
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".svg":  "image/svg+xml",
	".webp": "image/webp",
	".ico":  "image/x-icon",
}

func makeHandler(cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			serveIndex(w, r, cfg)
			return
		}
		serveMarkdown(w, r, cfg)
	}
}

type TreeNode struct {
	Name     string
	Path     string // slash-separated relative path for URL (files only)
	IsDir    bool
	Children []*TreeNode
}

func buildTree(root string) (*TreeNode, int, error) {
	rootNode := &TreeNode{Name: filepath.Base(root), IsDir: true}
	nodeMap := map[string]*TreeNode{".": rootNode}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil || rel == "." {
			return nil
		}
		rel = filepath.ToSlash(rel)
		parentKey := filepath.ToSlash(filepath.Dir(rel))
		parent, ok := nodeMap[parentKey]
		if !ok {
			return nil
		}
		if d.IsDir() {
			node := &TreeNode{Name: d.Name(), IsDir: true}
			nodeMap[rel] = node
			parent.Children = append(parent.Children, node)
			return nil
		}
		if strings.ToLower(filepath.Ext(d.Name())) == ".md" {
			parent.Children = append(parent.Children, &TreeNode{
				Name: d.Name(),
				Path: rel,
			})
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	count := pruneAndCount(rootNode)
	sortTree(rootNode)
	return rootNode, count, nil
}

// pruneAndCount removes empty directories and returns the total file count.
func pruneAndCount(node *TreeNode) int {
	if !node.IsDir {
		return 1
	}
	total := 0
	var kept []*TreeNode
	for _, child := range node.Children {
		n := pruneAndCount(child)
		if n > 0 {
			kept = append(kept, child)
			total += n
		}
	}
	node.Children = kept
	return total
}

// sortTree sorts children: directories first, then files, each alphabetically.
func sortTree(node *TreeNode) {
	sort.Slice(node.Children, func(i, j int) bool {
		a, b := node.Children[i], node.Children[j]
		if a.IsDir != b.IsDir {
			return a.IsDir
		}
		return strings.ToLower(a.Name) < strings.ToLower(b.Name)
	})
	for _, child := range node.Children {
		if child.IsDir {
			sortTree(child)
		}
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request, cfg *Config) {
	root, count, err := buildTree(cfg.Dir)
	if err != nil {
		http.Error(w, "cannot walk directory", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	renderIndex(w, cfg.Dir, root, count, cfg.Watch, cfg.Dark)
}

func resolveServedPath(cfg *Config, name string) (string, error) {
	absBase, err := filepath.Abs(cfg.Dir)
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(absBase, filepath.FromSlash(name))

	if !strings.HasPrefix(fullPath, absBase+string(filepath.Separator)) {
		return "", os.ErrPermission
	}

	return fullPath, nil
}

func serveMarkdown(w http.ResponseWriter, r *http.Request, cfg *Config) {
	name := strings.TrimPrefix(r.URL.Path, "/")

	if strings.Contains(name, "..") {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	ext := strings.ToLower(filepath.Ext(name))
	if ext != ".md" {
		if contentType, ok := imageContentTypes[ext]; ok {
			serveImage(w, cfg, name, contentType)
			return
		}
		http.Error(w, "only .md files are served", http.StatusNotFound)
		return
	}

	fullPath, err := resolveServedPath(cfg, name)
	if err != nil {
		if err == os.ErrPermission {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	fi, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "file not found: "+name, http.StatusNotFound)
			return
		}
		http.Error(w, "cannot read file", http.StatusInternalServerError)
		return
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "cannot read file", http.StatusInternalServerError)
		return
	}

	html := renderMarkdown(data)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	renderDocument(w, name, html, fi.ModTime().Format("2006-01-02 15:04:05"), cfg.Watch, cfg.Version, cfg.Dark)
}

func serveImage(w http.ResponseWriter, cfg *Config, name, contentType string) {
	fullPath, err := resolveServedPath(cfg, name)
	if err != nil {
		if err == os.ErrPermission {
			http.Error(w, "invalid path", http.StatusBadRequest)
			return
		}
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "file not found: "+name, http.StatusNotFound)
			return
		}
		http.Error(w, "cannot read file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	_, _ = w.Write(data)
}
