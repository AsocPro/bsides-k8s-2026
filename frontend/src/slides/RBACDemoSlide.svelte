<script>
  import Slide from '../lib/Slide.svelte';
  import SplitLayout from '../lib/SplitLayout.svelte';
  import Terminal from '../lib/Terminal.svelte';
  import DemoStep from '../lib/DemoStep.svelte';
</script>

<Slide padding={false}>
  <div class="p-8 pb-2">
    <h2 class="text-3xl font-bold text-white">RBAC Demo</h2>
    <p class="text-neutral-500 text-sm mt-1">Live — try creating roles and testing access</p>
  </div>

  <div class="flex-1 min-h-0 px-8 pb-8">
    <SplitLayout ratio="2fr 3fr">
      {#snippet left()}
        <div class="space-y-2 text-sm pt-4 overflow-y-auto">
          <DemoStep step={1} color="cyan"
            label="Try listing pods with no permissions"
            commands={[
              'kubectl --as=system:serviceaccount:demo:demo-user -n demo get pods',
            ]}
          />
          <DemoStep step={2} color="cyan"
            label="Create a Role allowing pod access"
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
            commands={[
              'kubectl --as=system:serviceaccount:demo:demo-user -n demo get pods',
            ]}
          />
          <DemoStep step={5} color="cyan"
            label="Try deleting — least privilege in action"
            commands={[
              'kubectl --as=system:serviceaccount:demo:demo-user -n demo delete pod web',
            ]}
          />
        </div>
      {/snippet}
      {#snippet right()}
        <Terminal name="rbac" class="w-full h-full" />
      {/snippet}
    </SplitLayout>
  </div>
</Slide>
