<script>
  import { onMount } from 'svelte';
  import { gsap } from 'gsap';
  import Slide from '../lib/Slide.svelte';
  import { registerSubsteps } from '../lib/slides.js';

  let timeline;
  let phase = 0;
  const totalPhases = 3; // bare metal, VMs, containers

  const eras = [
    {
      label: 'Bare Metal',
      icon: '&#x1F5A5;',
      color: 'text-red-400',
      borderColor: 'border-red-500/30',
      bgColor: 'bg-red-500/10',
      desc: 'Direct hardware install. Config drift, slow provisioning, inconsistent patching.',
    },
    {
      label: 'Virtual Machines',
      icon: '&#x1F4E6;',
      color: 'text-amber-400',
      borderColor: 'border-amber-500/30',
      bgColor: 'bg-amber-500/10',
      desc: 'Snapshots & templates improved things. Still heavyweight, still mutable.',
    },
    {
      label: 'Containers',
      icon: '&#x2693;',
      color: 'text-cyan-400',
      borderColor: 'border-cyan-500/30',
      bgColor: 'bg-cyan-500/10',
      desc: 'Immutable, lightweight, reproducible. Same image everywhere.',
    },
  ];

  function revealPhase(p) {
    if (!timeline) return;
    const cards = timeline.querySelectorAll('.era-card');
    const connectors = timeline.querySelectorAll('.connector');
    const idx = p - 1; // phase is 1-based

    gsap.to(cards[idx], {
      opacity: 1,
      x: 0,
      duration: 0.5,
      ease: 'power2.out',
    });

    // Reveal the connector before this card (connector[0] is between card 0 and 1)
    if (idx > 0 && connectors[idx - 1]) {
      gsap.to(connectors[idx - 1], {
        opacity: 1,
        scaleX: 1,
        duration: 0.3,
        ease: 'power1.out',
      });
    }
  }

  onMount(() => {
    if (!timeline) return;
    const cards = timeline.querySelectorAll('.era-card');
    const connectors = timeline.querySelectorAll('.connector');

    // Hide all cards and connectors initially
    gsap.set(cards, { opacity: 0, x: -40 });
    gsap.set(connectors, { opacity: 0, scaleX: 0 });

    const unregister = registerSubsteps(() => {
      if (phase < totalPhases) {
        phase++;
        revealPhase(phase);
        return true; // consumed
      }
      return false; // let slide engine advance
    });

    return unregister;
  });
</script>

<Slide>
  <h2 class="text-4xl font-bold text-white mb-2">How We Got Here</h2>
  <p class="text-neutral-500 mb-8">The evolution of application delivery</p>

  <div bind:this={timeline} class="flex items-stretch gap-0 flex-1 min-h-0">
    {#each eras as era, i}
      <div class="era-card flex-1 p-6 rounded-xl border {era.borderColor} {era.bgColor} flex flex-col">
        <div class="text-3xl mb-3">{@html era.icon}</div>
        <h3 class="text-xl font-semibold {era.color} mb-2">{era.label}</h3>
        <p class="text-neutral-400 text-sm leading-relaxed">{era.desc}</p>
      </div>
      {#if i < eras.length - 1}
        <div class="connector flex items-center px-2">
          <div class="w-8 h-0.5 bg-neutral-700 origin-left"></div>
          <div class="text-neutral-600">&#x25B6;</div>
        </div>
      {/if}
    {/each}
  </div>
</Slide>
