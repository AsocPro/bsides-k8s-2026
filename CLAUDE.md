# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

BSides K8s — a 25-minute interactive web presentation on Kubernetes security with live k3s cluster demos orchestrated by Dagger. Frontend is Svelte 5 + Vite, backend is Go, infrastructure is Dagger (Go SDK) spinning up k3s clusters with ttyd terminals.

## Build & Run Commands

Everything runs via Dagger — no local Node/Go installs needed.

```bash
# Development with Vite HMR
dagger call dev up --ports 8080:8080

# Production build (no k3s clusters)
dagger call serve up --ports 8080:8080

# Full presentation with live k3s clusters + terminals
dagger call present up --ports 8080:8080

# Debug a specific k3s cluster (profiles: base, kyverno, netpol)
dagger call k3s-debug --profile base terminal
```

There are no separate lint or test commands — the project is presentation-focused.

## Architecture

```
Frontend (Svelte 5, frontend/)
  └─ Slide engine with GSAP animations, embedded ttyd terminals, WebSocket event overlays

Backend (Go 1.24, backend/)
  └─ HTTP server: serves static frontend (or proxies Vite in dev), reverse-proxies ttyd
     terminals, runs WebSocket event hub (/ws/events), tails JSONL command logs

Dagger (.dagger/)
  └─ Orchestrates everything: builds frontend+backend, spins up k3s clusters per profile,
     attaches ttyd terminals, pre-seeds manifests, wires services together
```

**Event flow:** Terminal commands → JSONL file → Go watcher → WebSocket hub → Svelte stores → reactive overlays

**K3s cluster profiles** (defined in `.dagger/k3s.go`):
- `base` — plain k3s for RBAC demos
- `kyverno` — k3s + Kyverno admission controller for policy demos
- `netpol` — k3s with network policy controller for network segmentation demos

Each cluster shares kubeconfig with its ttyd terminal via Dagger cache volumes.

## Key Patterns

- **Dagger service binding:** ttyd terminals and k3s clusters are Dagger services bound together. Backend reads `TTYD_URLS` env var for reverse proxy routes.
- **Dagger AsService:** Always use `AsService(Args)` — never `WithDefaultArgs` + `AsService` or `WithExec` + `AsService`.
- **Frontend dev proxy:** In dev mode, backend reverse-proxies `/` to `http://vite:5173` via `VITE_DEV_URL` env var.
- **Manifests pre-seeding:** `manifests/{rbac,policy,netpol}/` are applied to clusters before demos start, via the seed script in k3s.go.

## Podman Note

The user runs Podman, not Docker. k3s requires systemd cgroup delegation — run `sudo ./scripts/setup-cgroup-delegation.sh` if cpuset errors occur, then restart Dagger.
