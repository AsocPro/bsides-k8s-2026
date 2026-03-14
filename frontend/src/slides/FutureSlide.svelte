<script>
  import { onMount } from 'svelte';
  import { gsap } from 'gsap';
  import Slide from '../lib/Slide.svelte';

  let container;

  const techs = [
    {
      name: 'eBPF / Cilium',
      desc: 'Kernel-level security — L7 network policies, runtime detection, zero overhead.',
      color: 'text-purple-400',
      border: 'border-purple-500/30',
      bg: 'bg-purple-500/10',
    },
    {
      name: 'Workload Identity',
      desc: 'Short-lived credentials, automatic rotation, no static secrets to leak.',
      color: 'text-blue-400',
      border: 'border-blue-500/30',
      bg: 'bg-blue-500/10',
    },
    {
      name: 'Talos Linux',
      desc: 'No SSH, no shell, no userspace. Minimal attack surface by design.',
      color: 'text-emerald-400',
      border: 'border-emerald-500/30',
      bg: 'bg-emerald-500/10',
    },
    {
      name: 'SPIFFE / SPIRE',
      desc: 'Universal workload identity — cryptographic proof of who is talking to whom.',
      color: 'text-amber-400',
      border: 'border-amber-500/30',
      bg: 'bg-amber-500/10',
    },
  ];

  onMount(() => {
    if (!container) return;
    const cards = container.querySelectorAll('.tech-card');
    gsap.from(cards, {
      opacity: 0,
      y: 30,
      duration: 0.5,
      stagger: 0.12,
      ease: 'power2.out',
    });
  });
</script>

<Slide>
  <h2 class="text-4xl font-bold text-white mb-2">Looking Forward</h2>
  <p class="text-neutral-500 mb-8">The next generation of secure deployments</p>

  <div bind:this={container} class="grid grid-cols-2 gap-4 flex-1 min-h-0">
    {#each techs as tech}
      <div class="tech-card p-5 rounded-xl border {tech.border} {tech.bg} flex flex-col">
        <h3 class="text-lg font-semibold {tech.color} mb-2">{tech.name}</h3>
        <p class="text-neutral-400 text-sm leading-relaxed">{tech.desc}</p>
      </div>
    {/each}
  </div>
</Slide>
