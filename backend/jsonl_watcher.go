package main

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
)

// JSONLWatcher tails a JSONL file and emits events to the hub.
type JSONLWatcher struct {
	hub  *EventHub
	path string
}

func NewJSONLWatcher(hub *EventHub, path string) *JSONLWatcher {
	return &JSONLWatcher{hub: hub, path: path}
}

// Watch starts tailing the JSONL file, emitting events for each new line.
// It starts reading from the end of the file (only new entries).
func (w *JSONLWatcher) Watch(ctx context.Context) {
	for {
		err := w.tailFile(ctx)
		if ctx.Err() != nil {
			return
		}
		if err != nil {
			log.Printf("jsonl watcher error: %v, retrying in 2s", err)
		}
		select {
		case <-time.After(2 * time.Second):
		case <-ctx.Done():
			return
		}
	}
}

func (w *JSONLWatcher) tailFile(ctx context.Context) error {
	f, err := os.Open(w.path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Seek to end — only read new lines
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	log.Printf("jsonl watcher: tailing %s", w.path)

	reader := bufio.NewReader(f)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				// No new data, poll
				select {
				case <-time.After(100 * time.Millisecond):
					continue
				case <-ctx.Done():
					return nil
				}
			}
			return err
		}

		if len(line) == 0 {
			continue
		}

		w.processLine(line)
	}
}

func (w *JSONLWatcher) processLine(line []byte) {
	// Try to parse as a structured event
	var raw map[string]any
	if err := json.Unmarshal(line, &raw); err != nil {
		// Not valid JSON, emit as raw command
		w.hub.Broadcast(Event{
			Type:   "command",
			Source: "jsonl",
			Data:   map[string]any{"raw": string(line)},
		})
		return
	}

	// If the JSONL entry has a "type" field, use it
	eventType, _ := raw["type"].(string)
	if eventType == "" {
		eventType = "command"
	}

	// If it has a "command" field, extract it for the command echo overlay
	cmd, _ := raw["command"].(string)
	exitCode, _ := raw["exit_code"].(float64)

	w.hub.Broadcast(Event{
		Type:   eventType,
		Source: "jsonl",
		Data: map[string]any{
			"command":   cmd,
			"exit_code": int(exitCode),
			"raw":       raw,
		},
	})
}
