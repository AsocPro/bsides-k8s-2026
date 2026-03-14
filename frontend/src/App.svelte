<script>
  import { onMount } from 'svelte';
  import SlideEngine from './lib/SlideEngine.svelte';
  import EventOverlay from './lib/EventOverlay.svelte';
  import { slides, currentSlide, currentIndex } from './lib/slides.js';
  import { connect, disconnect, send } from './lib/events.js';

  import TitleSlide from './slides/TitleSlide.svelte';
  import EvolutionSlide from './slides/EvolutionSlide.svelte';
  import K8sIntroSlide from './slides/K8sIntroSlide.svelte';
  import RBACSlide from './slides/RBACSlide.svelte';
  import RBACDemoSlide from './slides/RBACDemoSlide.svelte';
  import PolicySlide from './slides/PolicySlide.svelte';
  import NetworkPolicySlide from './slides/NetworkPolicySlide.svelte';
  import FutureSlide from './slides/FutureSlide.svelte';

  onMount(() => {
    slides.set([
      { component: TitleSlide, id: 'title' },
      { component: EvolutionSlide, id: 'evolution' },
      { component: K8sIntroSlide, id: 'k8s-intro' },
      { component: RBACSlide, id: 'rbac' },
      { component: RBACDemoSlide, id: 'rbac-demo' },
      { component: PolicySlide, id: 'policy' },
      { component: NetworkPolicySlide, id: 'netpol' },
      { component: FutureSlide, id: 'future' },
    ]);

    // Connect to event bridge
    connect();
    return disconnect;
  });

  // Broadcast slide changes to the event bridge
  $effect(() => {
    const slide = $currentSlide;
    if (slide) {
      send('slide-change', { id: slide.id, index: $currentIndex });
    }
  });
</script>

<SlideEngine>
  {#if $currentSlide}
    {#key $currentSlide.id}
      <div class="slide-container">
        <svelte:component this={$currentSlide.component} />
      </div>
    {/key}
  {/if}
</SlideEngine>

<EventOverlay />

<style>
  .slide-container {
    width: 100%;
    height: 100%;
    position: absolute;
    inset: 0;
  }
</style>
