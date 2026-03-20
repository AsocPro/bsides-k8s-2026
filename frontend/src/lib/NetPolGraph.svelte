<script>
  import { onMount } from 'svelte';
  import { gsap } from 'gsap';

  let { data } = $props();

  // Node positions (inverted triangle)
  const nodes = [
    { id: 'frontend', label: 'Frontend', x: 60, y: 25 },
    { id: 'api', label: 'API', x: 220, y: 25 },
    { id: 'database', label: 'Database', x: 140, y: 100 },
  ];

  const edges = [
    { from: 'frontend', to: 'api', key: 'frontend->api' },
    { from: 'frontend', to: 'database', key: 'frontend->database' },
    { from: 'api', to: 'database', key: 'api->database' },
  ];

  function getNode(id) {
    return nodes.find(n => n.id === id);
  }

  let lineEls = {};

  $effect(() => {
    if (!data?.connections) return;
    for (const edge of edges) {
      const el = lineEls[edge.key];
      if (el) {
        const connected = data.connections[edge.key];
        gsap.to(el, {
          attr: {
            stroke: connected ? '#34d399' : '#f87171',
            'stroke-dasharray': connected ? 'none' : '6 4',
          },
          duration: 0.5,
          ease: 'power2.out',
        });
      }
    }
  });
</script>

<svg viewBox="0 0 280 130" class="w-full max-w-xs">
  <!-- Edges -->
  {#each edges as edge}
    {@const from = getNode(edge.from)}
    {@const to = getNode(edge.to)}
    {@const connected = data?.connections?.[edge.key] ?? true}
    <line
      bind:this={lineEls[edge.key]}
      x1={from.x} y1={from.y}
      x2={to.x} y2={to.y}
      stroke={connected ? '#34d399' : '#f87171'}
      stroke-width="2"
      stroke-dasharray={connected ? 'none' : '6 4'}
      stroke-linecap="round"
    />
  {/each}

  <!-- Nodes -->
  {#each nodes as node}
    <rect
      x={node.x - 32} y={node.y - 12}
      width="64" height="24"
      rx="6"
      fill="#1e1e1e"
      stroke="#525252"
      stroke-width="1"
    />
    <text
      x={node.x} y={node.y + 4}
      text-anchor="middle"
      class="text-[10px] font-mono fill-neutral-300"
    >
      {node.label}
    </text>
  {/each}
</svg>
