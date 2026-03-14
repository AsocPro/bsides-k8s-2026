<script>
  import Slide from '../lib/Slide.svelte';
  import SplitLayout from '../lib/SplitLayout.svelte';
  import Terminal from '../lib/Terminal.svelte';
</script>

<Slide padding={false}>
  <div class="p-8 pb-2">
    <h2 class="text-3xl font-bold text-white">Network Policy Demo</h2>
    <p class="text-neutral-500 text-sm mt-1">Live — microsegmentation in action</p>
  </div>

  <div class="flex-1 min-h-0 px-8 pb-8">
    <SplitLayout ratio="2fr 3fr">
      {#snippet left()}
        <div class="space-y-4 text-sm text-neutral-400 pt-4">
          <div class="p-3 rounded-lg bg-neutral-800/50 border border-neutral-700/50">
            <p class="text-red-400 font-mono text-xs mb-1">Step 1</p>
            <p>Deploy 3-tier app: frontend, API, database</p>
          </div>
          <div class="p-3 rounded-lg bg-neutral-800/50 border border-neutral-700/50">
            <p class="text-red-400 font-mono text-xs mb-1">Step 2</p>
            <p>Show frontend can reach database directly — the problem</p>
          </div>
          <div class="p-3 rounded-lg bg-neutral-800/50 border border-neutral-700/50">
            <p class="text-red-400 font-mono text-xs mb-1">Step 3</p>
            <p>Apply default-deny — everything goes dark</p>
          </div>
          <div class="p-3 rounded-lg bg-neutral-800/50 border border-neutral-700/50">
            <p class="text-red-400 font-mono text-xs mb-1">Step 4</p>
            <p>Restore allowed paths one by one</p>
          </div>
          <div class="p-3 rounded-lg bg-neutral-800/50 border border-neutral-700/50">
            <p class="text-red-400 font-mono text-xs mb-1">Step 5</p>
            <p>Frontend → database still blocked — blast radius contained</p>
          </div>
        </div>
      {/snippet}
      {#snippet right()}
        <Terminal name="netpol" class="w-full h-full" />
      {/snippet}
    </SplitLayout>
  </div>
</Slide>
