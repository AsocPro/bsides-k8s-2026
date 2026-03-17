<script>
  import Slide from '../lib/Slide.svelte';
  import SplitLayout from '../lib/SplitLayout.svelte';
  import Terminal from '../lib/Terminal.svelte';
  import DemoStep from '../lib/DemoStep.svelte';
</script>

<Slide padding={false}>
  <div class="p-8 pb-2">
    <h2 class="text-3xl font-bold text-white">Policy Agent Demo</h2>
    <p class="text-neutral-500 text-sm mt-1">Live — Kyverno blocking non-compliant resources</p>
  </div>

  <div class="flex-1 min-h-0 px-8 pb-8">
    <SplitLayout ratio="2fr 3fr">
      {#snippet left()}
        <div class="space-y-2 text-sm pt-4 overflow-y-auto">
          <DemoStep step={1} color="amber"
            label="Inspect the active policies"
            commands={[
              'kubectl get clusterpolicy',
              'kubectl get clusterpolicy require-non-root -o yaml',
            ]}
          />
          <DemoStep step={2} color="amber"
            label="Try deploying a root container — rejected"
            commands={[
              'kubectl run bad-pod --image=nginx:alpine -n demo',
            ]}
          />
          <DemoStep step={3} color="amber"
            label="Deploy a non-root container — accepted"
            commands={[
              `kubectl apply -f - <<'EOF'
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
EOF`,
              'kubectl -n demo get pods',
            ]}
          />
          <DemoStep step={4} color="amber"
            label="Test registry policy with unapproved image"
            commands={[
              `kubectl apply -f - <<'EOF'
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
EOF`,
              `kubectl apply -f - <<'EOF'
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
EOF`,
            ]}
          />
        </div>
      {/snippet}
      {#snippet right()}
        <Terminal name="policy" class="w-full h-full" />
      {/snippet}
    </SplitLayout>
  </div>
</Slide>
