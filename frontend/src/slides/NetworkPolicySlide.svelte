<script>
  import Slide from '../lib/Slide.svelte';
  import AnimatedList from '../lib/AnimatedList.svelte';

  const benefits = [
    'Microsegmentation — each service only talks to its direct dependencies',
    'Defense in depth — network controls remain even if app controls fail',
    'Reduced lateral movement — attackers can\'t pivot from frontend to database',
    'Compliance — demonstrate network isolation for PCI-DSS, SOC2',
  ];
</script>

<Slide>
  <h2 class="text-4xl font-bold text-white mb-2">Network Policies</h2>
  <p class="text-neutral-500 mb-6">Limiting blast radius when something is exploited</p>

  <div class="mb-6">
    <p class="text-neutral-300 leading-relaxed">
      By default, all pods can communicate with each other.
      Network policies implement <span class="text-red-400 font-semibold">microsegmentation</span> — explicit
      allow-lists for traffic between pods.
    </p>
  </div>

  <AnimatedList items={benefits} />
</Slide>
