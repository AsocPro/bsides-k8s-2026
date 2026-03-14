<script>
  import { onMount } from 'svelte';
  import { gsap } from 'gsap';
  import { latestCommand, connected } from './events.js';

  let cmdEl;
  let prevCommand = null;

  // Animate when a new command comes in
  $effect(() => {
    const cmd = $latestCommand;
    if (cmd && cmd !== prevCommand && cmdEl) {
      prevCommand = cmd;
      // Flash in the command overlay
      gsap.fromTo(cmdEl,
        { opacity: 0, y: 10 },
        { opacity: 1, y: 0, duration: 0.3, ease: 'power2.out' }
      );
      // Fade out after a few seconds
      gsap.to(cmdEl, {
        opacity: 0,
        duration: 0.5,
        delay: 4,
        ease: 'power2.in',
      });
    }
  });
</script>

<!-- Connection indicator -->
<div class="fixed top-4 right-4 flex items-center gap-2 z-50">
  <div
    class="w-2 h-2 rounded-full"
    class:bg-emerald-500={$connected}
    class:bg-red-500={!$connected}
  ></div>
  <span class="text-xs text-neutral-600 font-mono">
    {$connected ? 'live' : 'offline'}
  </span>
</div>

<!-- Command echo overlay -->
{#if $latestCommand}
  <div
    bind:this={cmdEl}
    class="fixed bottom-16 left-1/2 -translate-x-1/2 z-50 opacity-0"
  >
    <div class="px-5 py-3 rounded-xl bg-neutral-900/90 border border-neutral-700/50 backdrop-blur-sm shadow-2xl">
      <div class="flex items-center gap-3">
        <span class="text-emerald-500 text-xs font-mono">$</span>
        <code class="text-neutral-200 text-sm font-mono">{$latestCommand.command}</code>
        {#if $latestCommand.exitCode !== undefined && $latestCommand.exitCode !== 0}
          <span class="text-red-400 text-xs font-mono">exit {$latestCommand.exitCode}</span>
        {/if}
      </div>
    </div>
  </div>
{/if}
