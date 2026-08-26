package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
)

type broadcaster struct {
	revision atomic.Uint64
}

func newBroadcaster() *broadcaster {
	return &broadcaster{}
}

func (b *broadcaster) broadcast() {
	revision := b.revision.Add(1)
	log.Printf("watch: revision %d", revision)
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

func watchHandler(b *broadcaster) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = fmt.Fprintln(w, strconv.FormatUint(b.revision.Load(), 10))
	}
}
