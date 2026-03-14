<script>
  import { onMount } from 'svelte';
  import { gsap } from 'gsap';
  import Slide from '../lib/Slide.svelte';

  let timeline;

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

  onMount(() => {
    if (!timeline) return;
    const cards = timeline.querySelectorAll('.era-card');
    const connector = timeline.querySelectorAll('.connector');
    const heading = timeline.querySelector('h2');

    const tl = gsap.timeline();
    tl.from(heading, { opacity: 0, y: -20, duration: 0.5 });

    cards.forEach((card, i) => {
      tl.from(card, {
        opacity: 0,
        x: -40,
        duration: 0.5,
        ease: 'power2.out',
      }, `-=0.1`);
      if (connector[i]) {
        tl.from(connector[i], {
          scaleX: 0,
          duration: 0.3,
          ease: 'power1.out',
        }, '-=0.2');
      }
    });
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
