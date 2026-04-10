<script>
  import { onMount } from 'svelte';
  import { gsap } from 'gsap';
  import Slide from '../lib/Slide.svelte';
  import feedbackQr from '../assets/feedback-qr.png';

  let heading;
  let name;
  let qrBlock;

  onMount(() => {
    const tl = gsap.timeline();
    tl.from(heading, { opacity: 0, y: 30, duration: 0.8, ease: 'power3.out' })
      .from(name, { opacity: 0, y: 20, duration: 0.6, ease: 'power3.out' }, '-=0.3')
      .from(qrBlock, { opacity: 0, scale: 0.9, duration: 0.5, ease: 'back.out(1.7)' }, '-=0.2');
  });
</script>

<Slide class="items-center justify-center" animate={false}>
  <h1 bind:this={heading} class="text-6xl font-bold text-white text-center">
    Questions?
  </h1>
  <p bind:this={name} class="text-2xl text-neutral-400 mt-4 text-center">
    Zach Gibbs
  </p>
  <div bind:this={qrBlock} class="mt-10 flex flex-col items-center gap-3">
    <img src={feedbackQr} alt="Session feedback QR code" class="w-48 h-48 rounded-lg" />
    <p class="text-sm text-neutral-500">Scan for session feedback</p>
  </div>
</Slide>
