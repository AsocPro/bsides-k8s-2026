<script>
  import { onDestroy } from 'svelte';
  import { checkState, startPolling, stopPolling } from './stores/stateCheck.js';

  let {
    demo,
    store,
    color = 'cyan',
    children,
  } = $props();

  let monitoring = $state(false);

  function toggleMonitor() {
    monitoring = !monitoring;
    if (monitoring) {
      startPolling(demo);
    } else {
      stopPolling(demo);
    }
  }

  onDestroy(() => stopPolling(demo));

  const colorMap = {
    cyan: {
      border: 'border-cyan-500/30',
      bg: 'bg-cyan-500/10',
      text: 'text-cyan-400',
      hover: 'hover:bg-cyan-500/20',
      ring: 'ring-cyan-500/30',
    },
    amber: {
      border: 'border-amber-500/30',
      bg: 'bg-amber-500/10',
      text: 'text-amber-400',
      hover: 'hover:bg-amber-500/20',
      ring: 'ring-amber-500/30',
    },
    red: {
      border: 'border-red-500/30',
      bg: 'bg-red-500/10',
      text: 'text-red-400',
      hover: 'hover:bg-red-500/20',
      ring: 'ring-red-500/30',
    },
  };

  const colors = $derived(colorMap[color] || colorMap.cyan);

  function handleCheck() {
    checkState(demo);
  }
</script>

<div class="rounded-lg border {colors.border} bg-neutral-800/30 overflow-hidden shrink-0">
  <div class="flex items-center gap-2 px-3 py-2">
    <button
      onclick={handleCheck}
      disabled={$store.loading || monitoring}
      class="px-3 py-1.5 rounded-md text-xs font-mono font-medium transition-all
             {colors.bg} {colors.text} {colors.hover} ring-1 {colors.ring}
             disabled:opacity-50 disabled:cursor-not-allowed"
    >
      Check State
    </button>
    <button
      onclick={toggleMonitor}
      class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs font-mono font-medium transition-all ring-1
             {monitoring
               ? 'bg-emerald-500/20 text-emerald-400 hover:bg-emerald-500/30 ring-emerald-500/30'
               : colors.bg + ' ' + colors.text + ' ' + colors.hover + ' ' + colors.ring}"
    >
      {#if monitoring}
        <svg class="w-3 h-3 animate-spin" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3" stroke-linecap="round" class="opacity-25" />
          <path d="M12 2a10 10 0 0 1 10 10" stroke="currentColor" stroke-width="3" stroke-linecap="round" class="opacity-75" />
        </svg>
        Monitoring
      {:else}
        Monitor
      {/if}
    </button>
    {#if $store.error}
      <span class="text-red-400 text-xs truncate">{$store.error}</span>
    {/if}
  </div>

  {#if $store.data}
    <div class="border-t {colors.border} px-3 py-2">
      {@render children()}
    </div>
  {/if}
</div>
