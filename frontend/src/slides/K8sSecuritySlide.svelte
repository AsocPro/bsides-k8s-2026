<script>
  import Slide from '../lib/Slide.svelte';
  import AnimatedList from '../lib/AnimatedList.svelte';
  import { onMount } from 'svelte';
  import { gsap } from 'gsap';

  let diagramEl;

  const securityFeatures = [
    'All changes go through the API server — fully auditable',
    'Secrets management is a first-class concept',
    'Resource isolation via namespaces, quotas, and security contexts',
    'Security posture defined in code — reviewable and version controlled',
  ];

  onMount(() => {
    if (!diagramEl) return;
    const tl = gsap.timeline({ delay: 0.3 });

    // API server appears first
    tl.from(diagramEl.querySelector('.api-server'), {
      scale: 0.8, opacity: 0, duration: 0.5, ease: 'power2.out',
    });

    // Spokes radiate out
    tl.from(diagramEl.querySelectorAll('.spoke'), {
      opacity: 0, strokeDashoffset: 80, duration: 0.4, stagger: 0.08, ease: 'power2.out',
    });

    // Feature nodes pop in
    tl.from(diagramEl.querySelectorAll('.feature-node'), {
      opacity: 0, scale: 0, duration: 0.35, stagger: 0.1, ease: 'back.out(1.7)',
    }, '-=0.2');

    // Icons/text inside nodes
    tl.from(diagramEl.querySelectorAll('.node-label'), {
      opacity: 0, duration: 0.3, stagger: 0.08, ease: 'power2.out',
    }, '-=0.2');

    // Shield outline pulses
    tl.from(diagramEl.querySelector('.shield'), {
      opacity: 0, scale: 0.9, duration: 0.5, ease: 'power2.out',
    }, '-=0.3');
  });
</script>

<Slide>
  <h2 class="text-4xl font-bold text-white mb-6">Kubernetes Built-in Security</h2>

  <div class="grid grid-cols-2 gap-8 flex-1 min-h-0">
    <!-- Left: Diagram -->
    <div bind:this={diagramEl} class="flex items-center justify-center">
      <svg viewBox="0 0 320 320" class="w-full h-full max-h-[420px]" xmlns="http://www.w3.org/2000/svg">
        <!-- Subtle shield outline behind everything -->
        <path class="shield" d="M160,28 L260,70 C260,180 220,260 160,295 C100,260 60,180 60,70 Z"
              fill="rgba(16,185,129,0.04)" stroke="rgb(16,185,129)" stroke-width="1" stroke-dasharray="8 4" opacity="0.5" />

        <!-- Spokes from center to each node -->
        <line class="spoke" x1="160" y1="160" x2="160" y2="62" stroke="rgb(100,116,139)" stroke-width="1" stroke-dasharray="4 3" />
        <line class="spoke" x1="160" y1="160" x2="262" y2="130" stroke="rgb(100,116,139)" stroke-width="1" stroke-dasharray="4 3" />
        <line class="spoke" x1="160" y1="160" x2="245" y2="245" stroke="rgb(100,116,139)" stroke-width="1" stroke-dasharray="4 3" />
        <line class="spoke" x1="160" y1="160" x2="75" y2="245" stroke="rgb(100,116,139)" stroke-width="1" stroke-dasharray="4 3" />
        <line class="spoke" x1="160" y1="160" x2="58" y2="130" stroke="rgb(100,116,139)" stroke-width="1" stroke-dasharray="4 3" />

        <!-- Center: API Server -->
        <g class="api-server">
          <circle cx="160" cy="160" r="36" fill="rgba(99,102,241,0.15)" stroke="rgb(99,102,241)" stroke-width="2" />
          <text x="160" y="155" text-anchor="middle" fill="rgb(165,180,252)" font-size="10" font-weight="700">API</text>
          <text x="160" y="170" text-anchor="middle" fill="rgb(165,180,252)" font-size="10" font-weight="700">Server</text>
        </g>

        <!-- Node 1: Audit Logs (top center) -->
        <g class="feature-node" transform="translate(160,52)">
          <rect x="-38" y="-18" width="76" height="36" rx="8" fill="rgba(251,191,36,0.1)" stroke="rgb(251,191,36)" stroke-width="1.2" />
          <text class="node-label" x="0" y="-3" text-anchor="middle" fill="rgb(251,191,36)" font-size="14">&#128221;</text>
          <text class="node-label" x="0" y="11" text-anchor="middle" fill="rgb(253,224,71)" font-size="8" font-weight="500">Audit Logs</text>
        </g>

        <!-- Node 2: Secrets (right top) -->
        <g class="feature-node" transform="translate(268,125)">
          <rect x="-35" y="-18" width="70" height="36" rx="8" fill="rgba(244,114,182,0.1)" stroke="rgb(244,114,182)" stroke-width="1.2" />
          <text class="node-label" x="0" y="-3" text-anchor="middle" fill="rgb(244,114,182)" font-size="14">&#128272;</text>
          <text class="node-label" x="0" y="11" text-anchor="middle" fill="rgb(251,146,201)" font-size="8" font-weight="500">Secrets</text>
        </g>

        <!-- Node 3: Security Contexts (right bottom) -->
        <g class="feature-node" transform="translate(248,252)">
          <rect x="-42" y="-18" width="84" height="36" rx="8" fill="rgba(34,211,238,0.1)" stroke="rgb(34,211,238)" stroke-width="1.2" />
          <text class="node-label" x="0" y="-3" text-anchor="middle" fill="rgb(34,211,238)" font-size="14">&#128737;</text>
          <text class="node-label" x="0" y="11" text-anchor="middle" fill="rgb(103,232,249)" font-size="8" font-weight="500">Sec Contexts</text>
        </g>

        <!-- Node 4: Namespaces/Quotas (left bottom) -->
        <g class="feature-node" transform="translate(72,252)">
          <rect x="-42" y="-18" width="84" height="36" rx="8" fill="rgba(6,182,212,0.1)" stroke="rgb(6,182,212)" stroke-width="1.2" />
          <text class="node-label" x="0" y="-3" text-anchor="middle" fill="rgb(6,182,212)" font-size="14">&#128230;</text>
          <text class="node-label" x="0" y="11" text-anchor="middle" fill="rgb(103,232,249)" font-size="8" font-weight="500">NS &amp; Quotas</text>
        </g>

        <!-- Node 5: GitOps / Code Review (left top) -->
        <g class="feature-node" transform="translate(52,125)">
          <rect x="-35" y="-18" width="70" height="36" rx="8" fill="rgba(16,185,129,0.1)" stroke="rgb(16,185,129)" stroke-width="1.2" />
          <text class="node-label" x="0" y="-3" text-anchor="middle" fill="rgb(16,185,129)" font-size="14">&#128196;</text>
          <text class="node-label" x="0" y="11" text-anchor="middle" fill="rgb(52,211,153)" font-size="8" font-weight="500">As Code</text>
        </g>
      </svg>
    </div>

    <!-- Right: Text list -->
    <div class="flex flex-col justify-center">
      <AnimatedList items={securityFeatures} delay={0.8} />
    </div>
  </div>
</Slide>
