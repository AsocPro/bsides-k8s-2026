<script>
  import { onMount } from 'svelte';
  import { gsap } from 'gsap';
  import Terminal from '../lib/Terminal.svelte';

  let container;

  onMount(() => {
    gsap.from(container.querySelector('h2'), {
      opacity: 0, x: -30, duration: 0.5, ease: 'power2.out',
    });
    gsap.from(container.querySelector('.terminal-panel'), {
      opacity: 0, y: 20, duration: 0.6, delay: 0.2, ease: 'power2.out',
    });
  });
</script>

<div bind:this={container} class="flex flex-col h-full p-8 gap-4">
  <h2 class="text-3xl font-bold text-white">Live Demo</h2>
  <div class="terminal-panel flex-1 min-h-0">
    <Terminal name="demo" class="w-full h-full" />
  </div>
</div>
