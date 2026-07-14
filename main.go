package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var version = "dev" // overridden at build time via -ldflags

type Config struct {
	Host    string
	Port    int
	Dir     string
	Watch   bool
	Dark    bool
	Version string
}

func main() {
	host := flag.String("host", "0.0.0.0", "Host to bind to")
	port := flag.Int("port", 8485, "Port to listen on")
	watch := flag.Bool("watch", true, "Reload browser on file changes")
	dark := flag.Bool("dark", false, "Render pages in dark mode")
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	dir := "."
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}

	info, err := os.Stat(dir)
	if err != nil {
		log.Fatalf("cannot access directory %q: %v", dir, err)
	}
	if !info.IsDir() {
		log.Fatalf("%q is not a directory", dir)
	}

	cfg := &Config{
		Host:    *host,
		Port:    *port,
		Dir:     dir,
		Watch:   *watch,
		Dark:    *dark,
		Version: version,
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Printf("Serving %s on http://%s", cfg.Dir, addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", makeHandler(cfg))

	if cfg.Watch {
		b := newBroadcaster()
		if err := startWatcher(cfg, b); err != nil {
			log.Fatalf("watcher error: %v", err)
		}
		mux.HandleFunc("/events", sseHandler(b))
		log.Printf("Watch mode enabled")
	}

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
