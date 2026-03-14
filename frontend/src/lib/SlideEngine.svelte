<script>
  import { onMount } from 'svelte';
  import { slides, currentIndex, currentSlide, totalSlides, next, prev } from './slides.js';

  let { children } = $props();

  onMount(() => {
    function handleKey(e) {
      if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') return;
      // Don't capture keys when terminal iframe is focused
      if (document.activeElement?.tagName === 'IFRAME') return;

      switch (e.key) {
        case 'ArrowRight':
        case ' ':
        case 'PageDown':
          e.preventDefault();
          next();
          break;
        case 'ArrowLeft':
        case 'PageUp':
          e.preventDefault();
          prev();
          break;
      }
    }

    window.addEventListener('keydown', handleKey);
    return () => window.removeEventListener('keydown', handleKey);
  });
</script>

<div class="slide-engine">
  {@render children()}

  <!-- Slide counter -->
  <div class="fixed bottom-4 right-4 text-xs text-neutral-600 font-mono select-none">
    {$currentIndex + 1} / {$totalSlides}
  </div>

  <!-- Nav arrows -->
  <button
    class="fixed left-4 top-1/2 -translate-y-1/2 text-neutral-700 hover:text-neutral-400 transition-colors text-2xl select-none"
    onclick={prev}
    disabled={$currentIndex === 0}
  >
    ‹
  </button>
  <button
    class="fixed right-4 top-1/2 -translate-y-1/2 text-neutral-700 hover:text-neutral-400 transition-colors text-2xl select-none"
    onclick={next}
    disabled={$currentIndex === $totalSlides - 1}
  >
    ›
  </button>
</div>

<style>
  .slide-engine {
    width: 100%;
    height: 100%;
    position: relative;
  }

  button:disabled {
    opacity: 0.2;
    cursor: default;
  }
</style>
