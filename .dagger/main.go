// BSides K8s Interactive Presentation
//
// Dagger module for building and running the interactive Kubernetes security
// presentation. Everything runs via dagger call — no local installs required.

package main

import (
	"context"
	_ "embed"

	"dagger/bsides-k-8-s/internal/dagger"
)

//go:embed scripts/capture.js
var captureScript string

//go:embed scripts/nginx.conf
var nginxConf string

type BsidesK8S struct{}

// BuildFrontend builds the Svelte frontend and returns the dist directory
func (m *BsidesK8S) BuildFrontend(
	// Frontend source directory
	// +defaultPath="./frontend"
	source *dagger.Directory,
) *dagger.Directory {
	return dag.Container().
		From("node:22-alpine").
		WithDirectory("/app", source).
		WithWorkdir("/app").
		WithMountedCache("/root/.npm", dag.CacheVolume("npm-cache")).
		WithExec([]string{"npm", "install"}).
		WithExec([]string{"npm", "run", "build"}).
		Directory("/app/dist")
}

// BuildBackend compiles the Go backend and returns the binary
func (m *BsidesK8S) BuildBackend(
	// Backend source directory
	// +defaultPath="./backend"
	source *dagger.Directory,
) *dagger.File {
	return dag.Container().
		From("golang:1.24-alpine").
		WithDirectory("/app", source).
		WithWorkdir("/app").
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod-cache")).
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("go-build-cache")).
		WithEnvVariable("CGO_ENABLED", "0").
		WithExec([]string{"go", "mod", "tidy"}).
		WithExec([]string{"go", "build", "-o", "server", "."}).
		File("/app/server")
}

// Build produces the complete production artifact: Go binary + static frontend assets
func (m *BsidesK8S) Build(
	// Frontend source directory
	// +defaultPath="./frontend"
	frontendSource *dagger.Directory,
	// Backend source directory
	// +defaultPath="./backend"
	backendSource *dagger.Directory,
) *dagger.Directory {
	frontend := m.BuildFrontend(frontendSource)
	backend := m.BuildBackend(backendSource)

	return dag.Directory().
		WithFile("server", backend).
		WithDirectory("static", frontend)
}

// Terminal creates a basic ttyd container (no k8s) for simple shell demos
func (m *BsidesK8S) Terminal(
	// Name for this terminal environment
	name string,
) *dagger.Service {
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "ttyd", "bash", "curl", "jq"}).
		WithExposedPort(7681).
		AsService(dagger.ContainerAsServiceOpts{
			Args: []string{
				"ttyd",
				"--port", "7681",
				"--writable",
				"--base-path", "/terminal/" + name,
				"bash",
			},
		})
}

// K8sTerminal creates a ttyd container connected to a k3s cluster for live demos.
// The cluster is configured based on the profile (base, kyverno, netpol) and
// pre-seeded with the appropriate manifests.
func (m *BsidesK8S) K8sTerminal(
	ctx context.Context,
	// Name for this environment
	name string,
	// Cluster profile: base, kyverno, or netpol
	// +default="base"
	profile string,
	// Manifests to pre-load into the cluster
	// +optional
	manifests *dagger.Directory,
) (*dagger.Service, error) {
	cluster := NewK3sCluster(name, ClusterProfile(profile), manifests)
	k3sSvc, err := cluster.Server().Start(ctx)
	if err != nil {
		return nil, err
	}

	demo := cluster.DemoContainer(name, k3sSvc)

	return demo.
		AsService(dagger.ContainerAsServiceOpts{
			Args: []string{
				"bash", "-c",
				"bash /usr/local/bin/seed.sh && exec ttyd --port 7681 --writable --base-path /terminal/" + name + " bash",
			},
		}), nil
}

// Dev starts the development environment with Vite HMR, Go backend proxy, and a demo terminal
func (m *BsidesK8S) Dev(
	// Frontend source directory
	// +defaultPath="./frontend"
	frontendSource *dagger.Directory,
	// Backend source directory
	// +defaultPath="./backend"
	backendSource *dagger.Directory,
) *dagger.Service {
	binary := m.BuildBackend(backendSource)

	// Vite dev server with HMR
	vite := dag.Container().
		From("node:22-alpine").
		WithDirectory("/app", frontendSource).
		WithWorkdir("/app").
		WithMountedCache("/root/.npm", dag.CacheVolume("npm-cache")).
		WithExec([]string{"npm", "install"}).
		WithExposedPort(5173).
		AsService(dagger.ContainerAsServiceOpts{
			Args: []string{"npx", "vite", "--host", "0.0.0.0", "--port", "5173"},
		})

	// Basic demo terminal (no k8s) for dev mode — fast startup
	ttydDemo := m.Terminal("demo")

	return dag.Container().
		From("alpine:latest").
		WithFile("/app/server", binary).
		WithWorkdir("/app").
		WithServiceBinding("vite", vite).
		WithServiceBinding("ttyd-demo", ttydDemo).
		WithEnvVariable("VITE_DEV_URL", "http://vite:5173").
		WithEnvVariable("TTYD_URLS", "demo=http://ttyd-demo:7681").
		WithExposedPort(8080).
		AsService(dagger.ContainerAsServiceOpts{
			Args: []string{"/app/server"},
		})
}

