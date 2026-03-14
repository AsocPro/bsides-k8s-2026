<script>
  /**
   * Base slide wrapper. Provides consistent padding, background,
   * and entry animation for any slide content.
   */
  import { onMount } from 'svelte';
  import { gsap } from 'gsap';

  let {
    children,
    class: className = '',
    padding = true,
    animate = true,
  } = $props();

  let el;

  onMount(() => {
    if (animate && el) {
      gsap.from(el, { opacity: 0, duration: 0.3, ease: 'power2.out' });
    }
  });
</script>

<div
  bind:this={el}
  class="slide {className}"
  class:p-12={padding}
>
  {@render children()}
</div>

<style>
  .slide {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
  }
</style>
