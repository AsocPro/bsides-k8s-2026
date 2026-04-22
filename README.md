# Kubernetes: Orchestrating a More Secure Future

A 25-minute BSides talk delivered as a **live, interactive web presentation**
backed by real k3s clusters. The slides themselves run the demos — every
RBAC rejection, every Kyverno deny, every network-policy block shown in the
browser is the actual output of `kubectl` against a freshly-spun k3s cluster.

Everything — frontend, backend, k3s clusters, PDF export — runs inside
[Dagger](https://dagger.io). The only thing you need on your machine is
Dagger itself and a container runtime.

See also: [OVERVIEW.md](OVERVIEW.md), [SETUP.md](SETUP.md),
[SPEAKER_NOTES.md](SPEAKER_NOTES.md), [PRESENTATION_PLAN.md](PRESENTATION_PLAN.md).

---

## Prerequisites

- [Dagger](https://docs.dagger.io/install) v0.20.1 or newer
- Docker **or** Podman as the container runtime
- ~4 GB RAM free (three small k3s clusters run concurrently)
- A modern browser for viewing the presentation at `http://localhost:8080`

No Node, Go, kubectl, or k3s installs required — Dagger pulls everything.

### Podman users: enable cgroup delegation first

k3s needs the `cpuset` cgroup controller, and Podman's default user slice
doesn't delegate it. Without this fix, the run fails with:

```
level=fatal msg="failed to find cpuset cgroup (v2)"
```

Run the included one-shot script, then restart the Dagger engine so it
picks up the new delegation:

```bash
sudo ./scripts/setup-cgroup-delegation.sh

podman rm -f $(podman ps -q --filter name=dagger-engine)
dagger version   # triggers engine restart
```

Verify:

```bash
podman exec $(podman ps -q --filter name=dagger-engine) \
  cat /sys/fs/cgroup/cgroup.controllers
# expect: cpuset cpu io memory hugetlb pids
```

Docker users can skip this section — cgroup delegation is handled for you.

---

## Run the presentation

```bash
dagger call present up --ports 8080:8080
```

Open `http://localhost:8080`. Navigate with **→ / Space** (next) and
**← / Shift+Space** (previous). Some slides (Evolution, K8s Intro) have
progressive substeps — pressing next reveals content before advancing.

The first run is slow: Dagger has to pull the k3s image, build the
frontend/backend, install Kyverno, etc. Subsequent runs reuse the cache
volumes and come up in ~30 seconds.

---

## Export the slides as a PDF

Everything runs inside Dagger — no local Playwright or Chromium install.
The function builds the frontend, serves it with nginx in a container,
drives Chromium via Playwright to step through each slide (including
substep animations), stitches the screenshots into a PDF, and prepends a
cover page with a clickable link to the repo.

```bash
dagger call export-pdf export --path ./bsides-k8s-2026.pdf
```

Override the cover-page link with `--repo-url https://…` if you fork.
Implementation: `.dagger/main.go` (function `ExportPdf`) and
`.dagger/scripts/capture.js`.

---

## Demo walkthrough

Three live demos are wired to real clusters. Every command below is what
the speaker types in the embedded terminal during the talk; the slide's
left panel updates reactively via goss-driven state checks.

### RBAC (slide 6) — cluster profile `base`

Pre-seeded: namespace `demo`, ServiceAccount `demo-user` (no perms), pod
`web` running nginx.

```bash
# 1. No permissions → denied
kubectl --as=system:serviceaccount:demo:demo-user -n demo get pods

# 2. Create a Role
kubectl apply -f - <<'EOF'
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata: { name: pod-reader, namespace: demo }
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list"]
EOF

# 3. Bind it to the ServiceAccount
kubectl apply -f - <<'EOF'
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata: { name: read-pods, namespace: demo }
roleRef: { apiGroup: rbac.authorization.k8s.io, kind: Role, name: pod-reader }
subjects:
- { kind: ServiceAccount, name: demo-user, namespace: demo }
EOF

# 4. Retry → succeeds
kubectl --as=system:serviceaccount:demo:demo-user -n demo get pods

# 5. Try to delete → still denied (least privilege)
kubectl --as=system:serviceaccount:demo:demo-user -n demo delete pod web
```

### Policy (slide 8) — cluster profile `kyverno`

Pre-seeded: Kyverno + two `ClusterPolicy` objects (`require-non-root`,
`require-approved-registries`) + namespace `demo`.

```bash
kubectl get clusterpolicy

# Root container → rejected by require-non-root
kubectl run bad-pod --image=nginx:alpine -n demo

# Compliant pod → accepted
kubectl apply -f - <<'EOF'
apiVersion: v1
kind: Pod
metadata: { name: good-pod, namespace: demo }
spec:
  containers:
  - name: web
    image: nginxinc/nginx-unprivileged:alpine
    securityContext: { runAsNonRoot: true, allowPrivilegeEscalation: false }
EOF

# Unapproved registry → rejected
kubectl apply -f - <<'EOF'
apiVersion: v1
kind: Pod
metadata: { name: sketchy, namespace: demo }
spec:
  containers:
  - name: app
    image: evil-registry.io/malware:latest
    securityContext: { runAsNonRoot: true, allowPrivilegeEscalation: false }
EOF
```

### Network Policies (slide 10) — cluster profile `netpol`

Pre-seeded: 3-tier app (`frontend`, `api`, `database`) with services.

```bash
# Everything talks to everything by default
kubectl -n demo exec deployment/frontend -- curl -s --max-time 3 api
kubectl -n demo exec deployment/frontend -- curl -s --max-time 3 database
kubectl -n demo exec deployment/api      -- curl -s --max-time 3 database

# Default-deny locks everything down
kubectl apply -f /root/manifests/deny-all.yaml
kubectl -n demo exec deployment/frontend -- curl -s --max-time 3 api   # times out

# Restore only allowed paths
kubectl apply -f /root/manifests/allow-frontend-to-api.yaml
kubectl apply -f /root/manifests/allow-api-to-db.yaml

# Frontend → database stays blocked (blast radius contained)
kubectl -n demo exec deployment/frontend -- curl -s --max-time 3 database
```

Full speaker script with timing is in [SPEAKER_NOTES.md](SPEAKER_NOTES.md).

---

## Architecture

```
dagger call present
├── BuildFrontend         Svelte 5 → /app/static
├── BuildBackend          Go → /app/server
├── K3sCluster "rbac"     base     + manifests/rbac/
├── K3sCluster "policy"   kyverno  + manifests/policy/
├── K3sCluster "netpol"   netpol   + manifests/netpol/
└── Backend container (+ kubectl, goss, kubeconfigs)
    ├── /                         static Svelte app
    ├── /terminal/{name}/         reverse proxy → ttyd in the k3s container
    └── POST /api/state/{demo}    runs goss checks → JSON → reactive UI
```

Each k3s cluster runs inside a Dagger container. ttyd runs alongside k3s
in the same container and exposes a shell with kubectl pre-wired to the
cluster's kubeconfig. The Go backend reverse-proxies `/terminal/{name}/`
to ttyd over the Dagger service network. Button clicks like "Check State"
trigger `POST /api/state/{demo}` → `goss validate` over a goss YAML in
`goss/` → JSON results → reactive state panel.

### Cluster profiles

| Profile   | What it adds                                | Used by      |
|-----------|---------------------------------------------|--------------|
| `base`    | Plain k3s                                   | RBAC demo    |
| `kyverno` | k3s + Kyverno admission controller          | Policy demo  |
| `netpol`  | k3s with a CNI that enforces NetworkPolicy  | Netpol demo  |

Profiles are defined in `.dagger/k3s.go`.

---

## Repository layout

```
.dagger/              Dagger module (Go SDK)
  main.go             Top-level functions (present, export-pdf, …)
  k3s.go              Cluster profile builder + seed script
  scripts/            Assets embedded into Dagger (PDF capture, nginx conf)
backend/              Go HTTP server + state-check API
frontend/             Svelte 5 slide engine
  src/slides/         One component per slide
  src/lib/            SlideEngine, Terminal, RBACMatrix, NetPolGraph, etc.
manifests/
  rbac/  policy/  netpol/    per-profile bootstrap manifests
goss/                 goss YAML files backing the /api/state/{demo} endpoints
scripts/              Host-side helpers (cgroup delegation for Podman)
```

---

## Troubleshooting

- **`failed to find cpuset cgroup (v2)`** — Podman without cgroup
  delegation. Run `sudo ./scripts/setup-cgroup-delegation.sh` and restart
  the Dagger engine.
- **First run hangs for a while** — Dagger is pulling `rancher/k3s` and
  installing Kyverno. Give it a few minutes; subsequent runs hit the cache.
- **Terminals show "connection refused"** — the k3s cluster for that
  profile is still starting. Wait for the state panel to go green, or
  refresh the page.
- **PDF export fails at `npm install`** — Playwright image version
  mismatch. The image tag in `.dagger/main.go` (`playwrightVersion`) must
  match the dependency pinned in the same file; bump both together.

---

## License

[MIT](LICENSE).
