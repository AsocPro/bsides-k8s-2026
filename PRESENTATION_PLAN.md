# Presentation Plan: Kubernetes — Orchestrating a More Secure Future

## Format
- 25-minute talk at BSides
- Web-based interactive presentation (not traditional slides)
- Embedded live terminals running in Dagger containers
- Reactive overlays/animations triggered by real demo state changes
- Goal: "How did you do that?" level polish without distracting from content

---

## Architecture

```
┌─────────────────────────────────────────────────┐
│                   Browser                        │
│  ┌───────────────┐  ┌────────────────────────┐  │
│  │  Slide Engine  │  │  Embedded Terminal(s)  │  │
│  │  (animations,  │  │  (xterm.js + WebSocket │  │
│  │   overlays,    │◄─┤   to Dagger shells)    │  │
│  │   transitions) │  │                        │  │
│  └───────▲───────┘  └────────────────────────┘  │
│          │                                       │
└──────────┼───────────────────────────────────────┘
           │ WebSocket events
┌──────────┴───────────────────────────────────────┐
│              Event Bridge (Backend)               │
│  - Watches demo environments via Goss/kubectl     │
│  - Emits events: "pod-created", "policy-blocked", │
│    "network-denied", "rbac-denied", etc.          │
│  - Drives slide transitions & overlay triggers    │
└──────────────────────────────────────────────────┘
           │
┌──────────┴───────────────────────────────────────┐
│              Dagger Runtime                       │
│  - Spins up ephemeral K8s clusters (k3s/kind)    │
│  - Pre-warms environments before each section    │
│  - Tears down after section completes             │
└──────────────────────────────────────────────────┘
```

---

## Sections & Timing

### 1. Opening — "How We Got Here" (4 min)
**Content:** Evolution from bare metal → VMs → containers, manual → scripts → IaC

**Presentation Ideas:**
- Animated timeline that builds itself as you talk — each era "stacks" on the previous
- Visual metaphor: a server rack that transforms — physical box → VM boxes inside it → tiny container cubes
- Each era slides in from the right, previous era fades/shrinks but stays visible
- Pain points appear as red warning icons that accumulate (config drift, inconsistency, etc.)
- When you reach containers + IaC, the warnings resolve with satisfying green checkmarks
- No demo here — pure visual storytelling, fast-paced

**Wow moment:** The timeline is one continuous animation, not discrete slides. It feels like a motion graphic, not a slideshow.

---

### 2. Kubernetes 101 — "The Orchestrator" (3 min)
**Content:** Quick intro to pods, deployments, services, namespaces — just enough to follow the demos

**Presentation Ideas:**
- Interactive cluster visualization: nodes appear, pods get scheduled onto them in real-time
- As you mention each concept, it highlights/animates in the cluster diagram
- Pods pulse, services draw connection lines, namespaces draw boundary boxes
- Keep it snappy — this is scaffolding for the audience, not the main event

**Wow moment:** The cluster diagram is live — it's actually querying a real (pre-warmed) cluster via the event bridge. Pods you see appearing are real pods. Casual mention: "oh by the way, that's a real cluster."

**Demo setup:** Dagger pre-warms a k3s cluster during the opening section. Event bridge starts streaming pod/node state to the visualization.

---

### 3. RBAC Deep Dive — "Who Can Do What" (5 min)
**Content:** Roles, RoleBindings, ServiceAccounts, least privilege

**Presentation Ideas:**
- Split screen: terminal on right, animated permission matrix on left
- The permission matrix updates live as you create roles/bindings in the terminal
- When you attempt a forbidden action, the matrix cell flashes red and a "DENIED" overlay animates in
- When you grant permission and retry, it flashes green with a "GRANTED" overlay

**Demo flow:**
1. Show a service account with no permissions — try to list pods → denied
2. Create a Role allowing pod list — matrix animates the new permission appearing
3. Create a RoleBinding — connection line draws between the SA and Role
4. Retry list pods → success, matrix cell turns green
5. Try to delete a pod → denied (not in the role) — show least privilege in action
6. Try to access another namespace → denied — show namespace isolation

**Wow moment:** The permission matrix is not a static graphic — it's reacting to actual RBAC state in the cluster. The audience sees the kubectl command, the denial, AND the visual representation update simultaneously. The event bridge watches `kubectl auth can-i` results and pushes state changes.

**Event bridge triggers:**
- `rbac-deny` → red flash overlay + denied animation
- `rbac-allow` → green flash overlay + granted animation
- `role-created` → matrix row appears with slide-in animation
- `rolebinding-created` → connection line draws between entities

---

### 4. Policy Agents — "What Changes Are Allowed" (5 min)
**Content:** OPA/Kyverno admission control, practical policies

**Presentation Ideas:**
- Visual "gateway" metaphor — API requests flow toward the cluster, policy agent is a checkpoint
- Requests that pass float through smoothly; rejected ones bounce back with the violation reason
- Show the policy YAML briefly, then demonstrate it live

**Demo flow:**
1. Deploy a Kyverno policy requiring non-root containers
2. Try to deploy a container running as root → rejected
   - The request visually hits the gateway and bounces back
   - Violation reason appears as an overlay annotation
