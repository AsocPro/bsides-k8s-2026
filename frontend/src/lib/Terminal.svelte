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
</div>

<style>
  .terminal-wrapper {
    background: #1a1a2e;
    border: 1px solid #333;
    border-radius: 0.5rem;
    overflow: hidden;
  }
</style>
