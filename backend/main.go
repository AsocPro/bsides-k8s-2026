package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	hub := NewEventHub()

	// Start heartbeat
	go hub.StartHeartbeat(ctx, 10*time.Second)

	// Start JSONL watcher if configured
	if jsonlPath := os.Getenv("JSONL_WATCH_PATH"); jsonlPath != "" {
		watcher := NewJSONLWatcher(hub, jsonlPath)
		go watcher.Watch(ctx)
	}

	mux := http.NewServeMux()

	// Event bridge WebSocket
	mux.Handle("/ws/events", hub)

	// REST API for pushing events programmatically
	mux.HandleFunc("POST /api/events", func(w http.ResponseWriter, r *http.Request) {
		var evt Event
		if err := decodeJSON(r.Body, &evt); err != nil {
			http.Error(w, "invalid event", http.StatusBadRequest)
			return
		}
		evt.Source = "api"
		hub.Broadcast(evt)
		w.WriteHeader(http.StatusAccepted)
	})

	// Environment API stubs
	mux.HandleFunc("/api/environments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"environments":[]}`))
	})

	// Terminal proxy: reverse proxy to ttyd instances
	if ttydURLs := os.Getenv("TTYD_URLS"); ttydURLs != "" {
		for _, entry := range strings.Split(ttydURLs, ",") {
			parts := strings.SplitN(entry, "=", 2)
			if len(parts) != 2 {
				continue
			}
			name, rawURL := parts[0], parts[1]
			target, err := url.Parse(rawURL)
			if err != nil {
				log.Printf("invalid ttyd URL for %s: %v", name, err)
				continue
			}

			prefix := "/terminal/" + name + "/"
			proxy := &httputil.ReverseProxy{
				Director: func(req *http.Request) {
					req.URL.Scheme = target.Scheme
					req.URL.Host = target.Host
					req.Host = target.Host
				},
			}
			mux.Handle(prefix, proxy)
			log.Printf("Terminal proxy: %s -> %s", prefix, rawURL)
		}
	}

	// Serve frontend: reverse proxy to Vite in dev, or static files in prod
	if viteURL := os.Getenv("VITE_DEV_URL"); viteURL != "" {
		target, err := url.Parse(viteURL)
		if err != nil {
			log.Fatalf("invalid VITE_DEV_URL: %v", err)
		}
		proxy := httputil.NewSingleHostReverseProxy(target)
		mux.Handle("/", proxy)
		log.Printf("Proxying frontend to %s", viteURL)
	} else {
		staticDir := os.Getenv("STATIC_DIR")
		if staticDir == "" {
			staticDir = "./static"
		}
		mux.Handle("/", http.FileServer(http.Dir(staticDir)))
		log.Printf("Serving static files from %s", staticDir)
	}

	addr := ":8080"
	log.Printf("Backend listening on %s", addr)

	server := &http.Server{Addr: addr, Handler: mux}
	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
