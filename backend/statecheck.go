package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"
)

// StateChecker handles POST /api/state/{demo} by running goss validate
// against the appropriate test file and returning the JSON output.
type StateChecker struct{}

func (s *StateChecker) Handle(w http.ResponseWriter, r *http.Request) {
	demo := r.PathValue("demo")

	// Validate demo name
	switch demo {
	case "rbac", "policy", "netpol":
	default:
		http.Error(w, fmt.Sprintf("unknown demo: %s", demo), http.StatusBadRequest)
		return
	}

	gossFile := fmt.Sprintf("/app/goss/%s.yaml", demo)

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "goss",
		"--gossfile", gossFile,
		"validate",
		"--format", "json",
	)

	output, err := cmd.CombinedOutput()

	// Validate that output is valid JSON before returning it.
	// Goss exits non-zero when tests fail (expected), but also when it
	// can't run at all (gossfile missing, binary not found, etc.).
	// In either case we need parseable JSON to send to the frontend.
	if !json.Valid(output) {
		msg := string(output)
		if err != nil {
			msg = fmt.Sprintf("%s: %v", msg, err)
		}
		log.Printf("goss error for %s: %s", demo, msg)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		errResp, _ := json.Marshal(map[string]string{
			"error": fmt.Sprintf("state check failed: %s", msg),
		})
		w.Write(errResp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}
