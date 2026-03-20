# BSides K8s Presentation — Setup Guide

## Prerequisites

- [Dagger](https://docs.dagger.io/install) v0.20.1+
- Podman **or** Docker as the container runtime
- ~4GB RAM available for k3s clusters

## Quick Start

```bash
# Dev mode (no k3s, fast iteration on slides)
dagger call dev up --ports 8080:8080

# Serve mode (production build, no k3s)
dagger call serve up --ports 8080:8080

# Present mode (full presentation with live k3s clusters)
dagger call present up --ports 8080:8080
```

## Podman-Specific Setup

If you use Podman as your container runtime, you **must** enable systemd cgroup
delegation before running k3s clusters. Without this, k3s fails with:

```
level=fatal msg="failed to find cpuset cgroup (v2)"
```

### Why

Podman runs under `user.slice` which only delegates `cpu memory pids` cgroup
controllers by default. k3s requires `cpuset`, which isn't delegated. Docker
handles this automatically; Podman does not.

### Fix

Run the included setup script (requires sudo):

```bash
sudo ./scripts/setup-cgroup-delegation.sh
```

This creates `/etc/systemd/system/user@.service.d/delegate.conf` with:

```ini
[Service]
Delegate=yes
```

Then restart the Dagger engine:

```bash
podman rm -f $(podman ps -q --filter name=dagger-engine)
dagger version  # triggers engine restart
```

### Verify

```bash
podman exec $(podman ps -q --filter name=dagger-engine) \
  cat /sys/fs/cgroup/cgroup.controllers
# Should include: cpuset cpu io memory pids
```

## Architecture

Everything runs via `dagger call` — no local installs beyond Dagger itself.

```
dagger call present
├── BuildFrontend   → Svelte 5 + Vite + Tailwind v4
├── BuildBackend    → Go HTTP server (proxy, state-check API)
├── K3sCluster "rbac"    → k3s (base profile) + RBAC manifests
├── K3sCluster "policy"  → k3s (kyverno profile) + policy manifests
├── K3sCluster "netpol"  → k3s (netpol profile) + network policy manifests
└── Go server container (+ kubectl, goss, kubeconfigs)
    ├── /                    → slides (static or Vite proxy)
    ├── /terminal/*          → reverse proxy to ttyd containers
    └── /api/state/{demo}   → POST: run goss checks, return JSON
```

### Cluster Profiles

| Profile  | What it adds                          | Demo section       |
|----------|---------------------------------------|---------------------|
| base     | Plain k3s                             | RBAC                |
| kyverno  | k3s + Kyverno policy engine           | Policy Agents       |
| netpol   | k3s (built-in network policy support) | Network Policies    |

### K3s in Dagger

The k3s clusters use the same pattern as [Dagger's own helm-dev toolchain](https://github.com/dagger/dagger/tree/main/toolchains/helm-dev/k3s):

- Cgroup v2 entrypoint script (sourced from Docker's DinD / Moby)
- `InsecureRootCapabilities: true` for k3s server service
- Cache volume at `/etc/rancher/k3s` shares kubeconfig between containers
- `Config()` polls the cache volume until `k3s.yaml` appears
- Demo containers use `WithServiceBinding` + `WithFile` for kubeconfig

### Dagger Functions

| Function         | Description                                           |
|------------------|-------------------------------------------------------|
| `build-frontend` | Build Svelte app → dist directory                     |
| `build-backend`  | Compile Go server → static binary                    |
| `build`          | Both combined into one directory                      |
| `terminal`       | Basic ttyd container (no k3s)                         |
| `k-8-s-terminal` | ttyd + kubectl connected to a k3s cluster             |
| `dev`            | Vite HMR + Go backend + basic terminal                |
| `serve`          | Production build + basic terminal                     |
| `present`        | Production build + 3 k3s clusters for live demos      |
