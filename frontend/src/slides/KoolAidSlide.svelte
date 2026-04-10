<script>
  import { onMount } from 'svelte';
  import { gsap } from 'gsap';
  import Slide from '../lib/Slide.svelte';
  import koolaidImg from '../assets/k8s-koolaid.webp';

  let imgEl;

  onMount(() => {
    gsap.from(imgEl, { opacity: 0, scale: 0.9, duration: 0.8, ease: 'back.out(1.7)' });
  });
</script>

<Slide class="items-center justify-center" animate={false}>
  <img bind:this={imgEl} src={koolaidImg} alt="Kubernetes Kool-Aid" class="max-h-[80vh] max-w-[80vw] rounded-xl" />
</Slide>
