package main

import (
	"fmt"
	"time"

	"dagger/bsides-k-8-s/internal/dagger"
)

// Entrypoint script to setup cgroup v2 nesting for k3s.
// Copied from Dagger's own k3s module (sourced from Docker's DinD).
const k3sEntrypoint = `#!/bin/sh
set -o errexit
set -o nounset

if [ -f /sys/fs/cgroup/cgroup.controllers ]; then
  echo "[$(date -Iseconds)] [CgroupV2 Fix] Evacuating Root Cgroup ..."
  mkdir -p /sys/fs/cgroup/init
  xargs -rn1 < /sys/fs/cgroup/cgroup.procs > /sys/fs/cgroup/init/cgroup.procs || :
  sed -e 's/ / +/g' -e 's/^/+/' <"/sys/fs/cgroup/cgroup.controllers" >"/sys/fs/cgroup/cgroup.subtree_control"
  echo "[$(date -Iseconds)] [CgroupV2 Fix] Done"
fi

exec "$@"
`

// ClusterProfile defines what components a k3s cluster should have.
type ClusterProfile string

const (
	ProfileBase    ClusterProfile = "base"    // Plain k3s, RBAC demos
	ProfileKyverno ClusterProfile = "kyverno" // k3s + Kyverno for policy demos
	ProfileNetpol  ClusterProfile = "netpol"  // k3s with network policy support for netpol demos
)

// K3sCluster holds the state for a k3s cluster running in Dagger.
type K3sCluster struct {
	Name        string
	Profile     ClusterProfile
	ConfigCache *dagger.CacheVolume
	Container   *dagger.Container
	Manifests   *dagger.Directory
}

// NewK3sCluster creates a new k3s cluster with the given profile.
func NewK3sCluster(name string, profile ClusterProfile, manifests *dagger.Directory) *K3sCluster {
	ccache := dag.CacheVolume("k3s_config_" + name)

	ctr := dag.Container().
		From("rancher/k3s:latest").
		WithNewFile("/usr/bin/entrypoint.sh", k3sEntrypoint, dagger.ContainerWithNewFileOpts{
			Permissions: 0o755,
		}).
		WithEntrypoint([]string{"entrypoint.sh"}).
		WithMountedCache("/etc/rancher/k3s", ccache).
		WithMountedTemp("/etc/lib/cni").
		WithMountedTemp("/var/lib/kubelet").
		WithMountedCache("/var/lib/rancher", dag.CacheVolume("k3s_cache_"+name)).
		WithEnvVariable("CACHEBUST", time.Now().String()).
		WithExec([]string{"rm", "-rf", "/var/lib/rancher/k3s/server/tls", "/etc/rancher/k3s/k3s.yaml"}).
		WithExec([]string{"rm", "-rf", "/var/lib/rancher/k3s/"}).
		WithMountedTemp("/var/log").
		WithExposedPort(6443)

	return &K3sCluster{
		Name:        name,
		Profile:     profile,
		ConfigCache: ccache,
		Container:   ctr,
		Manifests:   manifests,
	}
}

// Server returns the k3s server as a Dagger service.
func (k *K3sCluster) Server() *dagger.Service {
	return k.Container.AsService(dagger.ContainerAsServiceOpts{
		Args: []string{
			"sh", "-c",
			"k3s server --debug --bind-address $(ip route | grep src | awk '{print $NF}') --disable traefik --disable metrics-server --egress-selector-mode=disabled",
		},
		InsecureRootCapabilities: true,
		UseEntrypoint:            true,
	})
}

// Config returns the kubeconfig file for this cluster.
func (k *K3sCluster) Config() *dagger.File {
	return dag.Container().
		From("alpine").
		WithEnvVariable("CACHE", time.Now().String()).
		WithMountedCache("/cache/k3s", k.ConfigCache).
		WithExec([]string{"sh", "-c",
			`while [ ! -f "/cache/k3s/k3s.yaml" ]; do echo "waiting for k3s config..." && sleep 0.5; done`,
		}).
		WithExec([]string{"cp", "/cache/k3s/k3s.yaml", "k3s.yaml"}).
		File("k3s.yaml")
}

// seedScript returns a shell script that waits for the cluster to be ready,
// installs profile-specific components, and applies setup manifests.
func (k *K3sCluster) seedScript() string {
	base := `#!/bin/bash
set -e

echo "=== Waiting for cluster to be ready ==="
until kubectl get nodes 2>/dev/null | grep -q " Ready"; do
  echo "waiting for node..."
  sleep 2
done
echo "=== Node ready ==="
kubectl get nodes
`

	var profileSetup string
	switch k.Profile {
	case ProfileKyverno:
		profileSetup = `
echo "=== Installing Kyverno ==="
kubectl create -f https://github.com/kyverno/kyverno/releases/download/v1.13.4/install.yaml || true

echo "=== Waiting for Kyverno to be ready ==="
kubectl -n kyverno rollout status deployment kyverno-admission-controller --timeout=120s
kubectl -n kyverno rollout status deployment kyverno-background-controller --timeout=60s
echo "=== Kyverno ready ==="
`
	case ProfileNetpol:
		// k3s includes a built-in network policy controller by default
		// Just verify it's working
		profileSetup = `
echo "=== Network policy controller is built-in to k3s ==="
echo "=== Verifying kube-system pods ==="
kubectl -n kube-system get pods
`
	default:
		profileSetup = ""
	}

	manifests := `
echo "=== Applying setup manifests ==="
if [ -d /manifests ] && ls /manifests/*.yaml 1>/dev/null 2>&1; then
  kubectl apply -f /manifests/
  echo "=== Waiting for pods to be ready ==="
  kubectl -n demo wait --for=condition=Ready pods --all --timeout=120s 2>/dev/null || true
fi

echo "=== Setup complete ==="
kubectl -n demo get all 2>/dev/null || echo "no demo namespace resources"
`

	return base + profileSetup + manifests
}

// DemoContainer returns an Alpine container with kubectl, ttyd, and common tools,
// wired to this k3s cluster. It runs the seed script, then starts ttyd.
func (k *K3sCluster) DemoContainer(name string, svc *dagger.Service) *dagger.Container {
	ctr := dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache",
			"ttyd", "bash", "curl", "jq", "kubectl",
		}).
		WithServiceBinding(name, svc).
		WithFile("/.kube/config", k.Config(), dagger.ContainerWithFileOpts{Permissions: 0o644}).
		WithEnvVariable("KUBECONFIG", "/.kube/config").
		WithEnvVariable("PS1", fmt.Sprintf("[%s] \\w $ ", name))

	// Mount manifests for the demo section
	if k.Manifests != nil {
		ctr = ctr.WithDirectory("/manifests", k.Manifests)
	}

	// Run the seed script to set up the cluster
	ctr = ctr.
		WithNewFile("/usr/local/bin/seed.sh", k.seedScript(), dagger.ContainerWithNewFileOpts{
			Permissions: 0o755,
		}).
		WithExec([]string{"bash", "/usr/local/bin/seed.sh"})

	// Also copy manifests to home for easy access during demo
	if k.Manifests != nil {
		ctr = ctr.WithDirectory("/root/manifests", k.Manifests)
	}

	return ctr.WithExposedPort(7681)
}
