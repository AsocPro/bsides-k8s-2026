<script>
  import Slide from '../lib/Slide.svelte';
  import SplitLayout from '../lib/SplitLayout.svelte';
  import Terminal from '../lib/Terminal.svelte';
  import DemoStep from '../lib/DemoStep.svelte';
</script>

<Slide padding={false}>
  <div class="p-8 pb-2">
    <h2 class="text-3xl font-bold text-white">Network Policy Demo</h2>
    <p class="text-neutral-500 text-sm mt-1">Live — microsegmentation in action</p>
  </div>

  <div class="flex-1 min-h-0 px-8 pb-8">
    <SplitLayout ratio="2fr 3fr">
      {#snippet left()}
        <div class="space-y-2 text-sm pt-4 overflow-y-auto">
          <DemoStep step={1} color="red"
            label="Show the running 3-tier app"
            commands={[
              'kubectl -n demo get pods,svc',
            ]}
          />
          <DemoStep step={2} color="red"
            label="Everything can talk to everything — the problem"
            commands={[
              'kubectl -n demo exec deployment/frontend -- curl -s --max-time 3 api',
              'kubectl -n demo exec deployment/frontend -- curl -s --max-time 3 database',
              'kubectl -n demo exec deployment/api -- curl -s --max-time 3 database',
            ]}
          />
          <DemoStep step={3} color="red"
            label="Apply default-deny — everything goes dark"
            commands={[
              'kubectl apply -f /root/manifests/deny-all.yaml',
              'kubectl -n demo exec deployment/frontend -- curl -s --max-time 3 api',
            ]}
          />
          <DemoStep step={4} color="red"
            label="Restore allowed paths one by one"
            commands={[
              'kubectl apply -f /root/manifests/allow-frontend-to-api.yaml',
              'kubectl apply -f /root/manifests/allow-api-to-db.yaml',
              'kubectl -n demo exec deployment/frontend -- curl -s --max-time 3 api',
              'kubectl -n demo exec deployment/api -- curl -s --max-time 3 database',
            ]}
          />
          <DemoStep step={5} color="red"
            label="Frontend → database still blocked — blast radius contained"
            commands={[
              'kubectl -n demo exec deployment/frontend -- curl -s --max-time 3 database',
            ]}
          />
        </div>
      {/snippet}
      {#snippet right()}
        <Terminal name="netpol" class="w-full h-full" />
      {/snippet}
    </SplitLayout>
  </div>
</Slide>