3. Fix the deployment (add securityContext) → accepted
   - Request flows through the gateway smoothly
4. Deploy a policy requiring approved registries
5. Try to deploy from `random-registry.io` → rejected
6. Deploy from approved registry → accepted

**Wow moment:** The "gateway" visualization shows requests in flight. When you hit enter on a kubectl apply, the audience sees the request travel toward the cluster, hit the policy checkpoint, and either pass or bounce — all animated in real-time synced to the actual API response. The rejection reason from Kyverno appears as a floating annotation on the bounced request.

**Event bridge triggers:**
- `admission-reject` → request bounces back animation + reason overlay
- `admission-allow` → request passes through animation
- `policy-created` → new rule appears at the gateway checkpoint

---

### 5. Network Policies — "Limiting Blast Radius" (5 min)
**Content:** Microsegmentation, default deny, explicit allow

**Presentation Ideas:**
- Network topology diagram showing frontend → API → database
- Connection lines are animated data flows (particles moving along paths)
- As you apply network policies, unauthorized paths visually "cut" with an animation
- Attempted connections on blocked paths show red blocked indicators

**Demo flow:**
1. Start with 3-tier app deployed, all pods can talk to all pods
   - Topology shows all connections active (particles flowing everywhere)
2. Show a compromised frontend reaching the database directly — this is the problem
   - Malicious flow highlighted in red/orange
3. Apply default-deny network policy
   - ALL connection lines cut simultaneously — dramatic visual moment
   - App breaks (show it in terminal)
4. Apply targeted allow policies one by one
   - Frontend → API: connection restores with green particles
   - API → Database: connection restores
   - Frontend → Database: stays blocked
5. Retry the direct frontend→database attack → blocked
   - Red particles hit a wall and dissipate

**Wow moment:** The network topology is the star here. Particles flowing along paths, connections cutting and restoring in real-time synced to actual network policy application. When the default-deny hits and everything goes dark, then connections restore one by one — it's a visceral demonstration of microsegmentation. The event bridge uses network policy watches and actual connectivity tests (curl between pods) to drive the visualization state.

**Event bridge triggers:**
- `netpol-applied` → affected connection lines animate cut/restore
- `connectivity-test-pass` → green particles flow
- `connectivity-test-fail` → red particles hit wall
- `default-deny-applied` → all lines cut simultaneously (dramatic moment)

---

### 6. Looking Forward — "The Future" (2 min)
**Content:** eBPF, workload identity, zero trust

**Presentation Ideas:**
- Quick visual montage — each technology gets a brief animated card
- eBPF: show kernel-level visibility concept (layers peeling back)
- Workload identity: credentials appearing and expiring on a timeline
- Zero trust: everything from the talk converging into a unified security posture diagram
- End with the full architecture from the talk assembled — RBAC + Policies + Network Policies + Future tech all layered together

**No demo here** — bring the energy back to the big picture and close strong.

---

### 7. Close (1 min)
- Quick recap visual — the full security stack assembled
- Resources/links
- Q&A prompt

---

## Technical Stack

> **Constraint: Zero local installs.** Everything — frontend build, Go backend, K8s environments, dev workflow — runs via `dagger call`. The only prerequisite on the presenter's machine is Dagger itself.

### Frontend — Svelte + Vite + Tailwind
- **Svelte** (not SvelteKit — this is a single-page presentation, no SSR needed) + **Vite** for dev/build
- **Tailwind CSS** for styling — utility classes keep animation-heavy components clean
- **GSAP** for timeline-based animations (scrub, pause, sequence)
- **ttyd** for embedded terminals — runs inside each Dagger demo container, Go backend reverse-proxies to it. Handles PTY + WebSocket + xterm.js rendering out of the box
- Custom slide engine built in Svelte — each section is a component with its own animation timeline and event handlers
- Built inside a Dagger container — `dagger call build-frontend` produces the static assets

### Backend — Go Proxy Server
The Go backend is the central hub. It serves the frontend, proxies terminal sessions, and runs the event bridge.

```
┌─────────────────────────────────────────────────────────┐
│                    Go Backend                            │
│                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────────┐  │
│  │ Static File  │  │  Terminal    │  │  Event Bridge  │  │
│  │ Server       │  │  Reverse     │  │                │  │
│  │ (Svelte app) │  │  Proxy→ttyd  │  │  (WebSocket)   │  │
│  └──────────────┘  └──────┬───────┘  └───────┬───────┘  │
│                           │                  │           │
│                    ┌──────┴───────┐   ┌──────┴───────┐  │
│                    │ Dagger       │   │ K8s Watchers  │  │
│                    │ Containers   │   │ JSONL Tailer  │  │
│                    │ running ttyd │   │ Goss Checks   │  │
│                    └──────────────┘   └──────────────┘  │
└─────────────────────────────────────────────────────────┘
```

