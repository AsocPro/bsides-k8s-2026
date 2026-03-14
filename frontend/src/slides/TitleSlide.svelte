<script>
  import { onMount } from 'svelte';
  import { gsap } from 'gsap';
  import Slide from '../lib/Slide.svelte';

  let title;
  let subtitle;
  let badge;

  onMount(() => {
    const tl = gsap.timeline();
    tl.from(title, { opacity: 0, y: 50, duration: 1, ease: 'power3.out' })
      .from(subtitle, { opacity: 0, y: 20, duration: 0.7, ease: 'power3.out' }, '-=0.4')
      .from(badge, { opacity: 0, scale: 0.8, duration: 0.5, ease: 'back.out(1.7)' }, '-=0.2');
  });
</script>

<Slide class="items-center justify-center" animate={false}>
  <h1 bind:this={title} class="text-7xl font-bold tracking-tight text-white text-center">
    Kubernetes
  </h1>
  <p bind:this={subtitle} class="text-2xl text-neutral-400 mt-6 text-center">
    Orchestrating a More Secure Future
  </p>
  <div bind:this={badge} class="mt-12 px-5 py-2 rounded-full bg-neutral-800/80 border border-neutral-700 text-sm text-neutral-400">
    BSides 2026
  </div>
</Slide>
