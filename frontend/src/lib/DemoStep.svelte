<script>
  let {
    step,
    label,
    commands = [],
    color = 'cyan',
  } = $props();

  let expanded = $state(false);
  let copiedIndex = $state(-1);

  function toggle() {
    if (commands.length > 0) expanded = !expanded;
  }

  async function copy(text, index) {
    try {
      await navigator.clipboard.writeText(text);
      copiedIndex = index;
      setTimeout(() => { copiedIndex = -1; }, 1500);
    } catch {
      // Fallback: select text
    }
  }

  const colorMap = {
    cyan: 'text-cyan-400',
    amber: 'text-amber-400',
    red: 'text-red-400',
    emerald: 'text-emerald-400',
  };
</script>

<div class="rounded-lg bg-neutral-800/50 border border-neutral-700/50 overflow-hidden">
  <button
    class="w-full p-3 text-left flex items-start gap-2 hover:bg-neutral-700/30 transition-colors"
    onclick={toggle}
  >
    <div class="flex-1 min-w-0">
      <p class="{colorMap[color] || 'text-cyan-400'} font-mono text-xs mb-1">Step {step}</p>
      <p class="text-sm text-neutral-400">{label}</p>
    </div>
    {#if commands.length > 0}
      <span class="text-neutral-600 text-xs mt-1 shrink-0 select-none">
        {expanded ? '▾' : '▸'}
      </span>
    {/if}
  </button>

  {#if expanded && commands.length > 0}
    <div class="border-t border-neutral-700/50 bg-neutral-900/50 p-2 space-y-1.5">
      {#each commands as cmd, i}
        <div class="group flex items-start gap-1.5 rounded bg-neutral-950/50 p-1.5">
          <pre class="flex-1 text-xs font-mono text-emerald-400/80 whitespace-pre-wrap break-all overflow-hidden select-all">{cmd}</pre>
          <button
            class="shrink-0 p-1 rounded text-neutral-600 hover:text-white hover:bg-neutral-700 transition-colors"
            onclick={(e) => { e.stopPropagation(); copy(cmd, i); }}
            title="Copy to clipboard"
          >
            {#if copiedIndex === i}
              <svg class="w-3.5 h-3.5 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            {:else}
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <rect x="9" y="9" width="13" height="13" rx="2" stroke-width="2"/>
                <path stroke-width="2" d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
              </svg>
            {/if}
          </button>
        </div>
      {/each}
    </div>
  {/if}
</div>
