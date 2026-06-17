package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type broadcaster struct {
	mu      sync.Mutex
	clients map[chan struct{}]struct{}
}

func newBroadcaster() *broadcaster {
	return &broadcaster{clients: make(map[chan struct{}]struct{})}
}

func (b *broadcaster) subscribe() chan struct{} {
	ch := make(chan struct{}, 1)
	b.mu.Lock()
	b.clients[ch] = struct{}{}
	b.mu.Unlock()
	return ch
}

func (b *broadcaster) unsubscribe(ch chan struct{}) {
	b.mu.Lock()
	delete(b.clients, ch)
	b.mu.Unlock()
}

func (b *broadcaster) broadcast() {
	b.mu.Lock()
	n := len(b.clients)
	for ch := range b.clients {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
	b.mu.Unlock()
	log.Printf("watch: broadcasting reload to %d client(s)", n)
}

func startWatcher(cfg *Config, b *broadcaster) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// Watch the root dir and all existing subdirectories.
	err = filepath.WalkDir(cfg.Dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		watcher.Close()
		return err
	}

	go func() {
		defer watcher.Close()
		var debounce *time.Timer
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// Watch newly created subdirectories.
				if event.Has(fsnotify.Create) {
					if fi, err := os.Stat(event.Name); err == nil && fi.IsDir() {
						watcher.Add(event.Name)
					}
				}
				if strings.ToLower(filepath.Ext(event.Name)) != ".md" {
					continue
				}
				log.Printf("watch: %s %s", event.Op, event.Name)
				// Debounce: editors often emit multiple events on a single save.
				if debounce != nil {
					debounce.Stop()
				}
				debounce = time.AfterFunc(100*time.Millisecond, b.broadcast)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("watch error: %v", err)
			}
		}
	}()

	return nil
}

func sseHandler(b *broadcaster) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming not supported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Send initial comment to flush headers to the browser immediately.
		fmt.Fprintf(w, ": connected\n\n")
		flusher.Flush()

		ch := b.subscribe()
		defer b.unsubscribe(ch)

		// Keepalive ticker: prevents proxies and browsers from closing idle connections.
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ch:
				fmt.Fprintf(w, "data: reload\n\n")
				flusher.Flush()
			case <-ticker.C:
				fmt.Fprintf(w, ": ping\n\n")
				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}
	}
}
