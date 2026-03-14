<script>
  /**
   * Split layout for slides with content on one side and a demo/terminal on the other.
   * Uses named snippets: 'left' and 'right'.
   */
  let {
    left,
    right,
    ratio = '1fr 1fr',
    gap = '2rem',
    class: className = '',
  } = $props();
</script>

<div
  class="split-layout {className}"
  style="grid-template-columns: {ratio}; gap: {gap};"
>
  <div class="split-left">
    {@render left()}
  </div>
  <div class="split-right">
    {@render right()}
  </div>
</div>

<style>
  .split-layout {
    display: grid;
    height: 100%;
    min-height: 0;
  }
  .split-left, .split-right {
    min-height: 0;
    min-width: 0;
    display: flex;
    flex-direction: column;
  }
</style>
