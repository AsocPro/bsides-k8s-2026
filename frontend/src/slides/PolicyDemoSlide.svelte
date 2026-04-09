<script>
  import Slide from '../lib/Slide.svelte';
  import SplitLayout from '../lib/SplitLayout.svelte';
  import Terminal from '../lib/Terminal.svelte';
  import DemoStep from '../lib/DemoStep.svelte';
  import StatePanel from '../lib/StatePanel.svelte';
  import PolicyGate from '../lib/PolicyGate.svelte';
  import { policyState } from '../lib/stores/stateCheck.js';

  let openStep = $state(-1);
  function toggle(step) { openStep = openStep === step ? -1 : step; }
</script>

<Slide padding={false}>
  <div class="p-8 pb-2">
    <h2 class="text-3xl font-bold text-white">Policy Agent Demo</h2>
    <p class="text-neutral-500 text-sm mt-1">Live — Kyverno blocking non-compliant resources</p>
  </div>

  <div class="flex-1 min-h-0 px-8 pb-8">
    <SplitLayout ratio="2fr 3fr">
      {#snippet left()}
        <div class="flex flex-col gap-2 text-sm pt-4 h-full min-h-0">
          <StatePanel demo="policy" store={policyState} color="amber">
            <PolicyGate data={$policyState.data} />
          </StatePanel>

          <div class="space-y-2 overflow-y-auto min-h-0 flex-1">
          <DemoStep step={1} color="amber"
            label="Apply Kyverno policies"
            expanded={openStep === 1} ontoggle={toggle}
            commands={[
              'kubectl apply -f /root/manifests/require-non-root.yaml',
              'kubectl apply -f /root/manifests/require-approved-registries.yaml',
              'kubectl get clusterpolicy',
            ]}
          />
          <DemoStep step={2} color="amber"
            label="Inspect a policy"
            expanded={openStep === 2} ontoggle={toggle}
            commands={[
              'kubectl get clusterpolicy require-non-root -o yaml',
            ]}
          />
          <DemoStep step={3} color="amber"
            label="Try deploying a root container — rejected"
            expanded={openStep === 3} ontoggle={toggle}
            commands={[
              'kubectl run bad-pod --image=nginx:alpine -n demo',
            ]}
          />
          <DemoStep step={4} color="amber"
            label="Deploy a non-root container — accepted"
            expanded={openStep === 4} ontoggle={toggle}
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
          <DemoStep step={5} color="amber"
            label="Test registry policy with unapproved image"
            expanded={openStep === 5} ontoggle={toggle}
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
        </div>
      {/snippet}
      {#snippet right()}
        <Terminal name="policy" class="w-full h-full" />
      {/snippet}
    </SplitLayout>
  </div>
</Slide>
