<script>
  import Slide from '../lib/Slide.svelte';
  import SplitLayout from '../lib/SplitLayout.svelte';
  import Terminal from '../lib/Terminal.svelte';
  import DemoStep from '../lib/DemoStep.svelte';
  import StatePanel from '../lib/StatePanel.svelte';
  import RBACMatrix from '../lib/RBACMatrix.svelte';
  import { rbacState } from '../lib/stores/stateCheck.js';

  let openStep = $state(-1);
  function toggle(step) { openStep = openStep === step ? -1 : step; }
</script>

<Slide padding={false}>
  <div class="p-8 pb-2">
    <h2 class="text-3xl font-bold text-white">RBAC Demo</h2>
    <p class="text-neutral-500 text-sm mt-1">Live — try creating roles and testing access</p>
  </div>

  <div class="flex-1 min-h-0 px-8 pb-8">
    <SplitLayout ratio="2fr 3fr">
      {#snippet left()}
        <div class="flex flex-col gap-2 text-sm pt-4 h-full min-h-0">
          <StatePanel demo="rbac" store={rbacState} color="cyan">
            <RBACMatrix data={$rbacState.data} />
          </StatePanel>

          <div class="space-y-2 overflow-y-auto min-h-0 flex-1">
          <DemoStep step={1} color="cyan"
            label="Try listing pods with no permissions"
            expanded={openStep === 1} ontoggle={toggle}
            commands={[
              'kubectl --as=system:serviceaccount:demo:demo-user -n demo get pods',
            ]}
          />
          <DemoStep step={2} color="cyan"
            label="Create a Role allowing pod access"
            expanded={openStep === 2} ontoggle={toggle}
            commands={[
              `kubectl apply -f - <<'EOF'
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: pod-reader
  namespace: demo
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list"]
EOF`,
            ]}
          />
          <DemoStep step={3} color="cyan"
            label="Bind the role to demo-user"
            expanded={openStep === 3} ontoggle={toggle}
            commands={[
              `kubectl apply -f - <<'EOF'
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
EOF`,
            ]}
          />
          <DemoStep step={4} color="cyan"
            label="Retry — should now succeed"
            expanded={openStep === 4} ontoggle={toggle}
            commands={[
              'kubectl --as=system:serviceaccount:demo:demo-user -n demo get pods',
            ]}
          />
          <DemoStep step={5} color="cyan"
            label="Try deleting — least privilege in action"
            expanded={openStep === 5} ontoggle={toggle}
            commands={[
              'kubectl --as=system:serviceaccount:demo:demo-user -n demo delete pod web',
            ]}
          />
          </div>
        </div>
      {/snippet}
      {#snippet right()}
        <Terminal name="rbac" class="w-full h-full" />
      {/snippet}
    </SplitLayout>
  </div>
</Slide>
