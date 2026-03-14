<script>
  let { name = 'demo', class: className = '' } = $props();

  // ttyd is served at /terminal/{name}/ via the Go backend proxy
  const src = `/terminal/${name}/`;
</script>

<div class="terminal-wrapper {className}">
  <iframe
    {src}
    title="Terminal: {name}"
    class="w-full h-full border-0 rounded-lg"
    allow="clipboard-read; clipboard-write"
  ></iframe>
  <div class="terminal-label">
    <span class="text-[10px] font-mono text-neutral-600 uppercase tracking-wider">{name}</span>
  </div>
</div>

<style>
  .terminal-wrapper {
    background: #1a1a2e;
    border: 1px solid #333;
    border-radius: 0.5rem;
    overflow: hidden;
    position: relative;
  }
  .terminal-label {
    position: absolute;
    top: 0.25rem;
    right: 0.5rem;
    pointer-events: none;
  }
</style>
