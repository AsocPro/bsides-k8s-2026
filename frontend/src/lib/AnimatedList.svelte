<script>
  /**
   * Renders a list of items with staggered entry animations.
   * Each item fades and slides in sequentially.
   */
  import { onMount } from 'svelte';
  import { gsap } from 'gsap';

  let {
    items = [],
    class: className = '',
    stagger = 0.15,
    delay = 0.2,
  } = $props();

  let listEl;

  onMount(() => {
    if (listEl) {
      gsap.from(listEl.children, {
        opacity: 0,
        x: -20,
        duration: 0.4,
        stagger,
        delay,
        ease: 'power2.out',
      });
    }
  });
</script>

<ul bind:this={listEl} class="space-y-3 {className}">
  {#each items as item}
    <li class="flex items-start gap-3 text-neutral-300">
      <span class="text-cyan-500 mt-1 shrink-0">&#x25B8;</span>
      <span>{item}</span>
    </li>
  {/each}
</ul>