// Present runs the full presentation with k3s-backed demo environments.
// Each demo section gets its own k3s cluster with the appropriate profile
// and pre-seeded manifests. The backend container additionally gets kubectl
// and goss for the state-check API (POST /api/state/{demo}).
func (m *BsidesK8S) Present(
	ctx context.Context,
	// Frontend source directory
	// +defaultPath="./frontend"
	frontendSource *dagger.Directory,
	// Backend source directory
	// +defaultPath="./backend"
	backendSource *dagger.Directory,
	// Manifests directory containing per-demo subdirectories (rbac/, policy/, netpol/)
	// +defaultPath="./manifests"
	manifestsDir *dagger.Directory,
	// Goss test definitions
	// +defaultPath="./goss"
	gossDir *dagger.Directory,
) (*dagger.Service, error) {
	binary := m.BuildBackend(backendSource)
	static := m.BuildFrontend(frontendSource)

	// Create clusters inline so we can access both ttyd services AND kubeconfigs
	type clusterInfo struct {
		name    string
		profile ClusterProfile
	}

	clusters := []clusterInfo{
		{name: "rbac", profile: ProfileBase},
		{name: "policy", profile: ProfileKyverno},
		{name: "netpol", profile: ProfileNetpol},
	}

	ttydServices := make(map[string]*dagger.Service)
	k3sServices := make(map[string]*dagger.Service)
	kubeconfigs := make(map[string]*dagger.File)

	for _, ci := range clusters {
		cluster := NewK3sCluster(ci.name, ci.profile, manifestsDir.Directory(ci.name))
		k3sSvc, err := cluster.Server().Start(ctx)
		if err != nil {
			return nil, err
		}

		k3sServices[ci.name] = k3sSvc
		kubeconfigs[ci.name] = cluster.Config()

		demo := cluster.DemoContainer(ci.name, k3sSvc)
		ttydServices[ci.name] = demo.AsService(dagger.ContainerAsServiceOpts{
			Args: []string{
				"bash", "-c",
				"bash /usr/local/bin/seed.sh && exec ttyd --port 7681 --writable --base-path /terminal/" + ci.name + " bash",
			},
		})
	}

	// Build the backend container with ttyd proxies + goss state checks
	backend := dag.Container().
		From("alpine:latest").
		// Install kubectl, curl, and goss for state checks
		WithExec([]string{"apk", "add", "--no-cache", "kubectl", "curl"}).
		WithExec([]string{"sh", "-c", "curl -fsSL https://goss.rocks/install | GOSS_DST=/usr/local/bin sh"}).
		WithFile("/app/server", binary).
		WithDirectory("/app/static", static).
		WithDirectory("/app/goss", gossDir).
		WithWorkdir("/app").
		// Bind ttyd services
		WithServiceBinding("ttyd-rbac", ttydServices["rbac"]).
		WithServiceBinding("ttyd-policy", ttydServices["policy"]).
		WithServiceBinding("ttyd-netpol", ttydServices["netpol"]).
		// Bind k3s services so kubeconfig hostnames resolve
		WithServiceBinding("rbac", k3sServices["rbac"]).
		WithServiceBinding("policy", k3sServices["policy"]).
		WithServiceBinding("netpol", k3sServices["netpol"]).
		// Mount kubeconfigs for goss tests
		WithFile("/app/kubeconfig-rbac", kubeconfigs["rbac"]).
		WithFile("/app/kubeconfig-policy", kubeconfigs["policy"]).
		WithFile("/app/kubeconfig-netpol", kubeconfigs["netpol"]).
		WithEnvVariable("KUBECONFIG_RBAC", "/app/kubeconfig-rbac").
		WithEnvVariable("KUBECONFIG_POLICY", "/app/kubeconfig-policy").
		WithEnvVariable("KUBECONFIG_NETPOL", "/app/kubeconfig-netpol").
		WithEnvVariable("STATIC_DIR", "/app/static").
		WithEnvVariable("TTYD_URLS", "rbac=http://ttyd-rbac:7681,netpol=http://ttyd-netpol:7681,policy=http://ttyd-policy:7681").
		WithExposedPort(8080)

	return backend.AsService(dagger.ContainerAsServiceOpts{
		Args: []string{"/app/server"},
	}), nil
}

