<script>
  import { onMount } from 'svelte';
  import { gsap } from 'gsap';

  let { data } = $props();

  const resources = ['pods', 'roles', 'rolebindings'];
  const verbs = ['get', 'create', 'delete'];

  let cellEls = {};
  let prevState = {};

  function cellKey(resource, verb) {
    return `${resource}-${verb}`;
  }

  $effect(() => {
    if (!data?.permissions) return;
    for (const resource of resources) {
      for (const verb of verbs) {
        const key = cellKey(resource, verb);
        const el = cellEls[key];
        const newVal = data.permissions[resource]?.[verb];
        if (el && prevState[key] !== newVal) {
          prevState[key] = newVal;
          gsap.fromTo(el,
            { scale: 0, opacity: 0 },
            { scale: 1, opacity: 1, duration: 0.4, ease: 'back.out(2)' }
          );
        }
      }
    }
  });
</script>

<div class="grid gap-1" style="grid-template-columns: auto repeat(3, 1fr);">
  <!-- Header row -->
  <div></div>
  {#each resources as resource}
    <div class="text-center text-[10px] font-mono text-cyan-400/70 uppercase tracking-wider px-1 py-0.5">
      {resource}
    </div>
  {/each}

  <!-- Data rows -->
  {#each verbs as verb}
    <div class="text-right text-[10px] font-mono text-neutral-500 pr-2 flex items-center justify-end">
      {verb}
    </div>
    {#each resources as resource}
      {@const allowed = data.permissions[resource]?.[verb]}
      <div
        bind:this={cellEls[cellKey(resource, verb)]}
        class="flex items-center justify-center rounded h-7
               {allowed ? 'bg-emerald-500/15' : 'bg-red-500/15'}"
      >
        {#if allowed}
          <svg class="w-4 h-4 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7" />
          </svg>
        {:else}
          <svg class="w-4 h-4 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12" />
          </svg>
        {/if}
      </div>
    {/each}
  {/each}
</div>
