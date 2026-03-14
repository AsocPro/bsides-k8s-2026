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

// K3sCluster holds the state for a k3s cluster running in Dagger.
type K3sCluster struct {
	Name        string
	ConfigCache *dagger.CacheVolume
	Container   *dagger.Container
}

// NewK3sCluster creates a new k3s cluster configuration.
func NewK3sCluster(name string) *K3sCluster {
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
		ConfigCache: ccache,
		Container:   ctr,
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
// It polls until k3s has written the config.
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

// KubectlContainer returns a container with kubectl configured to talk to this cluster.
// Useful as a base for demo environment containers.
func (k *K3sCluster) KubectlContainer() *dagger.Container {
	return dag.Container().
		From("bitnami/kubectl").
		WithoutEntrypoint().
		WithMountedCache("/cache/k3s", k.ConfigCache).
		WithEnvVariable("CACHE", time.Now().String()).
		WithFile("/.kube/config", k.Config(), dagger.ContainerWithFileOpts{Permissions: 0o644}).
		WithEnvVariable("KUBECONFIG", "/.kube/config")
}

// DemoContainer returns an Alpine container with kubectl, ttyd, and common tools,
// wired to this k3s cluster. It runs ttyd so the terminal is accessible via the proxy.
func (k *K3sCluster) DemoContainer(name string, svc *dagger.Service) *dagger.Container {
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache",
			"ttyd", "bash", "curl", "jq",
			"kubectl", // Alpine package for kubectl
		}).
		WithServiceBinding(name, svc).
		WithFile("/.kube/config", k.Config(), dagger.ContainerWithFileOpts{Permissions: 0o644}).
		WithEnvVariable("KUBECONFIG", "/.kube/config").
		WithEnvVariable("PS1", fmt.Sprintf("[%s] \\w $ ", name)).
		WithExposedPort(7681)
}