// K3sDebug gives you an interactive terminal inside the k3s container for debugging.
// Note: cgroups will be read-only in terminal mode, but you can still run k3s
// manually to see the actual error output.
// Usage: dagger call k3s-debug terminal
func (m *BsidesK8S) K3sDebug(
	// Name for this cluster
	// +default="debug"
	name string,
	// Cluster profile: base, kyverno, or netpol
	// +default="base"
	profile string,
	// Manifests to pre-load
	// +optional
	manifests *dagger.Directory,
) *dagger.Container {
	cluster := NewK3sCluster(name, ClusterProfile(profile), manifests)
	return cluster.Container
}

// Serve runs the production presentation with pre-built assets and a basic demo terminal
func (m *BsidesK8S) Serve(
	// Frontend source directory
	// +defaultPath="./frontend"
	frontendSource *dagger.Directory,
	// Backend source directory
	// +defaultPath="./backend"
	backendSource *dagger.Directory,
) *dagger.Service {
	binary := m.BuildBackend(backendSource)
	static := m.BuildFrontend(frontendSource)

	ttydDemo := m.Terminal("demo")

	return dag.Container().
		From("alpine:latest").
		WithFile("/app/server", binary).
		WithDirectory("/app/static", static).
		WithWorkdir("/app").
		WithServiceBinding("ttyd-demo", ttydDemo).
		WithEnvVariable("STATIC_DIR", "/app/static").
		WithEnvVariable("TTYD_URLS", "demo=http://ttyd-demo:7681").
		WithExposedPort(8080).
		AsService(dagger.ContainerAsServiceOpts{
			Args: []string{"/app/server"},
		})
}

// ExportPdf renders the presentation as a PDF with a first-page link to the
// GitHub repo. Everything runs headlessly inside Dagger — Playwright drives
// a Chromium session against an nginx-served copy of the Svelte build,
// stepping through each slide (including substep animations) and capturing
// a screenshot per slide. The screenshots are then stitched into a single
// PDF.
//
// Usage:
//
//	dagger call export-pdf export --path bsides-k8s-2026.pdf
func (m *BsidesK8S) ExportPdf(
	ctx context.Context,
	// Frontend source directory
	// +defaultPath="./frontend"
	frontendSource *dagger.Directory,
	// GitHub repo URL shown on the cover page
	// +default="https://github.com/AsocPro/bsides-k8s-2026"
	repoUrl string,
) *dagger.File {
	// Build the Svelte frontend and serve it as static files. The real
	// backend (ttyd + /api/state) isn't available during capture — the
	// embedded nginx config returns friendly placeholders for those paths
	// so the demo slide iframes don't render as 404 pages.
	frontend := m.BuildFrontend(frontendSource)

	server := dag.Container().
		From("nginx:alpine").
		WithDirectory("/usr/share/nginx/html", frontend).
		WithNewFile("/etc/nginx/conf.d/default.conf", nginxConf).
		WithExposedPort(80).
		AsService(dagger.ContainerAsServiceOpts{})

	// Playwright image version must match the npm package version so the
	// preinstalled browsers at PLAYWRIGHT_BROWSERS_PATH are reused instead
	// of being downloaded again.
	const playwrightVersion = "1.49.0"

	capture := dag.Container().
		From("mcr.microsoft.com/playwright:v"+playwrightVersion+"-jammy").
		WithServiceBinding("presentation", server).
		WithEnvVariable("BASE_URL", "http://presentation/").
		WithEnvVariable("REPO_URL", repoUrl).
		WithEnvVariable("OUT_DIR", "/work/out").
		WithWorkdir("/work").
		WithNewFile("/work/capture.js", captureScript).
		WithNewFile("/work/package.json", `{"name":"bsides-pdf-capture","version":"1.0.0","private":true,"dependencies":{"playwright":"`+playwrightVersion+`"}}`).
		WithMountedCache("/root/.npm", dag.CacheVolume("pdf-capture-npm-cache")).
		WithExec([]string{"sh", "-c", "apt-get update -qq && apt-get install -y -qq --no-install-recommends img2pdf poppler-utils"}).
		WithExec([]string{"npm", "install", "--no-audit", "--no-fund", "--prefer-offline"}).
		WithExec([]string{"node", "capture.js"}).
		// Slides are PNG screenshots (for fidelity with GSAP animations and
		// iframe placeholders); the title page is a real PDF so its repo
		// link is a clickable annotation. img2pdf builds the slide pages
		// losslessly, then pdfunite concatenates title + slides.
		WithExec([]string{"sh", "-c", "img2pdf $(ls /work/out/page-*.png | sort) --output /work/slides.pdf"}).
		WithExec([]string{"pdfunite", "/work/out/title.pdf", "/work/slides.pdf", "/work/bsides-k8s-2026.pdf"})

	return capture.File("/work/bsides-k8s-2026.pdf")
}