**Routes:**
- `GET /` — serves the Svelte presentation app
- `/terminal/:env/*` — reverse proxies to the ttyd instance running inside the Dagger container for that environment (HTTP + WebSocket upgrade)
- `WS /ws/events` — streams event bridge events to the frontend for reactive animations
- `GET /api/environments` — lists available demo environments and their readiness status
- `POST /api/environments/:env/start` — triggers Dagger to spin up a specific demo environment
- `POST /api/environments/:env/teardown` — tears down a demo environment

**Terminal Proxy Details:**
- Each Dagger demo container runs **ttyd** (lightweight terminal server) — it handles PTY allocation, xterm.js client, and WebSocket transport natively
- The Go backend reverse-proxies `/terminal/:env/` to the ttyd port inside the corresponding Dagger container
- The Svelte frontend embeds each terminal as an iframe or directly uses ttyd's xterm.js client via the proxied path
- Multiple simultaneous terminals are supported (each env runs its own ttyd instance on a unique port)
- ttyd stdout can be teed to the JSONL command log via a shell wrapper so typed commands become events

**Event Bridge Details:**
- Runs as goroutines within the same Go process
- K8s informers/watches for RBAC, NetworkPolicy, Pod, and Kyverno policy resources
- JSONL file tailer (fsnotify or tail -f equivalent) for command log events
- Optional Goss runner for connectivity/state checks
- All events are fanned out to connected WebSocket clients on `/ws/events`

### Demo Runtime — Everything in Dagger
- **All environments run inside Dagger containers** — the Go backend, the K8s clusters (k3s), the build pipeline
- Dagger pipelines define each demo environment as a composable module
- Environments can be pre-warmed in parallel before the talk starts
- Each demo section gets its own isolated k3s instance inside a Dagger container
- The Go backend itself runs in Dagger — nothing is installed locally, everything is a `dagger call`
- Development, building, and running are all `dagger call` commands — no local Go, Node, npm, etc. required

**Example `dagger call` workflow:**
```bash
# Development — hot-reload frontend + backend
dagger call dev up

# Build everything
dagger call build

# Pre-warm all demo environments before the talk
dagger call environments up

# Run the full presentation (backend + frontend + all envs)
dagger call present up

# Tear down everything
dagger call environments down
```

**Environment lifecycle:**
```
[Pre-talk] `dagger call environments up` spins up all environments in parallel
    → k3s clusters boot inside Dagger containers
    → base workloads deploy
    → readiness checks pass
    → environments report "ready" to Go backend

[During talk] Presenter advances to RBAC section
    → frontend connects terminal to rbac env
    → event bridge starts watching rbac env's K8s API
    → demo proceeds

[Section end] Environment can be torn down or left running
    → next environment was already pre-warmed
```

### Command Log Watcher (JSONL)
A system that logs executed commands to a JSONL file provides a second event source for the event bridge. This is useful in several ways:

- **Universal trigger source** — Not everything emits a K8s event. A `kubectl apply -f policy.yaml` command being run is itself a meaningful event, even before the cluster processes it. The event bridge can tail the JSONL file and fire a "command-issued" event immediately, letting the presentation show a "deploying policy..." animation the instant you hit enter, before the K8s watch picks up the result.
- **Command echo overlay** — The presentation can display a stylized version of the command you just ran as a floating overlay or annotation, so the audience doesn't have to squint at the terminal. The JSONL watcher parses the command and the presentation renders it large and readable.
- **Two-phase animations** — Phase 1 triggers on the command log (intent: "applying network policy..."), Phase 2 triggers on the K8s watch (result: "network policy active"). This creates a natural request→response animation flow that matches what's actually happening.
- **Fallback for unobservable state** — Some things are hard to watch via K8s informers (e.g., a curl failing between pods for the network policy demo). If the command log captures `kubectl exec frontend -- curl api:8080` and the exit code, the event bridge can derive connectivity pass/fail without needing a separate Goss check.
- **Demo pacing** — The event bridge can track which commands have been run to know where you are in the demo script. If you're on step 3 of 5, the presentation can subtly highlight what's coming next or pre-load the next animation.
- **Rehearsal & debugging** — Replay the JSONL file to re-run the entire presentation's reactive layer without touching a cluster. Great for testing animations and timing.

### Resilience
- Every demo should have a pre-recorded fallback (seamless switch if something breaks)
- Event bridge should have timeout handling — if a demo step takes too long, show a graceful loading state
- Pre-warm ALL environments before the talk starts, not just-in-time
- Test the full flow on the actual presentation hardware before the talk

---

## Key Design Principles

1. **The demo IS the slide** — don't context-switch between "now look at the screen" and "now look at my terminal." They're the same thing.
2. **Animations serve understanding** — every animation should make a concept clearer, not just look cool. The network particles show data flow. The RBAC matrix shows permissions. The gateway shows admission control.
3. **Reactive, not scripted** — the overlays respond to real state, which means if something unexpected happens, the presentation adapts. This is more impressive AND more honest than pre-baked animations.
4. **Graceful degradation** — if a demo breaks, the fallback is seamless. The audience should never see you scrambling.
5. **25 minutes is short** — every second counts. No filler slides, no "agenda" slide, no "about me" unless it's <15 seconds. Jump straight into the timeline.
