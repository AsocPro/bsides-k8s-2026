<script>
  import { onMount } from 'svelte';
  import SlideEngine from './lib/SlideEngine.svelte';
  import { slides, currentSlide, currentIndex } from './lib/slides.js';

  import TitleSlide from './slides/TitleSlide.svelte';
  import KoolAidSlide from './slides/KoolAidSlide.svelte';
  import EvolutionSlide from './slides/EvolutionSlide.svelte';
  import K8sIntroSlide from './slides/K8sIntroSlide.svelte';
  import K8sSecuritySlide from './slides/K8sSecuritySlide.svelte';
  import RBACSlide from './slides/RBACSlide.svelte';
  import RBACDemoSlide from './slides/RBACDemoSlide.svelte';
  import PolicySlide from './slides/PolicySlide.svelte';
  import PolicyDemoSlide from './slides/PolicyDemoSlide.svelte';
  import NetworkPolicySlide from './slides/NetworkPolicySlide.svelte';
  import NetPolDemoSlide from './slides/NetPolDemoSlide.svelte';
  import FutureSlide from './slides/FutureSlide.svelte';
  import EndSlide from './slides/EndSlide.svelte';

  onMount(() => {
    slides.set([
      { component: TitleSlide, id: 'title' },
      { component: KoolAidSlide, id: 'koolaid' },
      { component: EvolutionSlide, id: 'evolution' },
      { component: K8sIntroSlide, id: 'k8s-intro' },
      { component: K8sSecuritySlide, id: 'k8s-security' },
      { component: RBACSlide, id: 'rbac' },
      { component: RBACDemoSlide, id: 'rbac-demo' },
      { component: PolicySlide, id: 'policy' },
      { component: PolicyDemoSlide, id: 'policy-demo' },
      { component: NetworkPolicySlide, id: 'netpol' },
      { component: NetPolDemoSlide, id: 'netpol-demo' },
      { component: FutureSlide, id: 'future' },
      { component: EndSlide, id: 'end' },
    ]);
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

<style>
  .slide-container {
    width: 100%;
    height: 100%;
    position: absolute;
    inset: 0;
  }
</style>
