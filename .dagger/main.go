// BSides K8s Interactive Presentation
//
// Dagger module for building and running the interactive Kubernetes security
// presentation. Everything runs via dagger call — no local installs required.

package main

import (
	"dagger/bsides-k-8-s/internal/dagger"
)

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

// Terminal creates a ttyd container that serves a shell over HTTP/WebSocket
func (m *BsidesK8S) Terminal(
	// Name for this terminal environment
	name string,
) *dagger.Service {
	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "ttyd", "bash", "curl", "jq"}).
		WithExposedPort(7681).
		WithDefaultArgs([]string{
			"ttyd",
			"--port", "7681",
			"--writable",
			"--base-path", "/terminal/" + name,
			"bash",
		}).
		AsService()
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
		WithDefaultArgs([]string{"npx", "vite", "--host", "0.0.0.0", "--port", "5173"}).
		AsService()

	// Demo terminal via ttyd
	ttydDemo := m.Terminal("demo")

	// Run pre-built backend binary proxying to Vite and ttyd
	return dag.Container().
		From("alpine:latest").
		WithFile("/app/server", binary).
		WithWorkdir("/app").
		WithServiceBinding("vite", vite).
		WithServiceBinding("ttyd-demo", ttydDemo).
		WithEnvVariable("VITE_DEV_URL", "http://vite:5173").
		WithEnvVariable("TTYD_URLS", "demo=http://ttyd-demo:7681").
		WithExposedPort(8080).
		WithDefaultArgs([]string{"/app/server"}).
		AsService()
}

// Serve runs the production presentation with pre-built assets and a demo terminal
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
		WithDefaultArgs([]string{"/app/server"}).
		AsService()
}
