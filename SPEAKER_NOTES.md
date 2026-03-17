# Speaker Notes — Kubernetes: Orchestrating a More Secure Future

## Starting the Presentation

```bash
dagger call present up --ports 8080:8080
```

Open http://localhost:8080 — navigate with arrow keys or spacebar.

---

## Slide 1: Title (no commands)

Just the title card. Take a breath, greet the audience.

---

## Slide 2: Evolution — "How We Got Here" (~4 min, no commands)

Animated cards walk through Bare Metal → VMs → Containers. No terminal needed.

---

## Slide 3: K8s Intro — "Kubernetes in 60 Seconds" (~3 min, no commands)

Core concepts and built-in security features. No terminal needed.

---

## Slide 4: RBAC Concepts (no commands)

Explain Roles, RoleBindings, ServiceAccounts, and Least Privilege before the demo.

---

## Slide 5: RBAC Demo (~5 min)

Terminal environment: `rbac` (base k3s cluster). Pre-deployed: namespace `demo`, ServiceAccount `demo-user` (no permissions), pod `web` running nginx.

### Step 1 — Show that demo-user has no permissions

```bash
kubectl --as=system:serviceaccount:demo:demo-user -n demo get pods
```

Expected: `Error from server (Forbidden)`

### Step 2 — Create a Role allowing pod read access

```bash
kubectl apply -f - <<EOF
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: pod-reader
  namespace: demo
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list"]
EOF
```

### Step 3 — Bind the Role to demo-user

```bash
kubectl apply -f - <<EOF
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: read-pods
  namespace: demo
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: pod-reader
subjects:
- kind: ServiceAccount
  name: demo-user
  namespace: demo
EOF
```

### Step 4 — Retry listing pods (should now succeed)

```bash
kubectl --as=system:serviceaccount:demo:demo-user -n demo get pods
```

Expected: see the `web` pod listed.

### Step 5 — Try deleting a pod (should be denied — least privilege)

```bash
kubectl --as=system:serviceaccount:demo:demo-user -n demo delete pod web
```

Expected: `Error from server (Forbidden)` — the role only grants get/list, not delete.

---

## Slide 6: Policy Concepts (no commands)

Explain that RBAC controls *who*, policies control *what*. Introduce Kyverno.

---

## Slide 7: Policy Demo (~5 min)

Terminal environment: `policy` (kyverno k3s cluster). Pre-deployed: Kyverno admission controller, `require-non-root` ClusterPolicy, `require-approved-registries` ClusterPolicy, namespace `demo`.

### Step 1 — Inspect the active policy

```bash
kubectl get clusterpolicy
```

```bash
kubectl get clusterpolicy require-non-root -o yaml
```

### Step 2 — Try deploying a container that runs as root (rejected)

```bash
kubectl run bad-pod --image=nginx:alpine -n demo
```

Expected: rejected by Kyverno — `require-non-root` policy blocks pods without `runAsNonRoot: true`.

### Step 3 — Deploy a properly configured non-root container (accepted)

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: good-pod
  namespace: demo
spec:
  containers:
  - name: web
    image: nginxinc/nginx-unprivileged:alpine
    securityContext:
      runAsNonRoot: true
      allowPrivilegeEscalation: false
EOF
```

Expected: pod created successfully.

```bash
kubectl -n demo get pods
```

### Step 4 — Test the registry policy

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: sketchy
  namespace: demo
spec:
  containers:
  - name: app
    image: evil-registry.io/malware:latest
    securityContext:
      runAsNonRoot: true
      allowPrivilegeEscalation: false
EOF
```

Expected: rejected — `require-approved-registries` only allows docker.io, ghcr.io, quay.io, registry.k8s.io.

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: legit
  namespace: demo
spec:
  containers:
  - name: web
    image: docker.io/nginxinc/nginx-unprivileged:1.27
    securityContext:
      runAsNonRoot: true
      allowPrivilegeEscalation: false
EOF
```

Expected: accepted — approved registry and non-root.

---

## Slide 8: Network Policy Concepts (no commands)

Explain microsegmentation, defense in depth, and blast radius reduction.

---

## Slide 9: Network Policy Demo (~5 min)

Terminal environment: `netpol` (netpol k3s cluster). Pre-deployed: namespace `demo` with a 3-tier app (frontend, api, database deployments + services).

### Step 1 — Show the running 3-tier app

```bash
kubectl -n demo get pods,svc
```

### Step 2 — Show that everything can talk to everything (the problem)

```bash
kubectl -n demo exec deployment/frontend -- curl -s --max-time 3 api
kubectl -n demo exec deployment/frontend -- curl -s --max-time 3 database
kubectl -n demo exec deployment/api -- curl -s --max-time 3 database
```

All three succeed — frontend can reach the database directly, which is bad.

### Step 3 — Apply default-deny to lock everything down

```bash
kubectl apply -f /root/manifests/deny-all.yaml
```

Verify everything is now blocked:

```bash
kubectl -n demo exec deployment/frontend -- curl -s --max-time 3 api
```

Expected: timeout / connection refused.

### Step 4 — Restore only the paths that should exist

```bash
kubectl apply -f /root/manifests/allow-frontend-to-api.yaml
kubectl apply -f /root/manifests/allow-api-to-db.yaml
```

Test the allowed paths:

```bash
kubectl -n demo exec deployment/frontend -- curl -s --max-time 3 api
kubectl -n demo exec deployment/api -- curl -s --max-time 3 database
```

Both succeed.

### Step 5 — Confirm frontend still cannot reach database

```bash
kubectl -n demo exec deployment/frontend -- curl -s --max-time 3 database
```

Expected: timeout — blast radius contained.

---

## Slide 10: The Future (no commands)

Talk about eBPF/Cilium, Workload Identity, Talos Linux, SPIFFE/SPIRE. Wrap up, take questions.

---

## Quick Reference

| Demo | Terminal | Cluster Profile | Key Manifests |
|------|----------|----------------|---------------|
| RBAC | `rbac` | base | `manifests/rbac/setup.yaml` |
| Policy | `policy` | kyverno | `manifests/policy/*.yaml` |
| Network Policy | `netpol` | netpol | `manifests/netpol/*.yaml` |
